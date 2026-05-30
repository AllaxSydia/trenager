package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Language    string    `json:"language" db:"language"`
	Defficulty  string    `json:"difficulty" db:"difficulty"`
	Template    string    `json:"template" db:"template"`
	StarterCode string    `json:"starter_code" db:"starter_code"`
	Tests       []Test    `json:"tests" db:"-"`
	AuthorID    uuid.UUID `json:"author_id" db:"author_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdateAt    time.Time `json:"updated_at" db:"updated_at"`
	IsPublished bool      `json:"is_published" db:"is_published"`
	Category    string    `json:"category" db:"category"`
	Points      int       `json:"points" db:"points"`
	Tags        []string  `json:"tags" db:"-"`
}

type Test struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	Description    string `json:"description,omitempty"`
	IsHidden       bool   `json:"is_hidden,omitempty"`
	Timeout        int    `json:"timeout,omitempty"`
}

type TaskFilter struct {
	Language      string   `json:"language,omitempty"`
	Difficulty    string   `json:"difficulty,omitempty"`
	Category      string   `json:"category,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Search        string   `json:"search,omitempty"`
	Page          int      `json:"page,omitempty"`
	PageSize      int      `json:"page_size,omitempty"`
	SortBy        string   `json:"sort_by,omitempty"`
	SortOrder     string   `json:"sort_order,omitempty"`
	OnlyPublished bool     `json:"only_published,omitempty"`
	AuthorID      string   `json:"author_id,omitempty"`
}

type TaskStats struct {
	TaskID          string     `json:"task_id"`
	TotalAttempts   int        `json:"total_attempts"`
	SuccessAttempts int        `json:"success_attempts"`
	SuccessRate     float64    `json:"success_rate"`
	AvgTime         float64    `json:"avg_time"`
	FirstSolved     *time.Time `json:"first_solved,omitempty"`
	LastAttempt     *time.Time `json:"last_attempt,omitempty"`
}
