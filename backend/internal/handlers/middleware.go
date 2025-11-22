package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware проверяет JWT токен и добавляет информацию о пользователе в контекст запроса
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, `{"error":"invalid_auth_header"}`, http.StatusUnauthorized)
			return
		}

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "please_change_this_secret"
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid_token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"invalid_token_claims"}`, http.StatusUnauthorized)
			return
		}

		// Добавляем информацию о пользователе в заголовки запроса для использования в handlers
		if userID, ok := claims["sub"].(float64); ok {
			r.Header.Set("X-User-ID", strconv.FormatInt(int64(userID), 10))
		}
		if username, ok := claims["usr"].(string); ok {
			r.Header.Set("X-Username", username)
		}
		if role, ok := claims["role"].(string); ok {
			r.Header.Set("X-User-Role", role)
		} else {
			r.Header.Set("X-User-Role", "student")
		}

		next(w, r)
	}
}

// TeacherOnlyMiddleware проверяет, что пользователь является преподавателем
func TeacherOnlyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-User-Role")
		if role != "teacher" {
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
			return
		}
		next(w, r)
	})
}

// GetUserIDFromRequest извлекает user_id из заголовков запроса
func GetUserIDFromRequest(r *http.Request) (int64, error) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0, errors.New("user_id not found in request")
	}
	// Парсим user_id из строки
	var userID int64
	_, err := fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

