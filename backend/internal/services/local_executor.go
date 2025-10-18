package services

import (
	"backend/internal/models"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type LocalExecutor struct{}

func NewLocalExecutor() *LocalExecutor {
	return &LocalExecutor{}
}

func (l *LocalExecutor) ExecuteCode(code, language string) (*models.ExecutionResult, error) {
	if language != "go" {
		return &models.ExecutionResult{
			Success: false,
			Output:  "На Render поддерживается только Go",
		}, nil
	}

	// Создаем временную директорию
	tmpDir, err := os.MkdirTemp("", "go_exec_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Записываем код в файл
	filePath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code: %v", err)
	}

	// Выполняем код
	cmd := exec.Command("go", "run", filePath)
	output, err := cmd.CombinedOutput()

	result := &models.ExecutionResult{
		Output:  string(output),
		Success: err == nil,
	}

	if err != nil {
		result.Error = err.Error()
	}

	log.Printf("Local execution - Success: %v, Output: %s", result.Success, result.Output)
	return result, nil
}
