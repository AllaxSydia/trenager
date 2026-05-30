package handler

import (
	"context"
	"log"
	"time"

	"auth-service/internal/service"

	pb "github.com/AllaxSydia/trenager/proto/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	log.Printf("📝 Register: username=%s, email=%s", req.Username, req.Email)

	user, tokens, err := h.authService.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		return &pb.AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.AuthResponse{
		Success:      true,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		AccessExp:    tokens.AccessExp,
		RefreshExp:   tokens.RefreshExp,
		User: &pb.User{
			Id:        user.ID.String(),
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		Message: "User registered successfully",
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	log.Printf("🔐 Login: email=%s", req.Email)

	user, tokens, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return &pb.AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.AuthResponse{
		Success:      true,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		AccessExp:    tokens.AccessExp,
		RefreshExp:   tokens.RefreshExp,
		User: &pb.User{
			Id:        user.ID.String(),
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
		Message: "Login successful",
	}, nil
}

func (h *AuthHandler) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.AuthResponse, error) {
	log.Printf("🔄 Refresh token")

	// TODO: Implement refresh logic
	return &pb.AuthResponse{
		Success: false,
		Error:   "refresh not implemented yet",
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	log.Printf("✅ Validate token")

	claims, err := h.authService.ValidateToken(ctx, req.AccessToken)
	if err != nil {
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid: true,
		User: &pb.User{
			Id:    claims.UserID,
			Email: claims.Email,
		},
	}, nil
}

func (h *AuthHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	log.Printf("👤 GetUser: user_id=%s", req.UserId)

	// TODO: Implement get user from database
	return &pb.UserResponse{
		User: &pb.User{
			Id:    req.UserId,
			Email: "user@example.com",
		},
	}, nil
}

func (h *AuthHandler) HealthCheck(ctx context.Context, req *pb.HealthRequest) (*pb.HealthResponse, error) {
	log.Printf("💚 HealthCheck")

	return &pb.HealthResponse{
		Healthy:   true,
		Status:    "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}
