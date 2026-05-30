package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/AllaxSydia/trenager/proto/ai"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type AIServer struct {
	pb.UnimplementedAIServiceServer
}

func (s *AIServer) GetHint(ctx context.Context, req *pb.GetHintRequest) (*pb.GetHintResponse, error) {
	log.Printf("🤖 GetHint: task=%s, level=%d", req.TaskId, req.HintLevel)
	hints := map[int32]string{
		1: "Попробуйте разбить задачу на подзадачи.",
		2: "Обратите внимание на входные и выходные данные.",
		3: "Возможно, стоит использовать цикл или рекурсию.",
		4: "Попробуйте использовать дополнительную структуру данных.",
		5: "Вот пример: пройдитесь по элементам и накопите результат.",
	}
	hint := hints[req.HintLevel]
	if hint == "" {
		hint = hints[3]
	}
	return &pb.GetHintResponse{
		Hint:        hint,
		HintLevel:   req.HintLevel,
		ExampleCode: "def solution(data):\n    return result",
	}, nil
}

func (s *AIServer) ReviewCode(ctx context.Context, req *pb.ReviewCodeRequest) (*pb.ReviewCodeResponse, error) {
	log.Printf("🔍 ReviewCode: language=%s", req.Language)
	return &pb.ReviewCodeResponse{
		QualityScore:    85,
		Issues:          []string{"Несколько длинных функций", "Отсутствуют комментарии"},
		Suggestions:     []string{"Разбейте код на функции", "Добавьте документацию"},
		OverallFeedback: "Хорошая работа! Код решает задачу, но есть место для улучшения.",
	}, nil
}

func (s *AIServer) GetRecommendations(ctx context.Context, req *pb.GetRecommendationsRequest) (*pb.GetRecommendationsResponse, error) {
	log.Printf("📚 GetRecommendations: user=%s", req.UserId)
	return &pb.GetRecommendationsResponse{
		Recommendations: []*pb.TaskRecommendation{
			{TaskId: "rec-001", Title: "Сортировка массива", Reason: "Популярная задача", DifficultyScore: 0.3},
			{TaskId: "rec-002", Title: "Динамическое программирование", Reason: "Следующий уровень", DifficultyScore: 0.7},
		},
	}, nil
}

func (s *AIServer) AskQuestion(ctx context.Context, req *pb.AskQuestionRequest) (*pb.AskQuestionResponse, error) {
	log.Printf("💬 AskQuestion: %s", req.Question[:min(50, len(req.Question))])
	return &pb.AskQuestionResponse{
		Answer:       "Попробуйте использовать цикл for для прохода по элементам.",
		HasCode:      true,
		CodeExamples: []string{"for i in range(len(arr)):\n    print(arr[i])"},
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	log.Println("🤖 Starting AIService gRPC server...")

	grpcServer := grpc.NewServer()
	pb.RegisterAIServiceServer(grpcServer, &AIServer{})
	reflection.Register(grpcServer) // Добавляем reflection

	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		log.Println("✅ AIService running on :50055")
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
