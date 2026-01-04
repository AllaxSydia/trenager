package handlers

import (
	"backend/internal/database"
	"backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Request / Response structs
type registerRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Success  bool   `json:"success"`
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

// RegisterHandler - регистрирует пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "invalid_request",
		})
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Username == "" || req.Email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "missing_fields",
		})
		return
	}

	// Хэшируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("bcrypt error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "server_error",
		})
		return
	}

	// Вставляем в БД (по умолчанию роль 'student')
	query := `INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, 'student') RETURNING id, created_at`
	var id int64
	var createdAt time.Time
	err = database.DB.QueryRow(query, req.Username, req.Email, string(hash)).Scan(&id, &createdAt)
	if err != nil {
		log.Printf("db insert error: %v", err)
		// Возможно уникальность нарушена
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "user_exists",
		})
		return
	}

	// Генерируем токен сразу после регистрации
	token, err := generateToken(id, req.Username, req.Email, "student")
	if err != nil {
		log.Printf("token error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "server_error",
		})
		return
	}

	res := authResponse{
		Success:  true,
		Token:    token,
		Username: req.Username,
		Email:    req.Email,
		Role:     "student",
		Message:  "Registration successful",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// LoginHandler - аутентификация
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "invalid_request",
		})
		return
	}

	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "missing_fields",
		})
		return
	}

	// Тестовые пользователи для быстрого входа (если нет в БД)
	if email == "teacher@mail.com" && req.Password == "123456789" {
		// Автоматический вход для тестового учителя
		token, err := generateToken(1, "teacher_avg", "teacher@mail.com", "teacher")
		if err != nil {
			log.Printf("token error for test teacher: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(authResponse{
				Success: false,
				Error:   "server_error",
			})
			return
		}

		res := authResponse{
			Success:  true,
			Token:    token,
			Username: "teacher_avg",
			Email:    "teacher@mail.com",
			Role:     "teacher",
			Message:  "Login successful (test teacher)",
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}

	if email == "student@trenager.ru" && req.Password == "123456789" {
		// Автоматический вход для тестового студента
		token, err := generateToken(2, "student_ivan", "student@trenager.ru", "student")
		if err != nil {
			log.Printf("token error for test student: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(authResponse{
				Success: false,
				Error:   "server_error",
			})
			return
		}

		res := authResponse{
			Success:  true,
			Token:    token,
			Username: "student_ivan",
			Email:    "student@trenager.ru",
			Role:     "student",
			Message:  "Login successful (test student)",
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}

	user, err := findUserByEmail(email)
	if err != nil {
		log.Printf("⚠️ User not found for email: %s, error: %v", email, err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "invalid_credentials",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Printf("⚠️ Invalid password for email: %s", email)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "invalid_credentials",
		})
		return
	}

	log.Printf("✅ Login successful for user: %s (role: %s)", user.Username, user.Role)

	token, err := generateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "server_error",
		})
		return
	}

	// Убеждаемся, что роль установлена (по умолчанию 'student')
	if user.Role == "" {
		user.Role = "student"
	}

	res := authResponse{
		Success:  true,
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Message:  "Login successful",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// helper: findUserByEmail
func findUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, COALESCE(role, 'student'), created_at, COALESCE(updated_at, created_at) FROM users WHERE email = $1 LIMIT 1`
	row := database.DB.QueryRow(query, email)

	var u models.User
	if err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
		log.Printf("⚠️ Error scanning user: %v", err)
		return nil, err
	}
	return &u, nil
}

// JWT handling
func generateToken(userID int64, username, email, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "please_change_this_secret" // ПОМЕНЯЙ на проде
	}
	if role == "" {
		role = "student"
	}
	claims := jwt.MapClaims{
		"sub":   userID,
		"usr":   username,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Middleware example (если потребуется позже)
func ParseTokenFromRequest(r *http.Request) (jwt.MapClaims, error) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return nil, errors.New("no auth")
	}
	parts := strings.Fields(auth)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return nil, errors.New("invalid auth header")
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "please_change_this_secret"
	}
	tokenStr := parts[1]
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// проверяем алгоритм
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}

// GuestAuthHandler — временная авторизация без регистрации
func GuestAuthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	guestID := time.Now().Unix()
	username := fmt.Sprintf("Guest-%d", guestID)

	// Создаем токен как у обычного пользователя
	token, err := generateToken(guestID, username, fmt.Sprintf("%s@guest.local", username), "student")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "server_error",
		})
		return
	}

	res := authResponse{
		Success:  true,
		Token:    token,
		Username: username,
		Email:    fmt.Sprintf("%s@guest.local", username),
		Role:     "student",
		Message:  "Guest login successful",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// QuickLoginHandler - быстрый вход тестовыми пользователями
func QuickLoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		UserType string `json:"user_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "invalid_request",
		})
		return
	}

	var email, username, role string
	var userID int64

	switch req.UserType {
	case "teacher":
		email = "teacher@mail.com"
		username = "teacher_avg"
		role = "teacher"
		userID = 1
	case "student":
		email = "student@trenager.ru"
		username = "student_ivan"
		role = "student"
		userID = 2
	case "admin":
		email = "admin@trenager.ru"
		username = "admin_root"
		role = "admin"
		userID = 3
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "invalid_user_type",
		})
		return
	}

	token, err := generateToken(userID, username, email, role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "server_error",
		})
		return
	}

	res := authResponse{
		Success:  true,
		Token:    token,
		Username: username,
		Email:    email,
		Role:     role,
		Message:  fmt.Sprintf("Quick login as %s successful", role),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// ValidateTokenHandler - проверка валидности токена
func ValidateTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "no_token",
		})
		return
	}

	claims, err := ParseTokenFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "invalid_token",
		})
		return
	}

	// Извлекаем данные пользователя из токена
	var userID int64
	var username, email, role string

	if sub, ok := claims["sub"].(float64); ok {
		userID = int64(sub)
		// Можно залогировать: log.Printf("Token validated for user ID: %d", userID)
	}
	if usr, ok := claims["usr"].(string); ok {
		username = usr
	}
	if eml, ok := claims["email"].(string); ok {
		email = eml
	}
	if rl, ok := claims["role"].(string); ok {
		role = rl
	}

	// Пример использования userID для дополнительной проверки
	if userID <= 0 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "invalid_user_id",
		})
		return
	}

	res := authResponse{
		Success:  true,
		Username: username,
		Email:    email,
		Role:     role,
		Message:  fmt.Sprintf("Token is valid for user ID: %d", userID),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// GetUserInfoHandler - получение информации о пользователе по токену
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims, err := ParseTokenFromRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "invalid_token",
		})
		return
	}

	// Получаем email из токена
	email, ok := claims["email"].(string)
	if !ok || email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "invalid_token_data",
		})
		return
	}

	// Ищем пользователя в БД
	user, err := findUserByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(authResponse{
			Success: false,
			Error:   "user_not_found",
		})
		return
	}

	res := authResponse{
		Success:  true,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Message:  "User info retrieved",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
