package executor

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"execution-service/internal/models"
)

type CodeExecutor struct {
	// Здесь будут настройки для изолированного выполнения
}

func NewCodeExecutor() *CodeExecutor {
	return &CodeExecutor{}
}

// Execute - выполняет код (mock версия)
func (e *CodeExecutor) Execute(ctx context.Context, req *models.ExecutionRequest) (*models.ExecutionResult, error) {
	// Имитируем выполнение кода
	startTime := time.Now()

	// Симулируем задержку выполнения
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}

	executionTime := time.Since(startTime).Milliseconds()

	// Mock результат
	result := &models.ExecutionResult{
		ID:              fmt.Sprintf("exec_%d", time.Now().UnixNano()),
		Success:         true,
		Output:          fmt.Sprintf("Mock output for %s code", req.Language),
		Error:           "",
		ExecutionTimeMs: executionTime,
		MemoryUsedBytes: rand.Int63n(50 * 1024 * 1024), // 0-50 MB
		Status:          "completed",
	}

	// Симулируем ошибку если код содержит "error"
	if len(req.Code) > 0 && req.Code[:5] == "error" {
		result.Success = false
		result.Error = "Runtime error: something went wrong"
		result.Status = "runtime_error"
	}

	// Симулируем таймаут
	if req.TimeLimitMs > 0 && executionTime > int64(req.TimeLimitMs) {
		result.Success = false
		result.Error = "Execution timeout"
		result.Status = "timeout"
	}

	return result, nil
}
