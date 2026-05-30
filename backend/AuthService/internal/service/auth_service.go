package service

import (
	"context"
	"errors"
	"log"
	"time"

	"auth-service/internal/models"
	"auth-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repository.UserRepository
	jwtSecret []byte
}

func NewAuthService(repo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		repo:      repo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*models.User, *models.TokenPair, error) {
	log.Printf("Registering user: %s", email)

	// Проверяем, существует ли пользователь
	existingUser, _ := s.repo.GetByEmail(ctx, email)
	if existingUser != nil {
		return nil, nil, errors.New("user already exists")
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, errors.New("failed to hash password")
	}

	user := &models.User{
		ID:           uuid.New(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	// Генерируем токены
	tokens, err := s.generateTokens(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return nil, nil, err
	}

	// Сохраняем refresh token
	if err := s.repo.SaveRefreshToken(ctx, user.ID, tokens.RefreshToken, tokens.RefreshExp); err != nil {
		log.Printf("Warning: failed to save refresh token: %v", err)
	}

	return user, tokens, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, *models.TokenPair, error) {
	log.Printf("Login attempt: %s", email)

	// Находим пользователя
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Обновляем время последнего входа
	s.repo.UpdateLastLogin(ctx, user.ID)

	// Генерируем токены
	tokens, err := s.generateTokens(user.ID.String(), user.Email, user.Role)
	if err != nil {
		return nil, nil, err
	}

	// Сохраняем refresh token
	if err := s.repo.SaveRefreshToken(ctx, user.ID, tokens.RefreshToken, tokens.RefreshExp); err != nil {
		log.Printf("Warning: failed to save refresh token: %v", err)
	}

	return user, tokens, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid claims")
	}

	return &models.Claims{
		UserID: claims.Subject,
		Email:  claims.Issuer,
	}, nil
}

func (s *AuthService) generateTokens(userID, email, role string) (*models.TokenPair, error) {
	accessExp := time.Now().Add(15 * time.Minute)
	refreshExp := time.Now().Add(7 * 24 * time.Hour)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    email,
		ExpiresAt: jwt.NewNumericDate(accessExp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(refreshExp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &models.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		AccessExp:    accessExp.Unix(),
		RefreshExp:   refreshExp.Unix(),
	}, nil
}
