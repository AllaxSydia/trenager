package handler

import (
	"auth-service/internal/service"
	pb "auth-service/proto"
	"context"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCHandler struct {
	pb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewGRPCHandler(authService service.AuthService) *GRPCHandler {
	return &GRPCHandler{
		authService: authService,
	}
}

func (h *GRPCHandler) Register(
	ctx context.Context,
	req *pb.RegisterRequest,
) (*pb.AuthResponse, error) {
	log.Printf("gRPC Register request: username=%s, email=%s", req.Username, req.Email)

	user, tokens, err := h.authService.Register(ctx, req.Username, req.Email, req.Password)
	if err != nil {
		log.Printf("Register error: %v", err)
		return &pb.AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(codes.InvalidArgument, err.Error())
	}

	log.Printf("Register successful: user_id=%s", user.ID.String())

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
		Message: "Registration successful",
	}, nil
}

func (h *GRPCHandler) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (*pb.AuthResponse, error) {
	log.Printf("gRPC Login request: email=%s", req.Email)

	user, tokens, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("Login error: %v", err)
		return &pb.AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(codes.Unauthenticated, err.Error())
	}

	log.Printf("Login successful: user_id=%s", user.ID.String())

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

func (h *GRPCHandler) Refresh(
	ctx context.Context,
	req *pb.RefreshRequest,
) (*pb.AuthResponse, error) {
	log.Printf("gRPC Refresh request")

	tokens, err := h.authService.RefreshTokens(ctx, req.RefreshToken)
	if err != nil {
		log.Printf("Refresh error: %v", err)
		return &pb.AuthResponse{
			Success: false,
			Error:   err.Error(),
		}, status.Error(codes.Unauthenticated, err.Error())
	}

	log.Printf("Refresh successful")

	return &pb.AuthResponse{
		Success:      true,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		AccessExp:    tokens.AccessExp,
		RefreshExp:   tokens.RefreshExp,
		Message:      "Tokens refreshed successfully",
	}, nil
}

func (h *GRPCHandler) ValidateToken(
	ctx context.Context,
	req *pb.ValidateTokenRequest,
) (*pb.ValidateTokenResponse, error) {
	log.Printf("gRPC ValidateToken request")

	user, err := h.authService.ValidateAccessToken(ctx, req.AccessToken)
	if err != nil {
		log.Printf("ValidateToken error: %v", err)
		return &pb.ValidateTokenResponse{
			Valid: false,
			Error: err.Error(),
		}, nil
	}

	log.Printf("ValidateToken successful: user_id=%s", user.ID.String())

	return &pb.ValidateTokenResponse{
		Valid: true,
		User: &pb.User{
			Id:        user.ID.String(),
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (h *GRPCHandler) GetUser(
	ctx context.Context,
	req *pb.GetUserRequest,
) (*pb.UserResponse, error) {
	log.Printf("gRPC GetUser request: user_id=%s", req.UserId)

	user, err := h.authService.GetUserByID(ctx, req.UserId)
	if err != nil {
		log.Printf("GetUser error: %v", err)
		return &pb.UserResponse{
			Error: err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	if user == nil {
		log.Printf("GetUser: user not found")
		return &pb.UserResponse{
			Error: "user not found",
		}, status.Error(codes.NotFound, "user not found")
	}

	log.Printf("GetUser successful: user_id=%s", user.ID.String())

	return &pb.UserResponse{
		User: &pb.User{
			Id:        user.ID.String(),
			Username:  user.Username,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		},
	}, nil
}

func (h *GRPCHandler) HealthCheck(
	ctx context.Context,
	req *pb.HealthRequest,
) (*pb.HealthResponse, error) {
	log.Printf("gRPC Health request")

	return &pb.HealthResponse{
		Healthy:   true,
		Status:    "Auth Service is healthy and running",
		Timestamp: time.Now().Format(time.RFC3339),
	}, nil
}
