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
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role,omitempty"`
	Error    string `json:"error,omitempty"`
}

// RegisterHandler - регистрирует пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_request"}`, http.StatusBadRequest)
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"missing_fields"}`, http.StatusBadRequest)
		return
	}

	// Хэшируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("bcrypt error: %v", err)
		http.Error(w, `{"error":"server_error"}`, http.StatusInternalServerError)
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
		http.Error(w, `{"error":"user_exists"}`, http.StatusConflict)
		return
	}

	// Генерируем токен сразу после регистрации
	token, err := generateToken(id, req.Username, req.Email, "student")
	if err != nil {
		log.Printf("token error: %v", err)
		http.Error(w, `{"error":"server_error"}`, http.StatusInternalServerError)
		return
	}

	res := authResponse{
		Token:    token,
		Username: req.Username,
		Email:    req.Email,
		Role:     "student",
	}

	json.NewEncoder(w).Encode(res)
}

// LoginHandler - аутентификация
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_request"}`, http.StatusBadRequest)
		return
	}
	email := strings.TrimSpace(strings.ToLower(req.Email))
	if email == "" || req.Password == "" {
		http.Error(w, `{"error":"missing_fields"}`, http.StatusBadRequest)
		return
	}

	user, err := findUserByEmail(email)
	if err != nil {
		log.Printf("⚠️ User not found for email: %s, error: %v", email, err)
		http.Error(w, `{"error":"invalid_credentials"}`, http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Printf("⚠️ Invalid password for email: %s", email)
		http.Error(w, `{"error":"invalid_credentials"}`, http.StatusUnauthorized)
		return
	}
	
	log.Printf("✅ Login successful for user: %s (role: %s)", user.Username, user.Role)

	token, err := generateToken(user.ID, user.Username, user.Email, user.Role)
	if err != nil {
		http.Error(w, `{"error":"server_error"}`, http.StatusInternalServerError)
		return
	}

	// Убеждаемся, что роль установлена (по умолчанию 'student')
	if user.Role == "" {
		user.Role = "student"
	}

	res := authResponse{
		Token:    token,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
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
		http.Error(w, `{"error":"server_error"}`, http.StatusInternalServerError)
		return
	}

	res := map[string]interface{}{
		"token":    token,
		"username": username,
		"email":    fmt.Sprintf("%s@guest.local", username),
		"guest":    true,
	}

	json.NewEncoder(w).Encode(res)
}