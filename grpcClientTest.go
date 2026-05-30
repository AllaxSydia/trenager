package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/AllaxSydia/trenager/proto/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50052",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	client := pb.NewTaskServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("=== Testing TaskService gRPC API ===\n")

	// 1. Test ListTasks
	fmt.Println("1. Calling ListTasks...")
	resp, err := client.ListTasks(ctx, &pb.ListTasksRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Success! Total tasks: %d\n", resp.Total)
	}

	// 2. Test CreateTask
	fmt.Println("\n2. Calling CreateTask...")
	createResp, err := client.CreateTask(ctx, &pb.CreateTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
		Difficulty:  "easy",
		Language:    "go",
		TestCases: []*pb.TestCase{
			{Input: "1 2", ExpectedOutput: "3"},
		},
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Success! Task ID: %s, Message: %s\n", createResp.TaskId, createResp.Message)
	}

	// 3. Test GetTask
	fmt.Println("\n3. Calling GetTask...")
	getResp, err := client.GetTask(ctx, &pb.GetTaskRequest{
		TaskId: "test-123",
	})
	if err != nil {
		fmt.Printf("   Error: %v\n", err)
	} else {
		fmt.Printf("   ✓ Success! Task: %s - %s\n", getResp.Id, getResp.Title)
	}

	fmt.Println("\n=== Test completed ===")
}
