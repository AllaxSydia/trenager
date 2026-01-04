package handlers

import (
	"backend/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type AIReviewRequest struct {
	Code        string `json:"code"`
	Language    string `json:"language"`
	TaskContext string `json:"task_context"`
}

func AIReviewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Добавим логирование
	fmt.Println("🎯 AI Review Handler: Received request")

	var req AIReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Printf("❌ AI Review Handler: JSON decode error: %v\n", err)
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	fmt.Printf("🔍 AI Review Handler: Code: %s, Language: %s\n", req.Language, req.TaskContext)

	if req.Code == "" || req.Language == "" {
		fmt.Println("❌ AI Review Handler: Missing code or language")
		http.Error(w, `{"error": "Code and language are required"}`, http.StatusBadRequest)
		return
	}

	reviewer := services.NewAIReviewer()
	reviewRequest := services.CodeReviewRequest{
		Code:        req.Code,
		Language:    req.Language,
		TaskContext: req.TaskContext,
		UserId:      1,
	}

	fmt.Println("🔧 AI Review Handler: Calling ReviewCode...")
	response, err := reviewer.ReviewCode(reviewRequest)
	if err != nil {
		fmt.Printf("❌ AI Review Handler: Review failed: %v\n", err)
		http.Error(w, `{"error": "AI review failed"}`, http.StatusInternalServerError)
		return
	}

	fmt.Printf("✅ AI Review Handler: Success, score: %d\n", response.Score)
	json.NewEncoder(w).Encode(response)
}

func AIHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":    "ok",
		"service":   "ai_code_review",
		"timestamp": time.Now().UTC(),
		"endpoints": map[string]string{
			"/api/ai/review": "POST - Analyze code with AI",
		},
	}

	// Проверяем, есть ли API ключ
	openrouterKey := os.Getenv("OPENROUTER_API_KEY")
	openaiKey := os.Getenv("OPENAI_API_KEY")

	if openrouterKey != "" {
		response["ai_provider"] = "OpenRouter"
		response["ai_status"] = "configured"
		response["has_key"] = true
	} else if openaiKey != "" {
		response["ai_provider"] = "OpenAI"
		response["ai_status"] = "configured"
		response["has_key"] = true
	} else {
		response["ai_provider"] = "mock"
		response["ai_status"] = "mock_mode"
		response["has_key"] = false
		response["message"] = "Add OPENROUTER_API_KEY or OPENAI_API_KEY to .env for real AI"
	}

	json.NewEncoder(w).Encode(response)
}
