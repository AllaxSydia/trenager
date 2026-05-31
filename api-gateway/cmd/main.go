package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"api-gateway/internal/client"
	"api-gateway/internal/config"
	"api-gateway/internal/router"
)

func main() {
	log.Println("🚪 Starting API Gateway...")

	cfg := config.LoadConfig()

	clients, err := client.NewGRPCClients(
		cfg.Services["auth"].URL,
		cfg.Services["task"].URL,
		cfg.Services["grading"].URL,
		cfg.Services["execution"].URL,
		cfg.Services["ai"].URL,
		cfg.Services["analytics"].URL,
	)
	if err != nil {
		log.Fatalf("Failed to create gRPC clients: %v", err)
	}
	defer clients.Close()

	r := router.NewRouter(clients)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("✅ API Gateway running on http://localhost:%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down API Gateway...")
	srv.Close()
	log.Println("Server stopped")
}
