package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gatewayapi/internal/http-server/server"
	"gatewayapi/internal/middleware/auth/pkg/jwt_checker"
	"gatewayapi/internal/middleware/ratelimiter"
	"gatewayapi/internal/repository"
	"gatewayapi/model"
	"gatewayapi/utils"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
)

var reqCount int64 // global counter
var stopRequestMonitor = make(chan struct{})

type App struct {
	httpserver  *server.HttpServer
	jwtchecker  *jwt_checker.JWTChecker
	ratelimiter *ratelimiter.RateLimiter
	gateWayrepo *repository.GateWayRepository
	gmodel      *model.GatewayModel
}

func NewAPIGatewayApp() *App {
	var app = &App{}
	app.init()
	return app
}

func (a *App) Start() {
	a.startWorkers(4)
	go a.requestsMonitor(stopRequestMonitor)
	if err := a.httpserver.Start(); err != nil {
		log.Fatalf("❌ Failed to start: %v", err)
	}
}

func (a *App) Stop() {
	close(a.gmodel.RequestQueue)
	close(stopRequestMonitor)
	if err := a.httpserver.Stop(); err != nil {
		log.Printf("⚠️ Error stopping server: %v", err)
	}
	log.Println("✅ Server stopped gracefully")
}

func (a *App) init() {
	a.jwtchecker = jwt_checker.NewJWTChecker()
	a.gmodel = model.NewGatewayModel()
	a.gateWayrepo = repository.NewGateWayRepository()
	a.ratelimiter = ratelimiter.NewRateLimiter(*a.gateWayrepo.RateLimiterModel)

	// Khởi tạo router
	router := mux.NewRouter()

	// Đăng ký route cho từng endpoint
	for _, sg := range a.gmodel.ServiceGroups {
		for _, ep := range sg.Endpoints {
			topic := sg.Name + "/" + ep.Name
			router.HandleFunc(ep.Path, a.makeHandler(topic)).Methods(ep.Method)
			log.Printf("Registered route: %s %s -> topic %s", ep.Method, ep.Path, topic)
		}
	}
	// 3. Khởi tạo http server
	a.httpserver = server.NewHttpServer("localhost:8080", router)
}

func (a *App) makeHandler(topic string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&reqCount, 1)

		var body []byte
		if r.Body != nil {
			b, _ := io.ReadAll(r.Body)
			_ = r.Body.Close()
			body = b
		}

		// Lấy path params từ mux
		vars := mux.Vars(r)
		pathWithParams := r.URL.Path
		for k, v := range vars {
			// Có thể replace {param} trong path gốc bằng value thực
			pathWithParams = utils.ReplaceParam(pathWithParams, k, v)
		}

		replyCh := make(chan model.GatewayResult, 1)
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println("Error parsing RemoteAddr:", err)
			return
		}

		job := model.RawRequestData{
			Ctx:     r.Context(),
			Method:  r.Method,
			Path:    pathWithParams,
			Header:  r.Header.Clone(),
			Body:    body,
			IP:      ip,
			Topic:   topic,
			Token:   r.Header.Get("Authorization"),
			ReplyCh: replyCh,
		}

		select {
		case a.gmodel.RequestQueue <- job:
		case <-r.Context().Done():
			utils.WritePlainError(w, http.StatusRequestTimeout, "Client canceled")
			return
		}

		select {
		case res := <-replyCh:
			for k, vs := range res.Headers {
				for _, v := range vs {
					w.Header().Add(k, v)
				}
			}
			w.WriteHeader(res.StatusCode)
			_, _ = w.Write(res.Body)
		case <-r.Context().Done():
			utils.WritePlainError(w, http.StatusRequestTimeout, "Gateway timeout waiting for pipeline")
		}
	}
}

// ===== Pipeline =====
func (a *App) startWorkers(n int) {
	for i := 0; i < n; i++ {
		go func(id int) {
			for req := range a.gmodel.RequestQueue {
				a.process(req)
			}
		}(i)
	}
}

func (a *App) requestsMonitor(stop <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			count := atomic.SwapInt64(&reqCount, 0)
			fmt.Printf("🔥 Requests per second: %d\n", count)
		case <-stop: // stop signal received
			fmt.Println("🛑 Monitor stopped")
			return
		}
	}
}

