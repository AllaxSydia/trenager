package executor

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type DockerExecutorImpl struct {
	client *client.Client
}

func NewDockerExecutor() (*DockerExecutorImpl, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	return &DockerExecutorImpl{client: cli}, nil
}

// Execute реализует интерфейс Executor - тот же что и у LocalExecutor
func (d *DockerExecutorImpl) Execute(code, language string, inputs []string) (map[string]interface{}, error) {
	log.Printf("🐳 DockerExecutor executing %s code, length: %d chars, inputs: %v", language, len(code), inputs)

	ctx := context.Background()

	// Создаем временный файл с кодом
	tmpDir, err := os.MkdirTemp("", "code_executor_docker")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Printf("⚠️ Failed to remove temp directory %s: %v", tmpDir, err)
		}
	}()

	// Определяем настройки для языка
	config, err := d.getLanguageConfig(language)
	if err != nil {
		return nil, err
	}

	// Записываем код в файл
	codeFile := filepath.Join(tmpDir, config.FileName)
	if err := os.WriteFile(codeFile, []byte(code), 0644); err != nil {
		return nil, fmt.Errorf("failed to write code file: %w", err)
	}

	// Подготавливаем входные данные
	inputData := strings.Join(inputs, "\n")

	// Создаем контейнер
	resp, err := d.client.ContainerCreate(ctx, &container.Config{
		Image:        config.Image,
		Cmd:          config.Command,
		WorkingDir:   "/code",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		StdinOnce:    true,
		Tty:          false,
	}, &container.HostConfig{
		Binds:       []string{tmpDir + ":/code:ro"},
		NetworkMode: "none",
		Resources: container.Resources{
			Memory:   100 * 1024 * 1024, // 100MB
			NanoCPUs: 500000000,         // 0.5 CPU
		},
		AutoRemove: false, // Не удалять автоматически, чтобы можно было получить статус
	}, nil, nil, "")

	if err != nil {
		return nil, fmt.Errorf("failed to create container: %w", err)
	}
	defer func() {
		// Удаляем контейнер после использования
		if err := d.client.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
			log.Printf("⚠️ Failed to remove container %s: %v", resp.ID, err)
		}
	}()

	// Запускаем контейнер
	if err := d.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	// Подключаемся к контейнеру
	attach, err := d.client.ContainerAttach(ctx, resp.ID, types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to attach to container: %w", err)
	}
	defer attach.Close()

	// Отправляем входные данные
	if inputData != "" {
		_, err = attach.Conn.Write([]byte(inputData + "\n"))
		if err != nil {
			log.Printf("⚠️ Failed to write input data: %v", err)
		}
	}
	attach.Conn.Close()

	// Ждем завершения с таймаутом
	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	statusCh, errCh := d.client.ContainerWait(timeoutCtx, resp.ID, container.WaitConditionNotRunning)

	var exitCode int64 = 1
	select {
	case err := <-errCh:
		if err != nil {
			if strings.Contains(err.Error(), "context deadline exceeded") {
				log.Printf("⏰ Docker execution timeout (30 seconds)")
				// Останавливаем контейнер при таймауте
				d.client.ContainerStop(ctx, resp.ID, container.StopOptions{})
				return map[string]interface{}{
					"output":   "",
					"error":    "Execution timeout (30 seconds)",
					"exitCode": 1,
				}, nil
			}
			return nil, fmt.Errorf("container wait error: %w", err)
		}
	case status := <-statusCh:
		exitCode = status.StatusCode
	}

	// Получаем логи
	out, err := d.client.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get container logs: %w", err)
	}
	defer out.Close()

	// Читаем вывод
	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, out)
	if err != nil {
		return nil, fmt.Errorf("failed to read container output: %w", err)
	}

	success := exitCode == 0
	if success {
		log.Printf("✅ Docker execution completed successfully")
	} else {
		log.Printf("⚠️ Docker execution completed with exit code %d", exitCode)
	}

	result := map[string]interface{}{
		"success":  success,
		"output":   stdout.String(),
		"error":    stderr.String(),
		"exitCode": int(exitCode),
	}

	log.Printf("📊 Docker execution result - Output: %d chars, Error: %d chars",
		len(stdout.String()), len(stderr.String()))

	return result, nil
}

type DockerLanguageConfig struct {
	Image    string
	Command  []string
	FileName string
}

func (d *DockerExecutorImpl) getLanguageConfig(language string) (*DockerLanguageConfig, error) {
	switch strings.ToLower(language) {
	case "python", "python3":
		return &DockerLanguageConfig{
			Image:    "python:3.11-alpine",
			Command:  []string{"python", "script.py"},
			FileName: "script.py",
		}, nil
	case "javascript", "node":
		return &DockerLanguageConfig{
			Image:    "node:18-alpine",
			Command:  []string{"node", "script.js"},
			FileName: "script.js",
		}, nil
	case "java":
		return &DockerLanguageConfig{
			Image:    "openjdk:11-jre-alpine",
			Command:  []string{"java", "Main"},
			FileName: "Main.java",
		}, nil
	case "cpp", "c++":
		return &DockerLanguageConfig{
			Image:    "gcc:latest",
			Command:  []string{"sh", "-c", "g++ -o program main.cpp && ./program"},
			FileName: "main.cpp",
		}, nil
	case "go":
		return &DockerLanguageConfig{
			Image:    "golang:1.21-alpine",
			Command:  []string{"sh", "-c", "go run main.go"},
			FileName: "main.go",
		}, nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}

// Cleanup для совместимости с интерфейсом
func (d *DockerExecutorImpl) Cleanup() {
	// Docker сам управляет контейнерами, но можно добавить cleanup старых контейнеров если нужно
	log.Printf("🧹 DockerExecutor cleanup completed")
}
