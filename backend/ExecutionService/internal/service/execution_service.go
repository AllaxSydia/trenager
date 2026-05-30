package service

import (
	"context"
	"fmt"
	"log"

	"execution-service/internal/executor"
	"execution-service/internal/models"

	"github.com/google/uuid"
)

type ExecutionService struct {
	executor *executor.CodeExecutor
}

func NewExecutionService() *ExecutionService {
	return &ExecutionService{
		executor: executor.NewCodeExecutor(),
	}
}

// ExecuteCode - выполняет код пользователя
func (s *ExecutionService) ExecuteCode(ctx context.Context, req *models.ExecutionRequest) (*models.ExecutionResult, error) {
	log.Printf("Executing code: language=%s, code_length=%d", req.Language, len(req.Code))

	// Выполняем код
	result, err := s.executor.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("execution failed: %w", err)
	}

	return result, nil
}

// ExecuteTests - выполняет тесты для проверки решения
func (s *ExecutionService) ExecuteTests(ctx context.Context, code string, language string, testCases []models.TestCase) ([]models.TestResult, error) {
	log.Printf("Running %d tests for language %s", len(testCases), language)

	var results []models.TestResult

	for i, test := range testCases {
		testResult := models.TestResult{
			ID:     uuid.New().String(),
			TestID: fmt.Sprintf("test_%d", i+1),
		}

		// Выполняем код с тестовым вводом
		req := &models.ExecutionRequest{
			Code:     code,
			Language: language,
			Input:    test.Input,
		}

		execResult, err := s.executor.Execute(ctx, req)
		if err != nil {
			testResult.Passed = false
			testResult.Error = err.Error()
		} else {
			testResult.ActualOutput = execResult.Output
			testResult.ExpectedOutput = test.ExpectedOutput
			testResult.Passed = execResult.Output == test.ExpectedOutput
			testResult.ExecutionTimeMs = execResult.ExecutionTimeMs
		}

		results = append(results, testResult)
	}

	return results, nil
}

// GetStatus - получает статус выполнения
func (s *ExecutionService) GetStatus(ctx context.Context, executionID string) (*models.ExecutionStatus, error) {
	// Здесь будет логика получения статуса из хранилища
	return &models.ExecutionStatus{
		ID:     executionID,
		Status: "completed",
	}, nil
}
