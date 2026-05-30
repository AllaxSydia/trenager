package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/AllaxSydia/trenager/proto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	client  pb.AuthServiceClient
	conn    *grpc.ClientConn
	mu      sync.RWMutex
	addr    string
	timeout time.Duration
}

type Config struct {
	Addr    string
	Timeout time.Duration
}

// NewAuthClient создает новый клиент для AuthService
func NewAuthClient(cfg *Config) (*AuthClient, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}

	conn, err := grpc.NewClient(cfg.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(cfg.Timeout),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}

	return &AuthClient{
		client:  pb.NewAuthServiceClient(conn),
		conn:    conn,
		addr:    cfg.Addr,
		timeout: cfg.Timeout,
	}, nil
}

// Close закрывает соединение
func (c *AuthClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ValidateToken проверяет валидность токена и возвращает пользователя
func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &pb.ValidateTokenRequest{
		AccessToken: token,
	}

	resp, err := c.client.ValidateToken(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %w", err)
	}

	if !resp.Valid {
		return nil, fmt.Errorf("invalid token: %s", resp.Error)
	}

	return resp.User, nil
}

// GetUser получает информацию о пользователе
func (c *AuthClient) GetUser(ctx context.Context, userID string) (*pb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req := &pb.GetUserRequest{
		UserId: userID,
	}

	resp, err := c.client.GetUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if resp.Error != "" {
		return nil, fmt.Errorf("user not found: %s", resp.Error)
	}

	return resp.User, nil
}

// HealthCheck проверяет доступность AuthService
func (c *AuthClient) HealthCheck(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	req := &pb.HealthRequest{}
	_, err := c.client.HealthCheck(ctx, req)
	if err != nil {
		return fmt.Errorf("auth service health check failed: %w", err)
	}

	return nil
}

// Interceptor для проверки аутентификации в gRPC
func (c *AuthClient) AuthInterceptor(ctx context.Context, method string, req interface{},
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	// Получаем токен из метаданных
	token, err := extractTokenFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("authentication required: %w", err)
	}

	// Валидируем токен
	user, err := c.ValidateToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Добавляем пользователя в контекст
	ctx = context.WithValue(ctx, "user", user)
	ctx = context.WithValue(ctx, "user_id", user.Id)

	return handler(ctx, req)
}

func extractTokenFromContext(ctx context.Context) (string, error) {
	// Здесь логика извлечения токена из gRPC метаданных
	// Временно возвращаем заглушку
	return "", fmt.Errorf("token extraction not implemented yet")
}

// AuthRequired проверяет, требуется ли аутентификация для метода
func AuthRequired(method string) bool {
	// Методы, не требующие аутентификации
	publicMethods := map[string]bool{
		"/task.TaskService/GetTask":   false, // Публичные задачи может смотреть любой
		"/task.TaskService/ListTasks": false,
	}

	if required, exists := publicMethods[method]; exists {
		return required
	}
	return true // По умолчанию требуем аутентификацию
}
