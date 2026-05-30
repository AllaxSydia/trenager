package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	TaskID          uuid.UUID  `json:"task_id" db:"task_id"`
	UserID          uuid.UUID  `json:"user_id" db:"user_id"`
	Code            string     `json:"code" db:"code"`
	Language        string     `json:"language" db:"language"`
	Status          string     `json:"status" db:"status"`
	Score           int        `json:"score" db:"score"`
	MaxScore        int        `json:"max_score" db:"max_score"`
	Feedback        string     `json:"feedback" db:"feedback"`
	ExecutionTimeMs int        `json:"execution_time_ms" db:"execution_time_ms"`
	MemoryUsedBytes int64      `json:"memory_used_bytes" db:"memory_used_bytes"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	GradedAt        *time.Time `json:"graded_at" db:"graded_at"`
}

type TestResult struct {
	ID              uuid.UUID `json:"id" db:"id"`
	SubmissionID    uuid.UUID `json:"submission_id" db:"submission_id"`
	TestName        string    `json:"test_name" db:"test_name"`
	Passed          bool      `json:"passed" db:"passed"`
	Output          string    `json:"output" db:"output"`
	Expected        string    `json:"expected" db:"expected"`
	Score           int       `json:"score" db:"score"`
	MaxScore        int       `json:"max_score" db:"max_score"`
	Error           string    `json:"error" db:"error"`
	ExecutionTimeMs int       `json:"execution_time_ms" db:"execution_time_ms"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type TaskStatistics struct {
	TaskID                uuid.UUID `json:"task_id" db:"task_id"`
	TotalSubmissions      int       `json:"total_submissions" db:"total_submissions"`
	SuccessfulSubmissions int       `json:"successful_submissions" db:"successful_submissions"`
	AverageScore          float64   `json:"average_score" db:"average_score"`
	TotalExecutionTimeMs  int64     `json:"total_execution_time_ms" db:"total_execution_time_ms"`
	LastUpdated           time.Time `json:"last_updated" db:"last_updated"`
}

type GradeResult struct {
	SubmissionID    uuid.UUID
	Score           int
	MaxScore        int
	Percentage      float64
	Status          string
	Feedback        string
	TestResults     []TestResult
	ExecutionTimeMs int
	MemoryUsedBytes int64
}

// GradeRequest - запрос на проверку решения
type GradeRequest struct {
	TaskID   uuid.UUID
	UserID   uuid.UUID
	Code     string
	Language string
}

// TaskStats - статистика по задаче (для ответа)
type TaskStats struct {
	TotalSubmissions      int
	SuccessfulSubmissions int
	AverageScore          float64
	ScoreDistribution     map[string]int32
}
