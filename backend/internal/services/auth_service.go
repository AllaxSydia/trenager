package services

import (
	"database/sql"
	"errors"
	"os"
	"time"

	"backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Register(username, email, password string) (*models.User, string, error) {
	// Проверка на существующего пользователя
	var exists bool
	err := s.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE email=$1 OR username=$2)",
		email, username,
	).Scan(&exists)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", errors.New("user already exists with given email or username")
	}

	// Хеширование пароля
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	now := time.Now().UTC()
	var id int64

	err = s.db.QueryRow(
		`INSERT INTO users (username, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		username, email, string(hash), now, now,
	).Scan(&id)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		ID:        id,
		Username:  username,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	token, err := generateJWT(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(emailOrUsername, password string) (*models.User, string, error) {
	var id, username, email, passwordHash string
	var createdAt, updatedAt time.Time

	err := s.db.QueryRow(
		`SELECT id, username, email, password_hash, created_at, updated_at
		FROM users WHERE email=$1 OR username=$1`,
		emailOrUsername,
	).Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", errors.New("invalid credentials")
		}
		return nil, "", err
	}

	// Сравнение пароля
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	user := &models.User{
		ID:        id,
		Username:  username,
		Email:     email,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	token, err := generateJWT(user.ID, user.Username)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func generateJWT(userID int64, username string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret-change-me-in-production"
	}

	claims := jwt.MapClaims{
		"sub":  userID,
		"name": username,
		"iat":  time.Now().Unix(),
		"exp":  time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
