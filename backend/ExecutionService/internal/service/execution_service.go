package service

import (
	"context"
	"log"

	"execution-service/internal/executor"
	"execution-service/internal/models"
)

type ExecutionService struct {
	executor *executor.CodeExecutor
}

func NewExecutionService() *ExecutionService {
	return &ExecutionService{
		executor: executor.NewCodeExecutor(),
	}
}

func (s *ExecutionService) ExecuteCode(ctx context.Context, req *models.ExecutionRequest) (*models.ExecutionResult, error) {
	log.Printf("Executing code: language=%s, code_length=%d", req.Language, len(req.Code))
	return s.executor.Execute(ctx, req)
}

func (s *ExecutionService) ExecuteTests(ctx context.Context, code string, language string, testCases []models.TestCase) ([]models.TestResult, error) {
	log.Printf("Running %d tests for language %s", len(testCases), language)
	return s.executor.ExecuteTests(ctx, code, language, testCases)
}

func (s *ExecutionService) GetStatus(ctx context.Context, executionID string) (*models.ExecutionStatus, error) {
	return &models.ExecutionStatus{
		ID:     executionID,
		Status: "completed",
	}, nil
}
