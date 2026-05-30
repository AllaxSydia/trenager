package service

import (
	"TaskService/internal/models"
	"TaskService/internal/repository"
	"context"

	"github.com/google/uuid"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) error {
	return s.repo.Create(ctx, task)
}

func (s *TaskService) GetTask(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) ListTasks(ctx context.Context, filter *models.TaskFilter) ([]models.Task, int, error) {
	return s.repo.List(ctx, filter)
}

func (s *TaskService) UpdateTask(ctx context.Context, task *models.Task) error {
	return s.repo.Update(ctx, task)
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
