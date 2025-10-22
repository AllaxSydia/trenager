package handlers

import (
	"backend/internal/models"
	"encoding/json"
	"net/http"
)

// Библиотека задач
// Временное хранилище задач в памяти
var tasks = []models.Task{
	{
		ID:          "1",
		Title:       "Hello World",
		Description: "Напишите программу которая выводит 'Hello, World!'",
		Template:    "print('Hello, World!')",
		Tests: []models.Test{
			{
				Input:          "",
				ExpectedOutput: "Hello, World!",
			},
		},
	},
	{
		ID:          "2",
		Title:       "Сумма двух чисел",
		Description: "Напишите функцию sum(a, b) которая возвращает сумму двух чисел",
		Template:    "def sum(a, b):\n    # Ваш код здесь\n    pass\n\n# Тестирование\nresult = sum(2, 3)\nprint(result)",
		Tests: []models.Test{
			{
				Input:          "2, 3",
				ExpectedOutput: "5",
			},
			{
				Input:          "10, -5",
				ExpectedOutput: "5",
			},
		},
	},
	{
		ID:          "3",
		Title:       "Факториал",
		Description: "Напишите функцию для вычисления факториала числа",
		Template:    "def factorial(n):\n    # Ваш код здесь\n    pass\n\n# Тестирование\nprint(factorial(5))",
		Tests: []models.Test{
			{
				Input:          "5",
				ExpectedOutput: "120",
			},
		},
	},
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Возвращаем задачи без тестов (для безопасности)
	var publicTasks []models.Task
	for _, task := range tasks {
		publicTasks = append(publicTasks, models.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Template:    task.Template,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publicTasks)
}
