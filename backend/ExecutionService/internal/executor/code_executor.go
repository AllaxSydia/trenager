package executor

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"execution-service/internal/models"

	"github.com/google/uuid"
)

type CodeExecutor struct {
	tempDir string
}

func NewCodeExecutor() *CodeExecutor {
	return &CodeExecutor{
		tempDir: "./temp",
	}
}

// Execute выполняет код в изолированной среде
func (e *CodeExecutor) Execute(ctx context.Context, req *models.ExecutionRequest) (*models.ExecutionResult, error) {
	executionID := uuid.New().String()
	workDir := filepath.Join(e.tempDir, executionID)

	// Создаем временную директорию
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(workDir)

	startTime := time.Now()

	// Создаем файл с кодом
	codeFile, err := e.createCodeFile(workDir, req.Code, req.Language)
	if err != nil {
		return nil, err
	}

	// Запускаем выполнение
	output, err := e.runCode(ctx, workDir, codeFile, req.Language, req.Input, req.TimeLimitMs)

	executionTime := time.Since(startTime).Milliseconds()

	result := &models.ExecutionResult{
		ID:              executionID,
		Success:         err == nil,
		Output:          output,
		ExecutionTimeMs: executionTime,
		MemoryUsedBytes: 0,
		Status:          "completed",
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		result.Status = "runtime_error"
	}

	return result, nil
}

func (e *CodeExecutor) createCodeFile(workDir, code, language string) (string, error) {
	var filename string
	switch language {
	case "python":
		filename = "main.py"
	case "go":
		filename = "main.go"
	case "javascript", "js":
		filename = "main.js"
	case "java":
		filename = "Main.java"
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	filePath := filepath.Join(workDir, filename)

	// Для Go нужно обернуть код в main функцию если её нет
	if language == "go" && !strings.Contains(code, "package main") {
		code = fmt.Sprintf("package main\n\nimport \"fmt\"\n\nfunc main() {\n    %s\n}", code)
	}

	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return "", fmt.Errorf("failed to write code file: %w", err)
	}

	return filePath, nil
}

func (e *CodeExecutor) runCode(ctx context.Context, workDir, codeFile, language, input string, timeLimitMs int) (string, error) {
	var cmd *exec.Cmd

	switch language {
	case "python":
		cmd = exec.CommandContext(ctx, "python3", codeFile)
	case "go":
		// Сначала компилируем
		outputFile := filepath.Join(workDir, "main")
		compileCmd := exec.CommandContext(ctx, "go", "build", "-o", outputFile, codeFile)
		if compileErr := compileCmd.Run(); compileErr != nil {
			return "", fmt.Errorf("compilation error: %w", compileErr)
		}
		cmd = exec.CommandContext(ctx, outputFile)
	case "javascript", "js":
		cmd = exec.CommandContext(ctx, "node", codeFile)
	case "java":
		// Компилируем Java
		compileCmd := exec.CommandContext(ctx, "javac", codeFile)
		if compileErr := compileCmd.Run(); compileErr != nil {
			return "", fmt.Errorf("compilation error: %w", compileErr)
		}
		cmd = exec.CommandContext(ctx, "java", "-cp", workDir, "Main")
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	// Устанавливаем таймаут
	if timeLimitMs > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeLimitMs)*time.Millisecond)
		defer cancel()
		cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
	}

	// Передаем stdin
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n"
		}
		output += stderr.String()
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return output, fmt.Errorf("execution timeout (%d ms)", timeLimitMs)
		}
		return output, err
	}

	return output, nil
}

// ExecuteTests выполняет тесты
func (e *CodeExecutor) ExecuteTests(ctx context.Context, code string, language string, testCases []models.TestCase) ([]models.TestResult, error) {
	var results []models.TestResult

	for i, test := range testCases {
		req := &models.ExecutionRequest{
			Code:     code,
			Language: language,
			Input:    test.Input,
		}

		result, err := e.Execute(ctx, req)

		testResult := models.TestResult{
			ID:              uuid.New().String(),
			TestID:          fmt.Sprintf("test_%d", i+1),
			Passed:          false,
			ExecutionTimeMs: result.ExecutionTimeMs,
		}

		if err != nil {
			testResult.Error = err.Error()
		} else if result.Success {
			testResult.ActualOutput = result.Output
			testResult.ExpectedOutput = test.ExpectedOutput
			testResult.Passed = strings.TrimSpace(result.Output) == strings.TrimSpace(test.ExpectedOutput)
		} else {
			testResult.Error = result.Error
		}

		results = append(results, testResult)
	}

	return results, nil
}
