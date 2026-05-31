package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/AllaxSydia/trenager/proto/execution"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type ExecutionServer struct {
	pb.UnimplementedExecutionServiceServer
}

// ExecuteCode - выполняет код пользователя
func (s *ExecutionServer) ExecuteCode(ctx context.Context, req *pb.ExecuteCodeRequest) (*pb.ExecuteCodeResponse, error) {
	log.Printf("💻 ExecuteCode: language=%s, code_length=%d", req.Language, len(req.Code))

	// Валидация
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}

	if req.Language == "" {
		return nil, status.Error(codes.InvalidArgument, "language is required")
	}

	// Имитация выполнения кода
	startTime := time.Now()

	// Здесь будет реальное выполнение кода через Docker
	// Пока возвращаем mock результат

	executionTime := time.Since(startTime).Milliseconds()

	return &pb.ExecuteCodeResponse{
		Success:         true,
		Output:          generateMockOutput(req.Code, req.Language),
		Error:           "",
		ExecutionTimeMs: executionTime,
		MemoryUsedBytes: 1024 * 1024, // 1 MB
		Status:          "completed",
	}, nil
}

// ExecuteTest - выполняет тесты
func (s *ExecutionServer) ExecuteTest(ctx context.Context, req *pb.ExecuteTestRequest) (*pb.ExecuteTestResponse, error) {
	log.Printf("🧪 ExecuteTest: language=%s, tests_count=%d", req.Language, len(req.Tests))

	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "code is required")
	}

	var results []*pb.TestResult
	passedCount := 0

	for i, test := range req.Tests {
		// Имитация выполнения теста
		passed := i%2 == 0 // Каждый второй тест проходит

		result := &pb.TestResult{
			TestId:          string(rune(i + 1)),
			Passed:          passed,
			ActualOutput:    generateMockOutputForTest(req.Code, test.Input),
			Error:           "",
			ExecutionTimeMs: 10,
		}

		if passed {
			passedCount++
		}

		results = append(results, result)
	}

	return &pb.ExecuteTestResponse{
		AllPassed:   passedCount == len(req.Tests),
		Results:     results,
		PassedCount: int32(passedCount),
		TotalCount:  int32(len(req.Tests)),
	}, nil
}

// GetExecutionStatus - получает статус выполнения
func (s *ExecutionServer) GetExecutionStatus(ctx context.Context, req *pb.GetExecutionStatusRequest) (*pb.ExecutionStatus, error) {
	log.Printf("📊 GetExecutionStatus: id=%s", req.ExecutionId)

	return &pb.ExecutionStatus{
		ExecutionId: req.ExecutionId,
		Status:      "completed",
		Result:      "Execution completed successfully",
		Error:       "",
	}, nil
}

func generateMockOutput(code, language string) string {
	outputs := map[string]string{
		"python":     "Hello from Python!\nOutput: 42",
		"go":         "Hello from Go!\nOutput: 42\n",
		"javascript": "Hello from JavaScript!\nOutput: 42\n",
		"java":       "Hello from Java!\nOutput: 42\n",
	}

	if output, ok := outputs[language]; ok {
		return output
	}

	return "Code executed successfully\nOutput: 42\n"
}

func generateMockOutputForTest(code, input string) string {
	if input == "" {
		return "42\n"
	}
	return "Processed: " + input + "\n"
}

func main() {
	log.Println("🚀 Starting ExecutionService gRPC server...")

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterExecutionServiceServer(grpcServer, &ExecutionServer{})

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50054"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Printf("✅ ExecutionService running on port %s", port)
		log.Println("⚠️  Running in MOCK mode - code is not actually executed")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
