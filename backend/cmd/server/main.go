package main

import (
	"backend/internal/config"
	"backend/internal/handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	// Обслуживание статических файлов фронтенда
	http.HandleFunc("/", serveFrontend)

	log.Printf("🚀 Сервер запущен на порту %s", cfg.Server.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, nil))
}

func serveFrontend(w http.ResponseWriter, r *http.Request) {
	// Путь к статическим файлам фронтенда
	staticDir := "./static"
	indexFile := filepath.Join(staticDir, "index.html")

	// Проверяем существование статической директории
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		// Если статики нет, отдаем JSON ответ для API запросов
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.NotFound(w, r)
			return
		}
		// Для корневого пути отдаем простое сообщение
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<!DOCTYPE html>
			<html>
			<head>
				<title>Trenager</title>
			</head>
			<body>
				<h1>🚀 Trenager Backend</h1>
				<p>Frontend static files not found. API is available at <code>/api/*</code></p>
			</body>
			</html>
		`))
		return
	}

	// Обслуживаем статические файлы
	fs := http.FileServer(http.Dir(staticDir))

	// Если файл не найден, отдаем index.html (для SPA роутинга)
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		filePath := filepath.Join(staticDir, r.URL.Path)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.ServeFile(w, r, indexFile)
			return
		}
	}

	fs.ServeHTTP(w, r)
}
