package app

import (
	"authservice/internal/api"
	auth "authservice/internal/core/authentication"
	"authservice/internal/core/http-server/server"
	"authservice/internal/core/session"
	"authservice/internal/core/userserviceclient"
	"authservice/internal/infra/store"
	"log"
	"time"

	"github.com/gorilla/mux"
)

type App struct {
	httpserver     *server.HttpServer
	authapi        *api.AuthAPI
	authentication auth.AuthenticationManager
	session        session.SessionManager
}

func NewAuthServiceApp() *App {
	var app = &App{}
	app.init()
	return app
}

func (a *App) Start() {
	if err := a.httpserver.Start(); err != nil {
		log.Fatalf("❌ Failed to start: %v", err)
	}
}

func (a *App) Stop() {
	if err := a.httpserver.Stop(); err != nil {
		log.Printf("⚠️ Error stopping server: %v", err)
	}
	log.Println("✅ Server stopped gracefully")
}

// ///////////////////////////////////////////////////////////////////////////////////////
func (a *App) init() {
	router := mux.NewRouter()
	cfg := &session.JwtConfig{
		AccessSecret:  []byte("my-access-secret"),
		RefreshSecret: []byte("my-refresh-secret"),
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
	}

	dbcredentalscfg := &store.PostGresConfig{
		Host:     "localhost", // IP
		Port:     "5432",      // Port
		User:     "taopq",     // user_name
		Password: "123456a@",  // password
		DBname:   "mydb",      // db
	}

	a.session = session.NewSessionManager(cfg, store.NewPostgresTokenStore(dbcredentalscfg))

	redisstorecfg := &store.RedisConfig{
		Host:     "localhost",
		Port:     "6379",
		Password: "",
		DBNumber: 1,
	}
	a.authentication = auth.NewAuthenticationManager(
		store.NewCredentialsStore(dbcredentalscfg, redisstorecfg),
		userserviceclient.NewUserServiceClient("http://localhost:9001"),
		a.session)
	a.authapi = api.NewAuthAPI(a.authentication, a.session)
	a.authapi.RegisterRoutes(router)
	// 3. Khởi tạo http server
	a.httpserver = server.NewHttpServer("localhost:9000", router)
}
