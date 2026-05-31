package router

import (
	"net/http"

	"api-gateway/internal/client"
	"api-gateway/internal/handler"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func NewRouter(clients *client.GRPCClients) http.Handler {
	r := mux.NewRouter()

	authHandler := handler.NewAuthHandler(clients)
	taskHandler := handler.NewTaskHandler(clients)
	executionHandler := handler.NewExecutionHandler(clients)

	// CORS
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")

	// Task routes
	api.HandleFunc("/tasks", taskHandler.ListTasks).Methods("GET")
	api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")

	// Execution routes
	api.HandleFunc("/execute", executionHandler.ExecuteCode).Methods("POST")
	api.HandleFunc("/execute/test", executionHandler.ExecuteTest).Methods("POST")

	// Health
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}).Methods("GET")

	return cors(r)
}
