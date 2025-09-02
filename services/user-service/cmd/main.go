package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"userservice/internal/app"
)

func main() {
	// 1. Create the API Gateway app
	apiApp := app.NewAuthServiceApp()

	// 2. Start the app (starts workers + HTTP server)
	go apiApp.Start()
	log.Println("🚀 API Gateway is running...")

	// 3. Graceful shutdown on SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("⚠️ Shutting down API Gateway...")
	apiApp.Stop()
}
