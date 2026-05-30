package handler

import (
	"context"
	"log"

	"grading-service/internal/models"
	"grading-service/internal/service"

	pb "github.com/AllaxSydia/trenager/proto/grading"
	"github.com/google/uuid"
)

type GradingHandler struct {
	pb.UnimplementedGradingServiceServer
	gradingService *service.GradingService
}

func NewGradingHandler(gradingService *service.GradingService) *GradingHandler {
	return &GradingHandler{
		gradingService: gradingService,
	}
}

func (h *GradingHandler) SubmitSolution(ctx context.Context, req *pb.SubmitSolutionRequest) (*pb.SubmitSolutionResponse, error) {
	log.Printf("📝 SubmitSolution: task_id=%s, user_id=%s, language=%s",
		req.TaskId, req.UserId, req.Language)

	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		return &pb.SubmitSolutionResponse{
			Success: false,
			Message: "Invalid task ID",
		}, nil
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &pb.SubmitSolutionResponse{
			Success: false,
			Message: "Invalid user ID",
		}, nil
	}

	gradeReq := &models.GradeRequest{
		TaskID:   taskID,
		UserID:   userID,
		Code:     req.Code,
		Language: req.Language,
	}

	result, err := h.gradingService.SubmitSolution(ctx, gradeReq)
	if err != nil {
		return &pb.SubmitSolutionResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SubmitSolutionResponse{
		Success:      true,
		SubmissionId: result.SubmissionID.String(),
		Message:      "Solution submitted successfully",
		Grade: &pb.Grade{
			SubmissionId: result.SubmissionID.String(),
			TaskId:       req.TaskId,
			UserId:       req.UserId,
			Score:        int32(result.Score),
			MaxScore:     int32(result.MaxScore),
			Status:       result.Status,
			Feedback:     result.Feedback,
		},
	}, nil
}

func (h *GradingHandler) GetGrade(ctx context.Context, req *pb.GetGradeRequest) (*pb.Grade, error) {
	log.Printf("📊 GetGrade: submission_id=%s", req.SubmissionId)

	submissionID, err := uuid.Parse(req.SubmissionId)
	if err != nil {
		return nil, err
	}

	grade, err := h.gradingService.GetGrade(ctx, submissionID)
	if err != nil {
		return nil, err
	}

	return &pb.Grade{
		SubmissionId: grade.ID.String(),
		TaskId:       grade.TaskID.String(),
		UserId:       grade.UserID.String(),
		Score:        int32(grade.Score),
		MaxScore:     int32(grade.MaxScore),
		Status:       grade.Status,
		Feedback:     grade.Feedback,
	}, nil
}

func (h *GradingHandler) GetUserGrades(ctx context.Context, req *pb.GetUserGradesRequest) (*pb.GetUserGradesResponse, error) {
	log.Printf("📋 GetUserGrades: user_id=%s, page=%d", req.UserId, req.Page)

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	grades, total, err := h.gradingService.GetUserGrades(ctx, userID, int(req.Page), int(req.PageSize))
	if err != nil {
		return nil, err
	}

	var pbGrades []*pb.Grade
	for _, grade := range grades {
		pbGrades = append(pbGrades, &pb.Grade{
			SubmissionId: grade.ID.String(),
			TaskId:       grade.TaskID.String(),
			UserId:       grade.UserID.String(),
			Score:        int32(grade.Score),
			MaxScore:     int32(grade.MaxScore),
			Status:       grade.Status,
			Feedback:     grade.Feedback,
		})
	}

	return &pb.GetUserGradesResponse{
		Grades: pbGrades,
		Total:  int32(total),
	}, nil
}

func (h *GradingHandler) GetTaskStatistics(ctx context.Context, req *pb.GetTaskStatisticsRequest) (*pb.TaskStatistics, error) {
	log.Printf("📈 GetTaskStatistics: task_id=%s", req.TaskId)

	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		return nil, err
	}

	stats, err := h.gradingService.GetTaskStatistics(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return &pb.TaskStatistics{
		TaskId:                req.TaskId,
		TotalSubmissions:      int32(stats.TotalSubmissions),
		SuccessfulSubmissions: int32(stats.SuccessfulSubmissions),
		AverageScore:          float32(stats.AverageScore),
		ScoreDistribution:     stats.ScoreDistribution,
	}, nil
}
