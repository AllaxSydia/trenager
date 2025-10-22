package handlers

import (
	"backend/internal/executor"
	"backend/internal/models"
	"backend/internal/services"
	"encoding/json"
	"log"
	"net/http"
)

// Короче, тут выбираем стратегию выполнения, либо Docker либо Локально
// потом преобразование результатов в единый формат ответа

var dockerService *services.DockerService // Изоляция
var localExecutor *executor.LocalExecutor // Быстро

func init() {
	var err error
	dockerService, err = services.NewDockerService()
	if err != nil {
		log.Printf("Warning: Docker service not available: %v", err)
		log.Println("Running in local execution mode")
	} else {
		log.Println("✅ Docker service initialized successfully")
	}

	// Инициализируем локальный исполнитель
	localExecutor = executor.NewLocalExecutor()
	// Если не получается создать Docker сервис для изолированного выполнения, то
	// Переходим в локальный режим
	// Локалка создаётся всегда
}

func ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"success": false, "message": "Only POST method allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	// Парсинг JSON
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
			response = executeCodeWithLocalExecutor(req.Code, req.Language)
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
		response = executeCodeWithLocalExecutor(req.Code, req.Language)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Новая функция использующая LocalExecutor для всех языков
func executeCodeWithLocalExecutor(code, language string) models.ExecutionResponse {
	log.Printf("🔧 Executing %s code with local executor", language)

	result, err := localExecutor.Execute(code, language)

	if err != nil {
		log.Printf("❌ Local execution error: %v", err)
		return models.ExecutionResponse{
			Success: false,
			Message: "Execution failed: " + err.Error(),
			Output:  "",
		}
	}

	// Преобразуем результат из LocalExecutor в models.ExecutionResponse
	exitCode := result["exitCode"].(int)
	output := result["output"].(string)
	errorMsg := result["error"].(string)

	success := exitCode == 0
	finalOutput := output
	if errorMsg != "" {
		finalOutput = errorMsg
		if output != "" {
			finalOutput = output + "\n" + errorMsg
		}
	}

	message := "Код выполнен успешно (локально)"
	if !success {
		message = "Ошибка выполнения кода"
	}

	log.Printf("✅ Local execution completed, success: %t, output length: %d", success, len(finalOutput))

	return models.ExecutionResponse{
		Success: success,
		Message: message,
		Output:  finalOutput,
	}
}