func (a *App) process(req model.RawRequestData) {
	start := time.Now()
	requestID := utils.NewRequestID()
	//fmt.Printf("req.IP = %s\n", req.IP)
	if !a.ratelimiter.MaxReqLimiter.Allow("Max-Request-Per-Second") {
		fmt.Printf("RATE_LIMIT_MAX_REQUEST Too Many Requests\n")
		req.ReplyCh <- a.normalizedError(requestID, http.StatusTooManyRequests, "RATE_LIMIT_MAX_REQUEST", "Too Many Requests", time.Since(start))
		return
	}

	if !a.ratelimiter.IPLimiter.Allow(req.IP) {
		fmt.Printf("RATE_LIMIT_IP Too Many Requests (IP) = %s\n", req.IP)
		req.ReplyCh <- a.normalizedError(requestID, http.StatusTooManyRequests, "RATE_LIMIT_IP", "Too Many Requests (IP)", time.Since(start))
		return
	}

	// Chỉ check JWT nếu endpoint cần auth
	var claims *jwt_checker.Claims // khai báo trước
	if a.gmodel.TopicAuthMap[req.Topic] {
		var ok bool
		claims, ok = a.jwtchecker.TokenCheck(req.Token)
		if !ok {
			fmt.Printf("UNAUTHENTICATED Unauthorized (JWT) = %s\n", time.Since(start).String())
			req.ReplyCh <- a.normalizedError(requestID, http.StatusUnauthorized, "UNAUTHENTICATED", "Unauthorized (JWT)", time.Since(start))
			return
		}
	}

	userID := "anonymous" // default nếu không auth
	if claims != nil {
		userID = claims.UserID
	}
	key := userID + ":" + req.Topic
	if !a.ratelimiter.FeatureLimiter.Allow(key) {
		fmt.Println("RATE_LIMIT_FEATURE Too Many Requests (Feature) %s", key)
		req.ReplyCh <- a.normalizedError(requestID, http.StatusTooManyRequests, "RATE_LIMIT_FEATURE", "Too Many Requests (Feature)", time.Since(start))
		return
	}

	res := a.routeToInternalService(req, requestID, start)
	req.ReplyCh <- res
}

// ===== Routing =====
func (a *App) routeToInternalService(req model.RawRequestData, requestID string, start time.Time) model.GatewayResult {
	// Lấy targetURL từ topic
	targetURL := a.getTargetURL(req)
	if targetURL == "" {
		return a.normalizedError(requestID, http.StatusBadGateway, "NO_ROUTE", "No internal service for topic "+req.Topic, time.Since(start))
	}

	ctx, cancel := context.WithTimeout(req.Ctx, 3*time.Second)
	defer cancel()

	ireq, err := http.NewRequestWithContext(ctx, req.Method, targetURL, bytes.NewReader(req.Body))
	if err != nil {
		return a.normalizedError(requestID, http.StatusInternalServerError, "BUILD_REQUEST_FAILED", err.Error(), time.Since(start))
	}

	utils.CopySafeHeaders(req.Header, ireq.Header)
	ireq.Header.Set("X-Request-ID", requestID)
	ireq.Header.Set("X-Trace-ID", utils.NewRequestID())

	client := &http.Client{}
	resp, err := client.Do(ireq)
	latency := time.Since(start)
	if err != nil {
		return a.normalizedError(requestID, http.StatusBadGateway, "BAD_GATEWAY", "Internal service unreachable: "+err.Error(), latency)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// log.Println("====== Internal Service Response ======")
	// log.Printf("RequestID: %s | Status: %d %s | Latency: %dms", requestID, resp.StatusCode, http.StatusText(resp.StatusCode), latency.Milliseconds())
	// for k, v := range resp.Header {
	// 	log.Printf("%s: %v", k, v)
	// }
	// log.Println("Body:", string(respBody))
	// log.Println("======================================")

	out := map[string]interface{}{
		"request_id": requestID,
		"status":     "SUCCESS",
		"latency_ms": latency.Milliseconds(),
		"data":       json.RawMessage(respBody),
		"error":      nil,
	}

	status := http.StatusOK
	if resp.StatusCode >= 400 {
		out["status"] = "ERROR"
		out["error"] = map[string]interface{}{
			"upstream_status": resp.StatusCode,
			"message":         string(respBody),
		}
		out["data"] = nil
		status = http.StatusOK
	}

	body, _ := json.Marshal(out)
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("X-Gateway", "api-gateway")
	h.Set("X-Request-ID", requestID)

	return model.GatewayResult{
		StatusCode: status,
		Headers:    h,
		Body:       body,
	}
}

func (a *App) normalizedError(requestID string, httpCode int, code string, message string, latency time.Duration) model.GatewayResult {
	payload := map[string]interface{}{
		"request_id": requestID,
		"status":     "ERROR",
		"latency_ms": latency.Milliseconds(),
		"data":       nil,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
	body, _ := json.Marshal(payload)
	h := http.Header{}
	h.Set("Content-Type", "application/json")
	h.Set("X-Gateway", "api-gateway")
	h.Set("X-Request-ID", requestID)
	return model.GatewayResult{
		StatusCode: httpCode,
		Headers:    h,
		Body:       body,
	}
}

func (a *App) getTargetURL(req model.RawRequestData) string {
	for _, sg := range a.gmodel.ServiceGroups {
		for _, ep := range sg.Endpoints {
			topic := sg.Name + "/" + ep.Name
			if req.Topic == topic {
				return "http://" + sg.IP + ":" + strconv.Itoa(sg.Port) + req.Path
			}
		}
	}
	// fallback
	return ""
}
