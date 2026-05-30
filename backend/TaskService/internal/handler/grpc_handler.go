package handler

import (
	"TaskService/internal/models"
	"TaskService/internal/service"
	"context"

	pb "github.com/AllaxSydia/trenager/proto/task" // путь к сгенерированному proto
	"github.com/google/uuid"
)

type TaskGrpcHandler struct {
	pb.UnimplementedTaskServiceServer
	taskService *service.TaskService
}

func NewTaskGrpcHandler(taskService *service.TaskService) *TaskGrpcHandler {
	return &TaskGrpcHandler{
		taskService: taskService,
	}
}

// CreateTask - создание задачи
func (h *TaskGrpcHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
	task := &models.Task{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Language:    req.Language,
		Defficulty:  req.Difficulty,
		Points:      calculatePoints(req.Difficulty),
	}

	// Конвертируем тесты
	for _, test := range req.TestCases {
		task.Tests = append(task.Tests, models.Test{
			Input:          test.Input,
			ExpectedOutput: test.ExpectedOutput,
			IsHidden:       test.IsHidden,
			Timeout:        int(req.TimeLimitMs),
		})
	}

	err := h.taskService.CreateTask(ctx, task)
	if err != nil {
		return &pb.TaskResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &pb.TaskResponse{
		Success: true,
		TaskId:  task.ID.String(),
		Message: "Task created successfully",
	}, nil
}

// GetTask - получение задачи по ID
func (h *TaskGrpcHandler) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.Task, error) {
	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		return nil, err
	}

	task, err := h.taskService.GetTask(ctx, taskID)
	if err != nil {
		return nil, err
	}

	return convertToProtoTask(task), nil
}

// ListTasks - список задач с фильтрацией
func (h *TaskGrpcHandler) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	filter := &models.TaskFilter{
		Language:   req.Language,
		Difficulty: req.Difficulty,
		Page:       int(req.Page),
		PageSize:   int(req.PageSize),
	}

	tasks, total, err := h.taskService.ListTasks(ctx, filter)
	if err != nil {
		return nil, err
	}

	var pbTasks []*pb.Task
	for _, task := range tasks {
		pbTasks = append(pbTasks, convertToProtoTask(&task))
	}

	return &pb.ListTasksResponse{
		Tasks:    pbTasks,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// UpdateTask - обновление задачи
func (h *TaskGrpcHandler) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.TaskResponse, error) {
	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		return &pb.TaskResponse{Success: false, Error: "invalid task id"}, nil
	}

	task := &models.Task{
		ID:          taskID,
		Title:       req.Title,
		Description: req.Description,
		Defficulty:  req.Difficulty,
	}

	err = h.taskService.UpdateTask(ctx, task)
	if err != nil {
		return &pb.TaskResponse{Success: false, Error: err.Error()}, nil
	}

	return &pb.TaskResponse{
		Success: true,
		Message: "Task updated successfully",
	}, nil
}

// DeleteTask - удаление задачи
func (h *TaskGrpcHandler) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		return &pb.DeleteTaskResponse{Success: false, Message: "invalid task id"}, nil
	}

	err = h.taskService.DeleteTask(ctx, taskID)
	if err != nil {
		return &pb.DeleteTaskResponse{Success: false, Message: err.Error()}, nil
	}

	return &pb.DeleteTaskResponse{
		Success: true,
		Message: "Task deleted successfully",
	}, nil
}

// Вспомогательные функции
func calculatePoints(difficulty string) int {
	switch difficulty {
	case "easy":
		return 100
	case "medium":
		return 200
	case "hard":
		return 300
	default:
		return 100
	}
}

func convertToProtoTask(task *models.Task) *pb.Task {
	var testCases []*pb.TestCase
	for _, test := range task.Tests {
		testCases = append(testCases, &pb.TestCase{
			Input:          test.Input,
			ExpectedOutput: test.ExpectedOutput,
			IsHidden:       test.IsHidden,
		})
	}

	return &pb.Task{
		Id:            task.ID.String(),
		Title:         task.Title,
		Description:   task.Description,
		Difficulty:    task.Defficulty,
		TestCases:     testCases,
		TimeLimitMs:   int32(getTimeLimit(task.Tests)),
		MemoryLimitMb: 256, // значение по умолчанию
		Language:      task.Language,
	}
}

func getTimeLimit(tests []models.Test) int {
	if len(tests) > 0 && tests[0].Timeout > 0 {
		return tests[0].Timeout
	}
	return 5000 // 5 секунд по умолчанию
}
