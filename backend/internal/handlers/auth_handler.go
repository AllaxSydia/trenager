package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"net/http"
)

var authService = services.NewAuthService()

// GuestAuthHandler обработчик гостевого доступа
func GuestAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем и GET и POST для простоты тестирования
	if r.Method != "POST" && r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, token, err := authService.GuestLogin()
	if err != nil {
		http.Error(w, `{"success": false, "message": "Failed to create guest session"}`, http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		Success: true,
		Message: "Добро пожаловать в гостевом режиме!",
		User:    user,
		Token:   token,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterHandler обработчик регистрации
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success": false, "message": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	response := models.AuthResponse{
		Success: false,
		Message: "Регистрация временно отключена. Используйте гостевой режим.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LoginHandler обработчик входа
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success": false, "message": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	response := models.AuthResponse{
		Success: false,
		Message: "Вход временно отключен. Используйте гостевой режим.",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
