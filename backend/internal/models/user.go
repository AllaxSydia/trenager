package models

import "time"

type User struct {
	ID           int64     `json:"id"`             // int64, чтобы совпадало с базой
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`              // не возвращаем в JSON
	Role         string    `json:"role"`            // 'student' или 'teacher'
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuthRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
	User    *User  `json:"user,omitempty"`
}
