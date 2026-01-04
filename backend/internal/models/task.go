package models

import (
	"time"
)

// Task - основная структура задачи для БД и API
type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Language    string `json:"language,omitempty"`
	Difficulty  string `json:"difficulty,omitempty"`
	Template    string `json:"template"`
	Tests       []Test `json:"tests"`

	// Новые поля для расширенной функциональности
	StarterCode string    `json:"starter_code,omitempty"` // Начальный код для студентов
	CreatedBy   int       `json:"created_by,omitempty"`   // ID пользователя-создателя
	CreatedAt   time.Time `json:"created_at,omitempty"`   // Дата создания
	UpdatedAt   time.Time `json:"updated_at,omitempty"`   // Дата обновления
	IsPublished bool      `json:"is_published,omitempty"` // Опубликована ли задача
	Category    string    `json:"category,omitempty"`     // Категория задачи
	Points      int       `json:"points,omitempty"`       // Очки за решение
	Tags        []string  `json:"tags,omitempty"`         // Теги для поиска
}

// Test - тест для задачи
type Test struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	Description    string `json:"description,omitempty"` // Описание теста (для учителей)
	IsHidden       bool   `json:"is_hidden,omitempty"`   // Скрытый тест (только для проверки)
	Timeout        int    `json:"timeout,omitempty"`     // Таймаут в мс
}

// ExecutionRequest - запрос на выполнение кода
type ExecutionRequest struct {
	TaskID   string   `json:"task_id"`
	Code     string   `json:"code"`
	Language string   `json:"language"`
	Inputs   []string `json:"inputs,omitempty"`
}

// ExecutionResponse - ответ от выполнения кода
type ExecutionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Output  string `json:"output"`
}

// CheckRequest - запрос на проверку решения
type CheckRequest struct {
	TaskID   interface{} `json:"task_id"` // принимает и строки и числа
	Code     string      `json:"code"`
	Language string      `json:"language"`
	Tests    []Test      `json:"tests,omitempty"`
}

// CheckResponse - ответ проверки решения
type CheckResponse struct {
	Success     bool         `json:"success"`
	Message     string       `json:"message"`
	TestResults []TestResult `json:"test_results"`
	TotalTests  int          `json:"total_tests"`
	PassedTests int          `json:"passed_tests"`
	Score       int          `json:"score,omitempty"`        // Оценка за решение
	TimeElapsed int64        `json:"time_elapsed,omitempty"` // Время выполнения в мс
}

// TestResult - результат выполнения одного теста
type TestResult struct {
	TestNumber int    `json:"test_number"`
	Passed     bool   `json:"passed"`
	Output     string `json:"output"`
	Error      string `json:"error,omitempty"`
	Expected   string `json:"expected"`
	Actual     string `json:"actual"`
	Input      string `json:"input,omitempty"`     // Входные данные теста
	IsHidden   bool   `json:"is_hidden,omitempty"` // Был ли тест скрытым
	Timeout    bool   `json:"timeout,omitempty"`   // Был ли превышен таймаут
}

// CheckResult - устаревшая структура (оставлена для обратной совместимости)
type CheckResult struct {
	Passed   bool   `json:"passed"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Message  string `json:"message"`
}

// ============ НОВЫЕ СТРУКТУРЫ ДЛЯ УПРАВЛЕНИЯ ЗАДАЧАМИ ============

// TaskRequest - запрос на создание/обновление задачи (от учителя)
type TaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Language    string `json:"language" binding:"required"`
	Difficulty  string `json:"difficulty" binding:"required"`
	Template    string `json:"template"`
	StarterCode string `json:"starter_code"`
	Tests       []Test `json:"tests" binding:"required,min=1"`
	Category    string `json:"category"`
	Points      int    `json:"points"`
	Tags        string `json:"tags"` // Теги через запятую
	IsPublished bool   `json:"is_published"`
}

// TaskResponse - ответ с задачей
type TaskResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Language    string    `json:"language"`
	Difficulty  string    `json:"difficulty"`
	Template    string    `json:"template"`
	StarterCode string    `json:"starter_code"`
	Tests       []Test    `json:"tests"`
	CreatedBy   int       `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsPublished bool      `json:"is_published"`
	Category    string    `json:"category"`
	Points      int       `json:"points"`
	Tags        []string  `json:"tags"`
	AuthorName  string    `json:"author_name,omitempty"`  // Имя создателя
	SolvedCount int       `json:"solved_count,omitempty"` // Сколько раз решили
}

