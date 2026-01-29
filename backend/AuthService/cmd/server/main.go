package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/repository"
	"auth-service/internal/service"
	"auth-service/pkg/database"
	pb "auth-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Загружаем конфигурацию
	cfg := config.Load()

	log.Printf("🚀 Starting Auth Service")
	log.Printf("📝 Configuration:")
	log.Printf("   Port: %d", cfg.Port)
	log.Printf("   Database: %s:%d/%s", cfg.DBHost, cfg.DBPort, cfg.DBName)
	log.Printf("   Environment: %s", cfg.Environment)

	// Подключаемся к базе данных
	dbConfig := database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("⚠️ Error closing database: %v", err)
		}
		log.Println("✅ Database connection closed")
	}()

	log.Println("✅ Connected to PostgreSQL database")

	// Проверяем соединение с базой
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
	}
	log.Println("✅ Database ping successful")

	// Создаем репозиторий
	userRepo := repository.NewPostgresUserRepository(db)
	log.Println("✅ User repository initialized")

	// Создаем сервис
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.RefreshSecret)
	log.Println("✅ Auth service initialized")

	// Создаем gRPC обработчик
	grpcHandler := handler.NewGRPCHandler(authService)
	log.Println("✅ gRPC handler initialized")

	// Настраиваем gRPC сервер с middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingInterceptor,
			recoveryInterceptor,
		),
	)

	// Регистрируем сервис
	pb.RegisterAuthServiceServer(grpcServer, grpcHandler)

	// Включаем reflection для отладки (только в development)
	if cfg.Environment == "development" {
		reflection.Register(grpcServer)
		log.Println("✅ gRPC reflection enabled (development mode)")
	}

	// Запускаем gRPC сервер
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("❌ Failed to create listener: %v", err)
	}

	// Канал для graceful shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем сервер в горутине
	go func() {
		log.Printf("✅ Auth Service gRPC server listening on port %d", cfg.Port)
		log.Printf("📡 Available endpoints:")
		log.Printf("   - /auth.AuthService/Register")
		log.Printf("   - /auth.AuthService/Login")
		log.Printf("   - /auth.AuthService/Refresh")
		log.Printf("   - /auth.AuthService/ValidateToken")
		log.Printf("   - /auth.AuthService/GetUser")
		log.Printf("   - /auth.AuthService/HealthCheck")

		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("❌ Failed to serve gRPC: %v", err)
		}
	}()

	// Ждем сигнал завершения
	<-stopChan
	log.Println("🛑 Received shutdown signal")

	// Graceful shutdown
	log.Println("⏳ Shutting down gracefully...")

	// Останавливаем gRPC сервер
	grpcServer.GracefulStop()
	log.Println("✅ gRPC server stopped")

	log.Println("👋 Auth Service shutdown complete")
}

// Middleware для логирования
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// Пропускаем health checks из логов
	if info.FullMethod != "/auth.AuthService/HealthCheck" {
		log.Printf("📥 gRPC call: %s", info.FullMethod)
	}

	resp, err := handler(ctx, req)

	if info.FullMethod != "/auth.AuthService/HealthCheck" {
		duration := time.Since(start)
		log.Printf("📤 gRPC call %s completed in %v", info.FullMethod, duration)
	}

	return resp, err
}

// Middleware для recovery (обработка паник)
func recoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("⚠️ PANIC recovered in gRPC handler %s: %v", info.FullMethod, r)
			err = fmt.Errorf("internal server error")
		}
	}()

	return handler(ctx, req)
}
