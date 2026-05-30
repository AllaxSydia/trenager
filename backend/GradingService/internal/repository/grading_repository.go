package repository

import (
	"context"
	"database/sql"
	"fmt"

	"grading-service/internal/models"
	"grading-service/pkg/database"

	"github.com/google/uuid"
)

type GradingRepository struct {
	db *database.Database
}

func NewGradingRepository(db *database.Database) *GradingRepository {
	return &GradingRepository{db: db}
}

func (r *GradingRepository) CreateSubmission(ctx context.Context, sub *models.Submission) error {
	query := `
        INSERT INTO submissions (id, task_id, user_id, code, language, status, 
                                score, max_score, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	_, err := r.db.GetDB().ExecContext(ctx, query,
		sub.ID, sub.TaskID, sub.UserID, sub.Code, sub.Language,
		sub.Status, sub.Score, sub.MaxScore, sub.CreatedAt, sub.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create submission: %w", err)
	}

	return nil
}

func (r *GradingRepository) UpdateSubmission(ctx context.Context, sub *models.Submission) error {
	query := `
        UPDATE submissions 
        SET status = $1, score = $2, max_score = $3, feedback = $4,
            execution_time_ms = $5, memory_used_bytes = $6, updated_at = $7, graded_at = $8
        WHERE id = $9
    `

	_, err := r.db.GetDB().ExecContext(ctx, query,
		sub.Status, sub.Score, sub.MaxScore, sub.Feedback,
		sub.ExecutionTimeMs, sub.MemoryUsedBytes, sub.UpdatedAt, sub.GradedAt, sub.ID)

	if err != nil {
		return fmt.Errorf("failed to update submission: %w", err)
	}

	return nil
}

func (r *GradingRepository) GetSubmission(ctx context.Context, id uuid.UUID) (*models.Submission, error) {
	query := `
        SELECT id, task_id, user_id, code, language, status, score, max_score,
            feedback, execution_time_ms, memory_used_bytes, created_at, updated_at, graded_at
        FROM submissions WHERE id = $1
    `

	sub := &models.Submission{}
	err := r.db.GetDB().QueryRowContext(ctx, query, id).Scan(
		&sub.ID, &sub.TaskID, &sub.UserID, &sub.Code, &sub.Language,
		&sub.Status, &sub.Score, &sub.MaxScore, &sub.Feedback,
		&sub.ExecutionTimeMs, &sub.MemoryUsedBytes, &sub.CreatedAt, &sub.UpdatedAt, &sub.GradedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get submission: %w", err)
	}

	return sub, nil
}

func (r *GradingRepository) SaveTestResult(ctx context.Context, result *models.TestResult) error {
	query := `
        INSERT INTO test_results (id, submission_id, test_name, passed, output, 
                                expected, score, max_score, error, execution_time_ms)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    `

	_, err := r.db.GetDB().ExecContext(ctx, query,
		result.ID, result.SubmissionID, result.TestName, result.Passed,
		result.Output, result.Expected, result.Score, result.MaxScore,
		result.Error, result.ExecutionTimeMs)

	return err
}

func (r *GradingRepository) GetUserSubmissions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.Submission, int, error) {
	countQuery := `SELECT COUNT(*) FROM submissions WHERE user_id = $1`
	var total int
	r.db.GetDB().QueryRowContext(ctx, countQuery, userID).Scan(&total)

	query := `
        SELECT id, task_id, user_id, code, language, status, score, max_score,
            feedback, execution_time_ms, memory_used_bytes, created_at, updated_at, graded_at
        FROM submissions
        WHERE user_id = $1
        ORDER BY created_at DESC
        LIMIT $2 OFFSET $3
    `

	rows, err := r.db.GetDB().QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var submissions []models.Submission
	for rows.Next() {
		var sub models.Submission
		err := rows.Scan(&sub.ID, &sub.TaskID, &sub.UserID, &sub.Code, &sub.Language,
			&sub.Status, &sub.Score, &sub.MaxScore, &sub.Feedback,
			&sub.ExecutionTimeMs, &sub.MemoryUsedBytes, &sub.CreatedAt, &sub.UpdatedAt, &sub.GradedAt)
		if err != nil {
			return nil, 0, err
		}
		submissions = append(submissions, sub)
	}

	return submissions, total, nil
}

func (r *GradingRepository) UpdateTaskStatistics(ctx context.Context, taskID uuid.UUID, score int, maxScore int, executionTimeMs int) error {
	query := `
        INSERT INTO task_statistics (task_id, total_submissions, successful_submissions, 
                                    average_score, total_execution_time_ms, last_updated)
        VALUES ($1, 1, $2, $3, $4, NOW())
        ON CONFLICT (task_id) DO UPDATE SET
            total_submissions = task_statistics.total_submissions + 1,
            successful_submissions = task_statistics.successful_submissions + $2,
            average_score = ((task_statistics.average_score * task_statistics.total_submissions) + $3) / (task_statistics.total_submissions + 1),
            total_execution_time_ms = task_statistics.total_execution_time_ms + $4,
            last_updated = NOW()
    `

	success := 0
	if score >= maxScore {
		success = 1
	}
	averageScore := float64(score) / float64(maxScore) * 100

	_, err := r.db.GetDB().ExecContext(ctx, query, taskID, success, averageScore, executionTimeMs)
	return err
}
