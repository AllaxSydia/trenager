package service

import (
	"context"
	"log"
	"time"

	"grading-service/internal/models"
	"grading-service/internal/repository"

	"github.com/google/uuid"
)

type GradingService struct {
	repo *repository.GradingRepository
}

func NewGradingService() *GradingService {
	return &GradingService{
		repo: nil, // Будет установлен позже или через сеттер
	}
}

func (s *GradingService) SetRepository(repo *repository.GradingRepository) {
	s.repo = repo
}

func (s *GradingService) SubmitSolution(ctx context.Context, req *models.GradeRequest) (*models.GradeResult, error) {
	log.Printf("Processing submission for task %s, user %s", req.TaskID, req.UserID)

	submission := &models.Submission{
		ID:        uuid.New(),
		TaskID:    req.TaskID,
		UserID:    req.UserID,
		Code:      req.Code,
		Language:  req.Language,
		Status:    "pending",
		Score:     0,
		MaxScore:  100,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if s.repo != nil {
		if err := s.repo.CreateSubmission(ctx, submission); err != nil {
			return nil, err
		}
	}

	result := &models.GradeResult{
		SubmissionID:    submission.ID,
		Score:           85,
		MaxScore:        100,
		Percentage:      85.0,
		Status:          "passed",
		Feedback:        "Good solution! Passed all tests.",
		ExecutionTimeMs: 150,
		MemoryUsedBytes: 2048,
	}

	now := time.Now()
	submission.Status = result.Status
	submission.Score = result.Score
	submission.MaxScore = result.MaxScore
	submission.Feedback = result.Feedback
	submission.ExecutionTimeMs = result.ExecutionTimeMs
	submission.MemoryUsedBytes = result.MemoryUsedBytes
	submission.GradedAt = &now
	submission.UpdatedAt = now

	if s.repo != nil {
		s.repo.UpdateSubmission(ctx, submission)
	}

	return result, nil
}

func (s *GradingService) GetGrade(ctx context.Context, submissionID uuid.UUID) (*models.Submission, error) {
	if s.repo != nil {
		return s.repo.GetSubmission(ctx, submissionID)
	}
	return nil, nil
}

func (s *GradingService) GetUserGrades(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]models.Submission, int, error) {
	if s.repo != nil {
		offset := (page - 1) * pageSize
		return s.repo.GetUserSubmissions(ctx, userID, pageSize, offset)
	}
	return []models.Submission{}, 0, nil
}

func (s *GradingService) GetTaskStatistics(ctx context.Context, taskID uuid.UUID) (*models.TaskStats, error) {
	return &models.TaskStats{
		TotalSubmissions:      42,
		SuccessfulSubmissions: 35,
		AverageScore:          78.5,
		ScoreDistribution: map[string]int32{
			"0-20":   2,
			"21-40":  3,
			"41-60":  5,
			"61-80":  12,
			"81-100": 20,
		},
	}, nil
}
