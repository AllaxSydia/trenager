package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var dockerService *services.DockerService

func init() {
	var err error
	dockerService, err = services.NewDockerService()
	if err != nil {
		log.Printf("Warning: Docker service not available: %v", err)
		log.Println("Running in mock mode - code execution will be simulated")
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

	var response models.ExecutionResponse

	// Если Docker доступен - выполняем код по-настоящему
	if dockerService != nil {
		result, err := dockerService.ExecuteCode(req.Code, req.Language)
		if err != nil {
			log.Printf("Docker execution failed: %v", err)
			response = models.ExecutionResponse{
				Success: false,
				Message: fmt.Sprintf("Execution error: %v", err),
				Output:  "",
			}
		} else {
			response = models.ExecutionResponse{
				Success: result.Success,
				Message: "Code executed successfully",
				Output:  result.Output,
			}
			if !result.Success {
				response.Message = "Code execution failed"
			}
		}
	} else {
		// Режим эмуляции когда Docker недоступен
		response = simulateExecution(req.Code, req.Language)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func simulateExecution(code, language string) models.ExecutionResponse {
	// Простая эмуляция выполнения кода
	return models.ExecutionResponse{
		Success: true,
		Message: "Код выполнен успешно! (Режим эмуляции - Docker недоступен)",
		Output:  fmt.Sprintf("Эмуляция выполнения %s кода:\n\n%s\n\n--- Результат ---\nHello, World! (эмуляция)", language, code),
	}
}
