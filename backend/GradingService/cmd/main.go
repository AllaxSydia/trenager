package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"grading-service/internal/handler"
	"grading-service/internal/repository"
	"grading-service/internal/service"
	"grading-service/pkg/database"

	pb "github.com/AllaxSydia/trenager/proto/grading"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	log.Println("📊 Starting GradingService gRPC server...")

	cfg := database.LoadConfig()
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	gradingRepo := repository.NewGradingRepository(db)
	gradingService := service.NewGradingService()
	gradingService.SetRepository(gradingRepo)
	gradingHandler := handler.NewGradingHandler(gradingService)

	grpcServer := grpc.NewServer()
	pb.RegisterGradingServiceServer(grpcServer, gradingHandler)
	reflection.Register(grpcServer) // Добавляем reflection

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Println("✅ GradingService running on :50053 with PostgreSQL")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
