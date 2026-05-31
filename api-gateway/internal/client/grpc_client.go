package client

import (
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	aipb "github.com/AllaxSydia/trenager/proto/ai"
	analyticspb "github.com/AllaxSydia/trenager/proto/analytics"
	authpb "github.com/AllaxSydia/trenager/proto/auth"
	executionpb "github.com/AllaxSydia/trenager/proto/execution"
	gradingpb "github.com/AllaxSydia/trenager/proto/grading"
	taskpb "github.com/AllaxSydia/trenager/proto/task"
)

type GRPCClients struct {
	Auth      authpb.AuthServiceClient
	Task      taskpb.TaskServiceClient
	Grading   gradingpb.GradingServiceClient
	Execution executionpb.ExecutionServiceClient
	AI        aipb.AIServiceClient
	Analytics analyticspb.AnalyticsServiceClient

	conns []*grpc.ClientConn
}

func NewGRPCClients(authURL, taskURL, gradingURL, executionURL, aiURL, analyticsURL string) (*GRPCClients, error) {
	clients := &GRPCClients{}

	// Auth Service
	authConn, err := grpc.NewClient(authURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %w", err)
	}
	clients.Auth = authpb.NewAuthServiceClient(authConn)
	clients.conns = append(clients.conns, authConn)
	log.Printf("Connected to AuthService at %s", authURL)

	// Task Service
	taskConn, err := grpc.NewClient(taskURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to task service: %w", err)
	}
	clients.Task = taskpb.NewTaskServiceClient(taskConn)
	clients.conns = append(clients.conns, taskConn)
	log.Printf("Connected to TaskService at %s", taskURL)

	// Grading Service
	gradingConn, err := grpc.NewClient(gradingURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grading service: %w", err)
	}
	clients.Grading = gradingpb.NewGradingServiceClient(gradingConn)
	clients.conns = append(clients.conns, gradingConn)
	log.Printf("Connected to GradingService at %s", gradingURL)

	// Execution Service
	executionConn, err := grpc.NewClient(executionURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to execution service: %w", err)
	}
	clients.Execution = executionpb.NewExecutionServiceClient(executionConn)
	clients.conns = append(clients.conns, executionConn)
	log.Printf("Connected to ExecutionService at %s", executionURL)

	// AI Service
	aiConn, err := grpc.NewClient(aiURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to AI service: %w", err)
	}
	clients.AI = aipb.NewAIServiceClient(aiConn)
	clients.conns = append(clients.conns, aiConn)
	log.Printf("Connected to AIService at %s", aiURL)

	// Analytics Service
	analyticsConn, err := grpc.NewClient(analyticsURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to analytics service: %w", err)
	}
	clients.Analytics = analyticspb.NewAnalyticsServiceClient(analyticsConn)
	clients.conns = append(clients.conns, analyticsConn)
	log.Printf("Connected to AnalyticsService at %s", analyticsURL)

	return clients, nil
}

func (c *GRPCClients) Close() {
	for _, conn := range c.conns {
		conn.Close()
	}
}
