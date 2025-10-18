package models

import "time"

// User представляет модель пользователя
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserProgress представляет прогресс пользователя
type UserProgress struct {
	UserID      string    `json:"user_id"`
	TaskID      string    `json:"task_id"`
	Completed   bool      `json:"completed"`
	Attempts    int       `json:"attempts"`
	BestScore   float64   `json:"best_score"`
	LastAttempt time.Time `json:"last_attempt"`
}

// AuthRequest представляет запрос аутентификации
type AuthRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse представляет ответ аутентификации
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}
