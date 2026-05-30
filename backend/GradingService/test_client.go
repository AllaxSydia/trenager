package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/AllaxSydia/trenager/proto/grading"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50053",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGradingServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("=== Testing GradingService ===\n")

	// 1. Отправка решения
	fmt.Println("1. Submitting solution...")
	resp, err := client.SubmitSolution(ctx, &pb.SubmitSolutionRequest{
		TaskId:   "task-123",
		UserId:   "user-456",
		Code:     "print('Hello World')",
		Language: "python",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("   ✓ Submission ID: %s\n", resp.SubmissionId)
	fmt.Printf("   ✓ Score: %d/%d\n", resp.Grade.Score, resp.Grade.MaxScore)
	fmt.Printf("   ✓ Status: %s\n", resp.Grade.Status)
	fmt.Printf("   ✓ Feedback: %s\n", resp.Grade.Feedback)

	// 2. Получение статистики по задаче
	fmt.Println("\n2. Getting task statistics...")
	stats, err := client.GetTaskStatistics(ctx, &pb.GetTaskStatisticsRequest{
		TaskId: "task-123",
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("   ✓ Total submissions: %d\n", stats.TotalSubmissions)
	fmt.Printf("   ✓ Success rate: %.1f%%\n",
		float64(stats.SuccessfulSubmissions)/float64(stats.TotalSubmissions)*100)
	fmt.Printf("   ✓ Average score: %.1f\n", stats.AverageScore)

	fmt.Println("\n✅ All tests passed!")
}
