package handlers

import (
	"backend/internal/services"
	"encoding/json"
	"fmt"
	"net/http"
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
