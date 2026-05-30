package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"TaskService/internal/models"
	"TaskService/pkg/database"

	"github.com/google/uuid"
)

type TaskRepository struct {
	db *database.Database
}

func NewTaskRepository(db *database.Database) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	query := `
        INSERT INTO tasks (id, title, description, language, difficulty, 
                        author_id, is_published, category, points, tags)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        RETURNING created_at, updated_at
    `

	tagsJSON, _ := json.Marshal(task.Tags)

	err := r.db.GetDB().QueryRowContext(ctx, query,
		task.ID, task.Title, task.Description, task.Language,
		task.Defficulty, task.AuthorID, task.IsPublished,
		task.Category, task.Points, tagsJSON,
	).Scan(&task.CreatedAt, &task.UpdateAt)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return r.saveTests(ctx, task.ID, task.Tests)
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	query := `
        SELECT id, title, description, language, difficulty, 
            author_id, created_at, updated_at, is_published, 
            category, points, tags
        FROM tasks WHERE id = $1
    `

	task := &models.Task{}
	var tagsJSON []byte

	err := r.db.GetDB().QueryRowContext(ctx, query, id).Scan(
		&task.ID, &task.Title, &task.Description, &task.Language,
		&task.Defficulty, &task.AuthorID, &task.CreatedAt, &task.UpdateAt,
		&task.IsPublished, &task.Category, &task.Points, &tagsJSON,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	json.Unmarshal(tagsJSON, &task.Tags)
	tests, _ := r.getTests(ctx, id)
	task.Tests = tests

	return task, nil
}

func (r *TaskRepository) List(ctx context.Context, filter *models.TaskFilter) ([]models.Task, int, error) {
	conditions := []string{"1=1"}
	args := []interface{}{}
	argIndex := 1

	if filter.Language != "" {
		conditions = append(conditions, fmt.Sprintf("language = $%d", argIndex))
		args = append(args, filter.Language)
		argIndex++
	}

	if filter.Difficulty != "" {
		conditions = append(conditions, fmt.Sprintf("difficulty = $%d", argIndex))
		args = append(args, filter.Difficulty)
		argIndex++
	}

	if filter.OnlyPublished {
		conditions = append(conditions, "is_published = true")
	}

	whereClause := strings.Join(conditions, " AND ")

	// Получаем общее количество
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tasks WHERE %s", whereClause)
	var total int
	r.db.GetDB().QueryRowContext(ctx, countQuery, args...).Scan(&total)

	// Пагинация
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}

	offset := (filter.Page - 1) * filter.PageSize

	query := fmt.Sprintf(`
        SELECT id, title, description, language, difficulty, 
            author_id, created_at, updated_at, is_published, 
            category, points, tags
        FROM tasks WHERE %s
        ORDER BY created_at DESC
        LIMIT $%d OFFSET $%d
    `, whereClause, argIndex, argIndex+1)

	args = append(args, filter.PageSize, offset)

	rows, err := r.db.GetDB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var tagsJSON []byte

		rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Language,
			&task.Defficulty, &task.AuthorID, &task.CreatedAt, &task.UpdateAt,
			&task.IsPublished, &task.Category, &task.Points, &tagsJSON,
		)

		json.Unmarshal(tagsJSON, &task.Tags)
		tasks = append(tasks, task)
	}

	return tasks, total, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	query := `
        UPDATE tasks 
        SET title = $1, description = $2, difficulty = $3, 
            category = $4, points = $5, tags = $6, updated_at = NOW()
        WHERE id = $7
    `

	tagsJSON, _ := json.Marshal(task.Tags)

	result, err := r.db.GetDB().ExecContext(ctx, query,
		task.Title, task.Description, task.Defficulty,
		task.Category, task.Points, tagsJSON, task.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := r.db.GetDB().ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *TaskRepository) saveTests(ctx context.Context, taskID uuid.UUID, tests []models.Test) error {
	query := `INSERT INTO tests (id, task_id, input, expected_output, description, is_hidden, timeout)
            VALUES ($1, $2, $3, $4, $5, $6, $7)`

	for _, test := range tests {
		testID := uuid.New()
		_, err := r.db.GetDB().ExecContext(ctx, query,
			testID, taskID, test.Input, test.ExpectedOutput,
			test.Description, test.IsHidden, test.Timeout,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TaskRepository) getTests(ctx context.Context, taskID uuid.UUID) ([]models.Test, error) {
	query := `SELECT input, expected_output, description, is_hidden, timeout
            FROM tests WHERE task_id = $1 ORDER BY created_at`

	rows, err := r.db.GetDB().QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []models.Test
	for rows.Next() {
		var test models.Test
		rows.Scan(&test.Input, &test.ExpectedOutput, &test.Description, &test.IsHidden, &test.Timeout)
		tests = append(tests, test)
	}
	return tests, nil
}
