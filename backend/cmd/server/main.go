package main

import (
	"backend/internal/database"
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
	database.Init()
	defer database.Close()

	port := getPort()

	log.Printf("🚀 Starting server on port %s", port)
	log.Printf("📁 Current directory: %s", getCurrentDir())
	log.Printf("🌐 Environment: %s", getEnvironment())

	// Проверяем существование статики
	if _, err := os.Stat("./static"); err != nil {
		log.Printf("💡 Running in API-only mode")
	} else {
		log.Println("✅ Static directory found")
	}

	// CORS middleware
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Разрешаем все локальные адреса для разработки
			allowedOrigins := getAllowedOrigins()
			origin := r.Header.Get("Origin")

			// Проверяем разрешенные origins
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-API-Key")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Security headers
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

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
	http.HandleFunc("/api/ai/review", loggingMiddleware(corsMiddleware(handlers.AIReviewHandler)))
	http.HandleFunc("/api/auth/guest", loggingMiddleware(corsMiddleware(handlers.GuestAuthHandler)))
	http.HandleFunc("/api/auth/register", loggingMiddleware(corsMiddleware(handlers.RegisterHandler)))

	// Admin routes (только для преподавателей)
	http.HandleFunc("/api/admin/statistics", loggingMiddleware(corsMiddleware(handlers.TeacherOnlyMiddleware(handlers.StatisticsHandler))))

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

	// Health check
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		healthStatus := "healthy"
		checks := map[string]interface{}{
			"api":         "ok",
			"environment": getEnvironment(),
			"timestamp":   time.Now().Format(time.RFC3339),
			"compilers":   []string{"python", "node", "g++", "javac"},
		}

		response := map[string]interface{}{
			"status":  healthStatus,
			"checks":  checks,
			"version": "1.0.0",
			"uptime":  time.Since(startTime).String(),
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
			"compilers":   []string{"python", "node", "g++", "javac"},
		}
		json.NewEncoder(w).Encode(response)
	})))

	// Task endpoint
	http.HandleFunc("/api/task/", loggingMiddleware(corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

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

	// Запускаем сервер
	server := &http.Server{
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}

var startTime = time.Now()

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

func getEnvironment() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	return env
}

func getAllowedOrigins() []string {
	env := getEnvironment()

	if env == "production" {
		return []string{
			"https://your-production-domain.com", // Замени на свой домен
		}
	}

	// Для development разрешаем все локальные адреса
	return []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:3001",
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:8080",
		"http://127.0.0.1:8080",
	}
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
