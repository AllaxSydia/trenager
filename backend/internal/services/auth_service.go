package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"backend/internal/models"
)

// AuthService сервис для аутентификации
type AuthService struct {
	// Временное хранилище в памяти (позже заменим на БД)
	users  map[string]*models.User
	tokens map[string]string // token -> userID
}

func NewAuthService() *AuthService {
	return &AuthService{
		users:  make(map[string]*models.User),
		tokens: make(map[string]string),
	}
}

// GuestLogin создает гостевую сессию
func (s *AuthService) GuestLogin() (*models.User, string, error) {
	guestID := "guest-" + generateRandomID(8)

	user := &models.User{
		ID:        guestID,
		Username:  "Гость_" + generateRandomID(4),
		Email:     "guest@" + generateRandomID(6) + ".com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Сохраняем пользователя
	s.users[user.ID] = user

	// Создаем токен
	token := generateRandomID(32)
	s.tokens[token] = user.ID

	return user, token, nil
}

// ValidateToken проверяет валидность токена
func (s *AuthService) ValidateToken(token string) (*models.User, error) {
	userID, exists := s.tokens[token]
	if !exists {
		return nil, errors.New("invalid token")
	}

	user, exists := s.users[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// Register регистрирует нового пользователя (заглушка)
func (s *AuthService) Register(username, email, password string) (*models.User, string, error) {
	return nil, "", errors.New("registration temporarily disabled")
}

// Login аутентифицирует пользователя (заглушка)
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	return nil, "", errors.New("login temporarily disabled")
}

// generateRandomID генерирует случайную строку
func generateRandomID(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback на timestamp если crypto недоступен
		return string(time.Now().UnixNano())
	}
	return hex.EncodeToString(bytes)
}
