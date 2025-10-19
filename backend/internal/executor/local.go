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
)

type LocalExecutor struct{}

func NewLocalExecutor() *LocalExecutor {
	return &LocalExecutor{}
}

func (e *LocalExecutor) Execute(code, language string) (map[string]interface{}, error) {
	switch strings.ToLower(language) {
	case "go":
		return e.executeGo(code)
	case "python", "python3":
		return e.executePython(code)
	case "javascript", "node":
		return e.executeJavaScript(code)
	case "cpp", "c++":
		return e.executeCpp(code)
	case "java":
		return e.executeJava(code)
	default:
		return map[string]interface{}{
			"output":   "",
			"error":    fmt.Sprintf("Unsupported language: %s. Supported: Go, Python, JavaScript, C++, Java", language),
			"exitCode": 1,
		}, nil
	}
}

func (e *LocalExecutor) executeGo(code string) (map[string]interface{}, error) {
	// Создаем временную директорию
	tmpDir, err := os.MkdirTemp("", "go_exec_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Записываем код в файл
	mainFile := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(mainFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code: %v", err)
	}

	// Выполняем с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", mainFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return map[string]interface{}{
			"output":   "",
			"error":    "Execution timeout (30 seconds exceeded)",
			"exitCode": 1,
		}, nil
	}

	if err != nil {
		return map[string]interface{}{
			"output":   stdout.String(),
			"error":    stderr.String(),
			"exitCode": 1,
		}, nil
	}

	return map[string]interface{}{
		"output":   stdout.String(),
		"error":    "",
		"exitCode": 0,
	}, nil
}

func (e *LocalExecutor) executePython(code string) (map[string]interface{}, error) {
	// Создаем временный файл
	tmpFile, err := os.CreateTemp("", "python_*.py")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(code)); err != nil {
		return nil, fmt.Errorf("failed to write code: %v", err)
	}
	tmpFile.Close()

	// Выполняем с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "python3", tmpFile.Name())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return map[string]interface{}{
			"output":   "",
			"error":    "Execution timeout (30 seconds exceeded)",
			"exitCode": 1,
		}, nil
	}

	if err != nil {
		return map[string]interface{}{
			"output":   stdout.String(),
			"error":    stderr.String(),
			"exitCode": 1,
		}, nil
	}

	return map[string]interface{}{
		"output":   stdout.String(),
		"error":    "",
		"exitCode": 0,
	}, nil
}

func (e *LocalExecutor) executeJavaScript(code string) (map[string]interface{}, error) {
	// Создаем временный файл
	tmpFile, err := os.CreateTemp("", "javascript_*.js")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(code)); err != nil {
		return nil, fmt.Errorf("failed to write code: %v", err)
	}
	tmpFile.Close()

	// Выполняем с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "node", tmpFile.Name())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return map[string]interface{}{
			"output":   "",
			"error":    "Execution timeout (30 seconds exceeded)",
			"exitCode": 1,
		}, nil
	}

	if err != nil {
		return map[string]interface{}{
			"output":   stdout.String(),
			"error":    stderr.String(),
			"exitCode": 1,
		}, nil
	}

	return map[string]interface{}{
		"output":   stdout.String(),
		"error":    "",
		"exitCode": 0,
	}, nil
}

func (e *LocalExecutor) executeCpp(code string) (map[string]interface{}, error) {
	// Создаем временную директорию
	tmpDir, err := os.MkdirTemp("", "cpp_exec_*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Записываем код в файл
	sourceFile := filepath.Join(tmpDir, "main.cpp")
	if err := os.WriteFile(sourceFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code: %v", err)
	}

	// Компилируем
	executable := filepath.Join(tmpDir, "main")
	compileCmd := exec.Command("g++", "-o", executable, sourceFile)
	var compileStderr bytes.Buffer
	compileCmd.Stderr = &compileStderr

	if err := compileCmd.Run(); err != nil {
		return map[string]interface{}{
			"output":   "",
			"error":    "Compilation failed: " + compileStderr.String(),
			"exitCode": 1,
		}, nil
	}

	// Выполняем с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, executable)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		return map[string]interface{}{
			"output":   "",
			"error":    "Execution timeout (30 seconds exceeded)",
			"exitCode": 1,
		}, nil
	}

	if err != nil {
		return map[string]interface{}{
			"output":   stdout.String(),
			"error":    stderr.String(),
			"exitCode": 1,
		}, nil
	}

	return map[string]interface{}{
		"output":   stdout.String(),
		"error":    "",
		"exitCode": 0,
	}, nil
}

func (e *LocalExecutor) executeJava(code string) (map[string]interface{}, error) {
	// Для Java возвращаем сообщение что компилятор не установлен
	return map[string]interface{}{
		"output":   "",
		"error":    "Java execution is not available in this environment. Please use Go, Python, JavaScript or C++.",
		"exitCode": 1,
	}, nil
}
