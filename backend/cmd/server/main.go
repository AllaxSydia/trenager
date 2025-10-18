package main

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"log"
	"net/http"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Настройка CORS middleware
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

	// Роуты API
	http.HandleFunc("/api/tasks", corsMiddleware(handlers.TasksHandler))
	http.HandleFunc("/api/execute", corsMiddleware(handlers.ExecuteHandler))
	http.HandleFunc("/api/auth/guest", corsMiddleware(handlers.GuestAuthHandler))
	http.HandleFunc("/api/auth/register", corsMiddleware(handlers.RegisterHandler))
	http.HandleFunc("/api/auth/login", corsMiddleware(handlers.LoginHandler))

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

	log.Printf("🚀 Сервер запущен на порту %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, nil))
}
