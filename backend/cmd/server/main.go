package main

import (
	"log"
	"net/http"
)

func main() {
	// Получаем порт от Render
	port := getPort()

	log.Printf("🚀 Запуск сервера на порту %s", port)

	// CORS middleware
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next(w, r)
		}
	}

	// Test endpoint
	http.HandleFunc("/api/test", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok", "message": "API is working"}`))
	}))

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "healthy"}`))
	})

	// Serve frontend static files
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Printf("✅ Сервер готов принимать запросы на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getPort() string {
	port := "10000"
	if envPort := "10000"; envPort != "" {
		port = envPort
	}
	return port
}
