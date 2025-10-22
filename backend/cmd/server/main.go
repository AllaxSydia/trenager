package main

import (
	"backend/internal/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	// Получаем порт из переменной окружения или используем 8080
	port := getPort() // для Render

	log.Printf("🚀 Starting server on port %s", port)
	log.Printf("📁 Current directory: %s", getCurrentDir())

	// Если не static то всё равно запуститься backend, только front не будет
	// Проверяем существование статики
	if _, err := os.Stat("./static"); err != nil {
		log.Printf("⚠️ Static directory not found: %v", err)
	} else {
		log.Println("✅ Static directory found")

		// Логируем содержимое static
		files, _ := os.ReadDir("./static")
		log.Printf("📂 Static files count: %d", len(files))
		for _, file := range files {
			log.Printf("   - %s", file.Name())
		}
	}

	// Упрощенный CORS middleware
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")

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

	// Serve frontend static files - ОБНОВЛЕНО ДЛЯ ЕДИНОГО КОНТЕЙНЕРА
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Fallback route for SPA - ОБНОВЛЕНО
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	log.Printf("✅ Server ready to accept requests on port %s", port)
	log.Printf("🌐 Frontend will be served from /static")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return dir
}
