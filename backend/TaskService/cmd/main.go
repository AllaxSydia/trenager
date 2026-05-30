package main

import (
	"log"
	"net"

	"TaskService/internal/handler"
	"TaskService/internal/repository"
	"TaskService/internal/service"
	"TaskService/pkg/database"

	pb "github.com/AllaxSydia/trenager/proto/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := database.LoadConfig()
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskGrpcHandler(taskService)

	grpcServer := grpc.NewServer()
	pb.RegisterTaskServiceServer(grpcServer, taskHandler)
	reflection.Register(grpcServer) // Добавлено

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("✅ TaskService running on :50052 with PostgreSQL")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
