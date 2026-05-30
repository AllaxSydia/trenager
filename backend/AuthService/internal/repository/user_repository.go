package repository

import (
	"context"
	"database/sql"
	"fmt"

	"auth-service/internal/models"
	"auth-service/pkg/database"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *database.Database
}

func NewUserRepository(db *database.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (id, username, email, password_hash, role, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := r.db.GetDB().ExecContext(ctx, query,
		user.ID, user.Username, user.Email, user.PasswordHash, user.Role,
		user.CreatedAt, user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at 
            FROM users WHERE email = $1`

	user := &models.User{}
	err := r.db.GetDB().QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, role, created_at, updated_at 
            FROM users WHERE id = $1`

	user := &models.User{}
	err := r.db.GetDB().QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (r *UserRepository) SaveRefreshToken(ctx context.Context, userID uuid.UUID, token string, expiresAt int64) error {
	query := `INSERT INTO refresh_tokens (id, user_id, token, expires_at) VALUES ($1, $2, $3, $4)`

	_, err := r.db.GetDB().ExecContext(ctx, query, uuid.New(), userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

func (r *UserRepository) ValidateRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
	query := `SELECT user_id FROM refresh_tokens WHERE token = $1 AND revoked = false AND expires_at > NOW()`

	var userID uuid.UUID
	err := r.db.GetDB().QueryRowContext(ctx, query, token).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, fmt.Errorf("invalid or expired refresh token")
		}
		return uuid.Nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	return userID, nil
}

func (r *UserRepository) RevokeRefreshToken(ctx context.Context, token string) error {
	query := `UPDATE refresh_tokens SET revoked = true WHERE token = $1`

	_, err := r.db.GetDB().ExecContext(ctx, query, token)
	if err != nil {
		return fmt.Errorf("failed to revoke refresh token: %w", err)
	}

	return nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET last_login = NOW() WHERE id = $1`

	_, err := r.db.GetDB().ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}

	return nil
}
