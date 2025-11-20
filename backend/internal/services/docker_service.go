package services

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"backend/internal/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// Docker файл сервиса. Изолятор выполнения кода

// Обертка над Docker API через официальный Go клиент
// Единственное поле - клиент для взаимодействия с Docker демоном
type DockerService struct {
	client *client.Client
}

func NewDockerService() (*DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Printf("⚠️ Docker client creation failed: %v", err)
		return nil, fmt.Errorf("Docker not available: %w", err)
	}

	// Проверяем подключение к Docker
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = cli.Ping(ctx)
	if err != nil {
		log.Printf("⚠️ Docker not available: %v", err)
		return nil, fmt.Errorf("Docker not available: %w", err)
	}

	log.Println("✅ Docker service initialized successfully")
	return &DockerService{client: cli}, nil
}

// LanguageConfigs конфигурация для разных языков программирования
var LanguageConfigs = map[string]models.LanguageConfig{
	"python": {
		DockerImage: "python:3.9-alpine",
		RunCmd:      []string{"python", "/app/code.py"},
		FileName:    "code.py",
		Timeout:     10 * time.Second,
	},
	"javascript": {
		DockerImage: "node:18-alpine",
		RunCmd:      []string{"node", "/app/code.js"},
		FileName:    "code.js",
		Timeout:     10 * time.Second,
	},
	"java": {
		DockerImage: "openjdk:17-alpine",
		CompileCmd:  []string{"javac", "/app/Main.java"}, // Компилируем Main.java
		RunCmd:      []string{"java", "Main"},            // Запускаем класс Main
		FileName:    "Main.java",                         // Файл должен называться Main.java
		Timeout:     15 * time.Second,
	},
	"cpp": {
		DockerImage: "gcc:latest",
		CompileCmd:  []string{"g++", "-o", "/app/code", "/app/code.cpp"},
		RunCmd:      []string{"/app/code"},
		FileName:    "code.cpp",
		Timeout:     15 * time.Second,
	},
	"go": {
		DockerImage: "golang:1.19-alpine",
		RunCmd:      []string{"go", "run", "/app/code.go"},
		FileName:    "code.go",
		Timeout:     10 * time.Second,
	},
}

func (s *DockerService) ExecuteCode(code, language string) (*models.ExecutionResult, error) {
	config, exists := LanguageConfigs[language]
	if !exists {
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	log.Printf("🔄 Executing %s code: %s", language, code)

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	// Создаем временный файл с кодом
	tempDir, err := os.MkdirTemp("", "code-execution")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Записываем код в файл
	filePath := filepath.Join(tempDir, config.FileName)
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code to file: %w", err)
	}

	log.Printf("📁 Code written to: %s", filePath)

	// Создаем контейнер
	containerID, err := s.createContainer(ctx, tempDir, config)
	if err != nil {
		log.Printf("❌ Failed to create container: %v", err)
		return nil, fmt.Errorf("failed to create container: %w", err)
	}
	defer s.removeContainer(ctx, containerID)

	log.Printf("🐳 Container created: %s", containerID)

	// Запускаем контейнер
	if err := s.startContainer(ctx, containerID); err != nil {
		log.Printf("❌ Failed to start container: %v", err)
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	log.Printf("🚀 Container started: %s", containerID)

	// Ждем завершения и получаем результат
	result, err := s.waitForCompletion(ctx, containerID, config)
	if err != nil {
		log.Printf("❌ Failed to wait for completion: %v", err)
		return nil, fmt.Errorf("failed to wait for completion: %w", err)
	}

	log.Printf("✅ Execution result: success=%v, output=%s", result.Success, result.Output)

	return result, nil
}

func (s *DockerService) createContainer(ctx context.Context, codePath string, config models.LanguageConfig) (string, error) {
	// Подготавливаем команды
	cmd := config.RunCmd
	if len(config.CompileCmd) > 0 {
		// Если нужна компиляция, объединяем команды
		compileCmd := strings.Join(config.CompileCmd, " ")
		runCmd := strings.Join(config.RunCmd, " ")
		cmd = []string{"/bin/sh", "-c", fmt.Sprintf("%s && %s", compileCmd, runCmd)}
	}

	resp, err := s.client.ContainerCreate(ctx, &container.Config{
		Image:      config.DockerImage,
		Cmd:        cmd,
		Tty:        false,
		WorkingDir: "/app",
	}, &container.HostConfig{
		Resources: container.Resources{
			Memory:    100 * 1024 * 1024, // 100MB limit
			CPUShares: 512,               // CPU limit
		},
		AutoRemove:  false,
		NetworkMode: "none", // Без сети для безопасности
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: codePath,
				Target: "/app",
			},
		},
	}, nil, nil, "")
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (s *DockerService) startContainer(ctx context.Context, containerID string) error {
	return s.client.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
}

func (s *DockerService) waitForCompletion(ctx context.Context, containerID string, config models.LanguageConfig) (*models.ExecutionResult, error) {
	statusCh, errCh := s.client.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		if err != nil {
			return nil, err
		}
	case <-statusCh:
	}

	// Получаем логи
	logs, err := s.getContainerLogs(ctx, containerID)
	if err != nil {
		return nil, err
	}

	// Проверяем статус выполнения
	inspect, err := s.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}

	result := &models.ExecutionResult{
		Output:  logs,
		Success: inspect.State.ExitCode == 0,
	}

	if inspect.State.ExitCode != 0 {
		result.Error = fmt.Sprintf("Exit code: %d", inspect.State.ExitCode)
	}

	return result, nil
}

func (s *DockerService) getContainerLogs(ctx context.Context, containerID string) (string, error) {
	reader, err := s.client.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     false,
	})
	if err != nil {
		return "", err
	}
	defer reader.Close()

	var output strings.Builder
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		// Docker добавляет префиксы к логам, пропускаем их
		line := scanner.Text()
		if len(line) > 8 {
			output.WriteString(line[8:])
			output.WriteString("\n")
		}
	}

	return strings.TrimSpace(output.String()), nil
}

func (s *DockerService) removeContainer(ctx context.Context, containerID string) {
	err := s.client.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		log.Printf("Warning: failed to remove container %s: %v", containerID, err)
	}
}
