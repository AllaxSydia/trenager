package main

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	cfg := config.Load()

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

	// API routes
	http.HandleFunc("/api/tasks", corsMiddleware(handlers.TasksHandler))
	http.HandleFunc("/api/execute", corsMiddleware(handlers.ExecuteHandler))
	http.HandleFunc("/api/auth/guest", corsMiddleware(handlers.GuestAuthHandler))
	http.HandleFunc("/api/auth/register", corsMiddleware(handlers.RegisterHandler))
	http.HandleFunc("/api/auth/login", corsMiddleware(handlers.LoginHandler))

	// Serve frontend with SPA routing
	http.HandleFunc("/", serveSPA)

	log.Printf("🚀 Сервер запущен на порту %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, nil))
}

func serveSPA(w http.ResponseWriter, r *http.Request) {
	// Если это API запрос - 404
	if strings.HasPrefix(r.URL.Path, "/api/") {
		http.NotFound(w, r)
		return
	}

	// Путь к статическим файлам
	staticDir := "./static"
	indexPath := filepath.Join(staticDir, "index.html")

	// Полный путь к запрашиваемому файлу
	filePath := filepath.Join(staticDir, r.URL.Path)

	// Проверяем существует ли файл
	if _, err := os.Stat(filePath); err == nil {
		// Файл существует - отдаем его
		http.ServeFile(w, r, filePath)
		return
	}

	// Для всех остальных маршрутов отдаем index.html (SPA)
	http.ServeFile(w, r, indexPath)
}
