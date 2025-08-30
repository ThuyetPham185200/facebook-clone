package app

import (
	"log"
	"userservice/internal/api"
	"userservice/internal/core/http-server/server"
	"userservice/internal/infra/store"

	"github.com/gorilla/mux"
)

type App struct {
	httpserver *server.HttpServer
	userapi    *api.UserAPI
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
	us := store.NewUserStore(
		"localhost", // IP
		"5432",      // Port
		"taopq",     // user_name
		"123456a@",  // password
		"mydb",      // db
	)
	a.userapi = api.NewUserAPI(us)
	router := mux.NewRouter()
	a.userapi.RegisterRoutes(router)
	a.httpserver = server.NewHttpServer("localhost:9001", router)
}
