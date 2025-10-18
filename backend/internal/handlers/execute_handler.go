package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var dockerService *services.DockerService

func init() {
	var err error
	dockerService, err = services.NewDockerService()
	if err != nil {
		log.Printf("Warning: Docker service not available: %v", err)
		log.Println("Running in local execution mode")
	} else {
		log.Println("✅ Docker service initialized successfully")
	}
}

func ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"success": false, "message": "Only POST method allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req models.ExecutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success": false, "message": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	log.Printf("🔧 Executing code for language: %s", req.Language)

	var response models.ExecutionResponse

	// Пробуем Docker сначала
	if dockerService != nil {
		log.Println("🐳 Attempting Docker execution...")
		result, err := dockerService.ExecuteCode(req.Code, req.Language)
		if err != nil {
			log.Printf("❌ Docker execution failed: %v", err)
			log.Println("🔄 Falling back to local execution...")
			response = executeCodeLocally(req.Code, req.Language)
		} else {
			log.Printf("✅ Docker execution successful, output: %s", result.Output)
			response = models.ExecutionResponse{
				Success: result.Success,
				Message: "Code executed successfully via Docker",
				Output:  result.Output,
			}
			if !result.Success {
				response.Message = "Code execution failed in Docker"
			}
		}
	} else {
		log.Println("🔄 Docker not available, using local execution...")
		response = executeCodeLocally(req.Code, req.Language)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func executeCodeLocally(code, language string) models.ExecutionResponse {
	// Поддерживаем только Go для локального выполнения
	if language != "go" {
		return models.ExecutionResponse{
			Success: false,
			Message: "Локальное выполнение поддерживает только Go",
			Output:  "Для выполнения кода на других языках требуется Docker",
		}
	}

	log.Printf("🔧 Executing Go code locally: %s", code)

	// Создаем временную директорию
	tmpDir, err := os.MkdirTemp("", "go_exec_*")
	if err != nil {
		return models.ExecutionResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create temp dir: %v", err),
			Output:  "",
		}
	}
	defer os.RemoveAll(tmpDir)

	// Создаем файл с кодом
	filePath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(filePath, []byte(code), 0644); err != nil {
		return models.ExecutionResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to write code file: %v", err),
			Output:  "",
		}
	}

	log.Printf("📁 Code written to: %s", filePath)

	// Выполняем с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", filePath)
	output, err := cmd.CombinedOutput()

	// Обрабатываем результат
	if ctx.Err() == context.DeadlineExceeded {
		return models.ExecutionResponse{
			Success: false,
			Message: "Execution timeout (30 seconds exceeded)",
			Output:  "Код выполнялся слишком долго",
		}
	}

	if err != nil {
		return models.ExecutionResponse{
			Success: false,
			Message: fmt.Sprintf("Execution error: %v", err),
			Output:  string(output),
		}
	}

	log.Printf("✅ Local execution successful, output: %s", string(output))

	return models.ExecutionResponse{
		Success: true,
		Message: "Код выполнен успешно (локально)",
		Output:  string(output),
	}
}