// TaskListResponse - список задач
type TaskListResponse struct {
	Tasks      []TaskResponse `json:"tasks"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}

// TaskFilter - фильтры для поиска задач
type TaskFilter struct {
	Language      string `json:"language,omitempty"`
	Difficulty    string `json:"difficulty,omitempty"`
	Category      string `json:"category,omitempty"`
	Tags          string `json:"tags,omitempty"`
	Search        string `json:"search,omitempty"`
	Page          int    `json:"page,omitempty"`
	PageSize      int    `json:"page_size,omitempty"`
	SortBy        string `json:"sort_by,omitempty"`    // created_at, difficulty, title
	SortOrder     string `json:"sort_order,omitempty"` // asc, desc
	OnlyPublished bool   `json:"only_published,omitempty"`
}

// TaskStats - статистика по задаче
type TaskStats struct {
	TaskID          string     `json:"task_id"`
	TotalAttempts   int        `json:"total_attempts"`
	SuccessAttempts int        `json:"success_attempts"`
	SuccessRate     float64    `json:"success_rate"`
	AvgTime         float64    `json:"avg_time"` // Среднее время решения в секундах
	FirstSolved     *time.Time `json:"first_solved,omitempty"`
	LastAttempt     *time.Time `json:"last_attempt,omitempty"`
}

// UserTaskProgress - прогресс пользователя по задаче
type UserTaskProgress struct {
	TaskID      string     `json:"task_id"`
	TaskTitle   string     `json:"task_title"`
	Language    string     `json:"language"`
	Attempts    int        `json:"attempts"`
	IsSolved    bool       `json:"is_solved"`
	BestScore   int        `json:"best_score"`
	LastAttempt *time.Time `json:"last_attempt,omitempty"`
	FirstSolved *time.Time `json:"first_solved,omitempty"`
}

// BulkTestRequest - запрос на массовое тестирование
type BulkTestRequest struct {
	TaskID   string `json:"task_id"`
	Code     string `json:"code"`
	Language string `json:"language"`
}

// SolutionRequest - запрос на сохранение решения
type SolutionRequest struct {
	TaskID    string `json:"task_id" binding:"required"`
	Code      string `json:"code" binding:"required"`
	Language  string `json:"language" binding:"required"`
	IsCorrect bool   `json:"is_correct"`
	Score     int    `json:"score"`
	TimeSpent int64  `json:"time_spent"` // Время в мс
}

// SolutionResponse - ответ с решением
type SolutionResponse struct {
	ID          string       `json:"id"`
	TaskID      string       `json:"task_id"`
	UserID      int          `json:"user_id"`
	Code        string       `json:"code"`
	Language    string       `json:"language"`
	IsCorrect   bool         `json:"is_correct"`
	Score       int          `json:"score"`
	TimeSpent   int64        `json:"time_spent"`
	CreatedAt   time.Time    `json:"created_at"`
	TestResults []TestResult `json:"test_results,omitempty"`
}

// AIAnalysisRequest - запрос на AI анализ кода
type AIAnalysisRequest struct {
	Code      string `json:"code" binding:"required"`
	Language  string `json:"language" binding:"required"`
	TaskID    string `json:"task_id,omitempty"`
	TaskTitle string `json:"task_title,omitempty"`
}

// AIAnalysisResponse - ответ AI анализа
type AIAnalysisResponse struct {
	Score                int      `json:"score"`
	Comments             []string `json:"comments"`
	Suggestions          []string `json:"suggestions"`
	BestPractices        []string `json:"best_practices"`
	Complexity           string   `json:"complexity"`
	AlternativeSolutions []string `json:"alternative_solutions"`
	CodeQuality          string   `json:"code_quality,omitempty"` // excellent, good, average, poor
	SecurityIssues       []string `json:"security_issues,omitempty"`
	PerformanceTips      []string `json:"performance_tips,omitempty"`
}
