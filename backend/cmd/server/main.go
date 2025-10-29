package main

import (
	"backend/internal/handlers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	// Получаем порт из переменной окружения (Railway автоматически устанавливает PORT)
	port := getPort()

	log.Printf("🚀 Starting server on port %s", port)
	log.Printf("📁 Current directory: %s", getCurrentDir())
	log.Printf("🌐 Environment: %s", getEnvironment())

	// Проверяем существование статики (для монолитного деплоя)
	if _, err := os.Stat("./static"); err != nil {
		log.Printf("⚠️ Static directory not found: %v", err)
		log.Printf("💡 Running in API-only mode")
	} else {
		log.Println("✅ Static directory found")
		files, _ := os.ReadDir("./static")
		log.Printf("📂 Static files count: %d", len(files))
	}

	// Production CORS middleware
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Для production разрешаем запросы с любых источников
			// В будущем можно ограничить конкретными доменами
			origin := r.Header.Get("Origin")
			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-API-Key")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			// Добавляем security headers для production
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			next(w, r)
		}
	}

	// Request logging middleware
	loggingMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Пропускаем health checks в логах чтобы не засорять
			if r.URL.Path != "/health" && r.URL.Path != "/api/health" {
				log.Printf("📥 %s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
			}

			next(w, r)

			if r.URL.Path != "/health" && r.URL.Path != "/api/health" {
				log.Printf("📤 %s %s completed in %v", r.Method, r.URL.Path, time.Since(start))
			}
		}
	}

	// API Routes with CORS and logging
	http.HandleFunc("/api/tasks", loggingMiddleware(corsMiddleware(handlers.TasksHandler)))
	http.HandleFunc("/api/check", loggingMiddleware(corsMiddleware(handlers.CheckHandler)))
	http.HandleFunc("/api/execute", loggingMiddleware(corsMiddleware(handlers.ExecuteHandler)))
	http.HandleFunc("/api/auth/login", loggingMiddleware(corsMiddleware(handlers.LoginHandler)))
	http.HandleFunc("/api/auth/guest", loggingMiddleware(corsMiddleware(handlers.GuestAuthHandler)))
	http.HandleFunc("/api/auth/register", loggingMiddleware(corsMiddleware(handlers.RegisterHandler)))

	// Test endpoint
	http.HandleFunc("/api/test", loggingMiddleware(corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"status":      "ok",
			"message":     "API is working",
			"timestamp":   time.Now().Format(time.RFC3339),
			"version":     "1.0.0",
			"environment": getEnvironment(),
		}
		json.NewEncoder(w).Encode(response)
	})))

	// Health check (без CORS для load balancers)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Проверяем критичные компоненты
		healthStatus := "healthy"
		checks := map[string]string{
			"api":       "ok",
			"memory":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		}

		response := map[string]interface{}{
			"status":  healthStatus,
			"checks":  checks,
			"version": "1.0.0",
		}

		json.NewEncoder(w).Encode(response)
	})

	// API Health check
	http.HandleFunc("/api/health", loggingMiddleware(corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]interface{}{
			"status":      "api_healthy",
			"timestamp":   time.Now().Format(time.RFC3339),
			"environment": getEnvironment(),
			"port":        port,
			"version":     "1.0.0",
		}
		json.NewEncoder(w).Encode(response)
	})))

	// Task endpoint
	http.HandleFunc("/api/task/", loggingMiddleware(corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Парсим параметры из URL
		path := strings.TrimPrefix(r.URL.Path, "/api/task/")
		parts := strings.Split(path, "/")

		if len(parts) < 3 {
			http.Error(w, `{"error": "Invalid task path. Use /api/task/lang/topic/id"}`, http.StatusBadRequest)
			return
		}

		lang := parts[0]
		topic := parts[1]
		taskId := parts[2]

		// Валидация языка
		supportedLanguages := map[string]bool{
			"python":     true,
			"javascript": true,
			"cpp":        true,
			"java":       true,
		}

		if !supportedLanguages[lang] {
			http.Error(w, `{"error": "Unsupported language. Use: python, javascript, cpp, java"}`, http.StatusBadRequest)
			return
		}

		task := map[string]interface{}{
			"id":          taskId,
			"title":       fmt.Sprintf("Task %s in %s", taskId, lang),
			"description": fmt.Sprintf("Write a %s program for %s topic", lang, topic),
			"language":    lang,
			"topic":       topic,
			"difficulty":  "beginner",
			"defaultCode": getDefaultCode(lang),
			"supported":   true,
			"environment": getEnvironment(),
		}

		json.NewEncoder(w).Encode(task)
	})))

	// Serve frontend static files (если есть)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// Fallback route for SPA
	http.HandleFunc("/index.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// 404 handler for API routes
	http.HandleFunc("/api/", loggingMiddleware(corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error":         "API endpoint not found",
			"path":          r.URL.Path,
			"timestamp":     time.Now().Format(time.RFC3339),
			"documentation": "Available endpoints: /api/execute, /api/check, /api/task/:lang/:topic/:id",
		})
	})))

	log.Printf("✅ Server ready to accept requests on port %s", port)
	log.Printf("🌐 Environment: %s", getEnvironment())
	log.Printf("📡 Available endpoints:")
	log.Printf("   GET  /health")
	log.Printf("   GET  /api/health")
	log.Printf("   POST /api/execute")
	log.Printf("   POST /api/check")
	log.Printf("   GET  /api/task/:lang/:topic/:id")

	if _, err := os.Stat("./static"); err == nil {
		log.Printf("   GET  / (frontend)")
	}

	// Запускаем сервер
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("ℹ️  PORT environment variable not set, using default: %s", port)
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

func getEnvironment() string {
	env := os.Getenv("RAILWAY_ENVIRONMENT")
	if env == "" {
		env = os.Getenv("ENVIRONMENT")
		if env == "" {
			env = "development"
		}
	}
	return env
}

func getDefaultCode(lang string) string {
	switch lang {
	case "python":
		return "# Write your Python code here\nprint(\"Hello World\")"
	case "javascript":
		return "// Write your JavaScript code here\nconsole.log(\"Hello World\")"
	case "cpp":
		return `// Write your C++ code here
#include <iostream>
using namespace std;

int main() {
    std::cout << "Hello World" << std::endl;
    return 0;
}`
	case "java":
		return `// Write your Java code here
public class Main {
    public static void main(String[] args) {
        System.out.println("Hello World");
    }
}`
	default:
		return "// Write your code here"
	}
}
