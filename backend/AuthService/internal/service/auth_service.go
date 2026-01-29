package service

import (
	"auth-service/internal/models"
	"auth-service/internal/repository"
	"auth-service/pkg/jwt"
	"context"
	"errors"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, username, email, password string) (*models.User, *jwt.TokenPair, error)
	Login(ctx context.Context, email, password string) (*models.User, *jwt.TokenPair, error)
	ValidateAccessToken(ctx context.Context, token string) (*models.User, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*jwt.TokenPair, error)
	GetUserByID(ctx context.Context, userID string) (*models.User, error)
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	userRepo repository.UserRepository
	jwt      *jwt.Manager
}

func NewAuthService(
	userRepo repository.UserRepository,
	accessSecret string,
	refreshSecret string,
) AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      jwt.NewManager(accessSecret, refreshSecret),
	}
}

func (s *authService) Register(
	ctx context.Context,
	username string,
	email string,
	password string,
) (*models.User, *jwt.TokenPair, error) {
	// Проверяем существование пользователя
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if existingUser != nil {
		return nil, nil, errors.New("user already exists")
	}

	// Создаем пользователя
	user, err := models.NewUser(username, email, password)
	if err != nil {
		return nil, nil, err
	}

	// Сохраняем в БД
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, nil, err
	}

	// Генерируем пару токенов
	tokenPair, err := s.jwt.GenerateTokenPair(user.ID.String(), user.Role)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}

func (s *authService) Login(
	ctx context.Context,
	email,
	password string,
) (*models.User, *jwt.TokenPair, error) {
	// Находим пользователя по email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}
	if user == nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Проверяем пароль
	if !user.CheckPassword(password) {
		return nil, nil, errors.New("invalid credentials")
	}

	// Генерируем пару токенов
	tokenPair, err := s.jwt.GenerateTokenPair(user.ID.String(), user.Role)
	if err != nil {
		return nil, nil, err
	}

	return user, tokenPair, nil
}

func (s *authService) ValidateAccessToken(
	ctx context.Context,
	token string,
) (*models.User, error) {
	claims, err := s.jwt.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("Invalid token claims")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(ctx, userID)
}

func (s *authService) RefreshTokens(
	ctx context.Context,
	refreshToken string,
) (*jwt.TokenPair, error) {
	return s.jwt.RefreshTokens(refreshToken)
}

func (s *authService) GetUserByID(
	ctx context.Context,
	userID string,
) (*models.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	return s.userRepo.FindByID(ctx, id)
}

func (s *authService) Logout(
	ctx context.Context,
	refreshToken string,
) error {
	return nil
}
