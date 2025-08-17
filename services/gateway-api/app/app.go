package app

import (
	"context"
	"gateway-api/internal/api"
	"gateway-api/internal/config"
	"gateway-api/internal/middleware"
	"gateway-api/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	httpServer *server.HttpServer
}

func NewApp(cfg *config.Config) *App {
	// Middleware init
	jwtChecker := middleware.NewJWTChecker(cfg.JWTSecret)
	rateLimiter := middleware.NewRateLimiter(cfg.RateLimit)

	// Router + API handler
	router := api.NewRouter(jwtChecker, rateLimiter)

	// Server
	httpSrv := server.NewHttpServer(cfg.HTTPAddr, router)

	return &App{httpServer: httpSrv}
}

func (a *App) Run() error {
	// Start server
	if err := a.httpServer.Start(); err != nil {
		return err
	}

	// Graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Stop(ctx); err != nil {
		log.Printf("❌ Failed to stop: %v", err)
	}
	log.Println("✅ Server stopped gracefully")
	return nil
}
