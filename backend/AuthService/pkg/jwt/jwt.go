package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	accessSecret  string
	refreshSecret string
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessExp    int64  `json:"access_exp"`
	RefreshExp   int64  `json:"refresh_exp"`
}

func NewManager(accessSecret, refreshSecret string) *Manager {
	if refreshSecret == "" {
		refreshSecret = accessSecret + "-refresh" // Можно использовать один секрет с суффиксом
	}

	return &Manager{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (m *Manager) GenerateTokenPair(userID, role string) (*TokenPair, error) {
	// Access Token (30 минут)
	accessExp := time.Now().Add(30 * time.Minute)
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"type":    "access",
		"exp":     accessExp.Unix(),
		"iat":     time.Now().Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(m.accessSecret))
	if err != nil {
		return nil, err
	}

	// Refresh Token (7 дней)
	refreshExp := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     refreshExp.Unix(),
		"iat":     time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(m.refreshSecret))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		AccessExp:    accessExp.Unix(),
		RefreshExp:   refreshExp.Unix(),
	}, nil
}

// ValidateAccessToken проверяет access токен
func (m *Manager) ValidateAccessToken(tokenString string) (jwt.MapClaims, error) {
	return m.validateToken(tokenString, m.accessSecret, "access")
}

// ValidateRefreshToken проверяет refresh токен
func (m *Manager) ValidateRefreshToken(tokenString string) (jwt.MapClaims, error) {
	return m.validateToken(tokenString, m.refreshSecret, "refresh")
}

// RefreshTokens обновляет пару токенов
func (m *Manager) RefreshTokens(refreshToken string) (*TokenPair, error) {
	claims, err := m.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, nil
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid refresh token: no user_id")
	}

	role := "student"
	if r, ok := claims["role"].(string); ok {
		role = r
	}

	return m.GenerateTokenPair(userID, role)
}

// Вспомогательный метод для валидации
func (m *Manager) validateToken(tokenString, secret, expectedType string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Проверяем тип токена
	if tokenType, ok := claims["type"].(string); !ok || tokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}

func (m *Manager) ExtractUserId(tokenString string) (string, error) {
	claims, err := m.ValidateAccessToken(tokenString)
	if err != nil {
		return "", nil
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token")
	}

	return userID, nil
}
