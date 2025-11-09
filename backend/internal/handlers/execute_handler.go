package handlers

import (
	"backend/internal/executor"
	"backend/internal/models"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var codeExecutor executor.Executor

func init() {
	codeExecutor = executor.NewExecutor()
}

func ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"success": false, "message": "Only POST method allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Парсинг JSON
	var req struct {
		Code     string   `json:"code"`
		Language string   `json:"language"`
		Inputs   []string `json:"inputs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Failed to parse execute request: %v", err)
		http.Error(w, `{"success": false, "message": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Code == "" {
		http.Error(w, `{"success": false, "message": "Code is required"}`, http.StatusBadRequest)
		return
	}
	if req.Language == "" {
		http.Error(w, `{"success": false, "message": "Language is required"}`, http.StatusBadRequest)
		return
	}

	log.Printf("🔧 Executing code for language: %s, code length: %d, inputs: %v", req.Language, len(req.Code), req.Inputs)

	// Выполняем код через выбранный executor
	result, err := codeExecutor.Execute(req.Code, req.Language, req.Inputs)
	if err != nil {
		log.Printf("❌ Execution error: %v", err)
		http.Error(w, `{"success": false, "message": "Execution failed: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Форматируем ответ
	success := false
	if s, ok := result["success"].(bool); ok {
		success = s
	} else if ec, ok := result["exitCode"].(int); ok {
		success = ec == 0
	}

	output := ""
	if out, ok := result["output"].(string); ok {
		output = out
	}

	errorMsg := ""
	if err, ok := result["error"].(string); ok {
		errorMsg = err
	}

	// Комбинируем output и error если нужно
	finalOutput := output
	if errorMsg != "" {
		if output != "" {
			finalOutput = output + "\n" + errorMsg
		} else {
			finalOutput = errorMsg
		}
	}

	finalOutput = strings.TrimSpace(finalOutput)

	message := "✅ Код выполнен успешно"
	if !success {
		message = "❌ Ошибка выполнения кода"
	}

	response := models.ExecutionResponse{
		Success: success,
		Message: message,
		Output:  finalOutput,
	}

	log.Printf("📊 Execution completed - Success: %t, Output length: %d", success, len(finalOutput))

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("❌ Failed to encode response: %v", err)
		http.Error(w, `{"success": false, "message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
}
