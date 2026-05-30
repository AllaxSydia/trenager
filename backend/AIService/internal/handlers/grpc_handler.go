package handler

import (
	"context"
	"log"

	"ai-service/internal/ai"

	pb "github.com/AllaxSydia/trenager/proto/ai"
)

type AIHandler struct {
	pb.UnimplementedAIServiceServer
	aiService *ai.AIService
}

func NewAIHandler(aiService *ai.AIService) *AIHandler {
	return &AIHandler{
		aiService: aiService,
	}
}

func (h *AIHandler) GetHint(ctx context.Context, req *pb.GetHintRequest) (*pb.GetHintResponse, error) {
	log.Printf("📝 GetHint: task=%s, level=%d", req.TaskId, req.HintLevel)

	hint, exampleCode, err := h.aiService.GetHint(ctx, req.TaskId, req.UserCode, req.Language, int(req.HintLevel))
	if err != nil {
		return nil, err
	}

	return &pb.GetHintResponse{
		Hint:        hint,
		HintLevel:   req.HintLevel,
		ExampleCode: exampleCode,
	}, nil
}

func (h *AIHandler) ReviewCode(ctx context.Context, req *pb.ReviewCodeRequest) (*pb.ReviewCodeResponse, error) {
	log.Printf("🔍 ReviewCode: language=%s, code_length=%d", req.Language, len(req.Code))

	score, issues, suggestions, feedback, err := h.aiService.ReviewCode(ctx, req.Code, req.Language, req.TaskDescription)
	if err != nil {
		return nil, err
	}

	return &pb.ReviewCodeResponse{
		QualityScore:    int32(score),
		Issues:          issues,
		Suggestions:     suggestions,
		OverallFeedback: feedback,
	}, nil
}

func (h *AIHandler) GetRecommendations(ctx context.Context, req *pb.GetRecommendationsRequest) (*pb.GetRecommendationsResponse, error) {
	log.Printf("📚 GetRecommendations: user=%s", req.UserId)

	recs, err := h.aiService.GetRecommendations(ctx, req.UserId, req.CompletedTasks, req.WeakTopics)
	if err != nil {
		return nil, err
	}

	var pbRecs []*pb.TaskRecommendation
	for _, rec := range recs {
		pbRecs = append(pbRecs, &pb.TaskRecommendation{
			TaskId:          rec["task_id"].(string),
			Title:           rec["title"].(string),
			Reason:          rec["reason"].(string),
			DifficultyScore: float32(rec["difficulty_score"].(float64)),
		})
	}

	return &pb.GetRecommendationsResponse{
		Recommendations: pbRecs,
	}, nil
}

func (h *AIHandler) AskQuestion(ctx context.Context, req *pb.AskQuestionRequest) (*pb.AskQuestionResponse, error) {
	log.Printf("💬 AskQuestion: question=%s", req.Question[:min(50, len(req.Question))])

	answer, codeExamples, err := h.aiService.AskQuestion(ctx, req.Question, req.CodeContext, req.Language)
	if err != nil {
		return nil, err
	}

	return &pb.AskQuestionResponse{
		Answer:       answer,
		CodeExamples: codeExamples,
		HasCode:      len(codeExamples) > 0,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
