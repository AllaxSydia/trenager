package handlers

import (
	"backend/internal/models"
	"encoding/json"
	"net/http"
	"strings"
)

// Библиотека задач с поддержкой разных языков
var tasks = map[string]models.Task{
	"python_1": {
		ID:          "1",
		Title:       "Hello World",
		Description: "Напишите программу которая выводит 'Hello, World!'",
		Language:    "python",
		Difficulty:  "easy",
		Template:    "print('Hello, World!')",
		Tests: []models.Test{
			{
				Input:          "",
				ExpectedOutput: "Hello, World!",
			},
		},
	},
	"python_2": {
		ID:          "2",
		Title:       "Сумма двух чисел",
		Description: "Напишите программу которая складывает два числа (5 и 3)",
		Language:    "python",
		Difficulty:  "easy",
		Template:    "# Сложите два числа\nnum1 = 5\nnum2 = 3\nresult = num1 + num2\nprint(result)",
		Tests: []models.Test{
			{
				Input:          "",
				ExpectedOutput: "8",
			},
		},
	},
	"python_3": {
		ID:          "3",
		Title:       "Факториал",
		Description: "Напишите программу которая вычисляет факториал числа 5",
		Language:    "python",
		Difficulty:  "medium",
		Template:    "# Вычислите факториал числа 5\nn = 5\nfactorial = 1\nfor i in range(1, n + 1):\n    factorial *= i\nprint(factorial)",
		Tests: []models.Test{
			{
				Input:          "",
				ExpectedOutput: "120",
			},
		},
	},
	"javascript_1": {
		ID:          "1",
		Title:       "Hello World",
		Description: "Напишите программу которая выводит 'Hello, World!'",
		Language:    "javascript",
		Difficulty:  "easy",
		Template:    "console.log('Hello, World!');",
		Tests: []models.Test{
			{
				Input:          "",
				ExpectedOutput: "Hello, World!",
			},
		},
	},
	"javascript_2": {
		ID:          "2",
		Title:       "Сумма двух чисел",
		Description: "Напишите программу которая складывает два числа (5 и 3)",
		Language:    "javascript",
		Difficulty:  "easy",
		Template:    "// Сложите два числа\nconst num1 = 5;\nconst num2 = 3;\nconst result = num1 + num2;\nconsole.log(result);",
		Tests: []models.Test{
			{
				Input:          "",
				ExpectedOutput: "8",
			},
		},
	},
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры языка и ID из query string
	language := r.URL.Query().Get("lang")
	taskID := r.URL.Query().Get("id")

	// Если указаны язык и ID - возвращаем конкретную задачу
	if language != "" && taskID != "" {
		key := language + "_" + taskID
		if task, exists := tasks[key]; exists {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		} else {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
	}

	// Если указан только язык - возвращаем все задачи для этого языка
	if language != "" {
		var languageTasks []models.Task
		for key, task := range tasks {
			if strings.HasPrefix(key, language+"_") {
				languageTasks = append(languageTasks, task)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(languageTasks)
		return
	}

	// Если параметров нет - возвращаем все задачи (без тестов для безопасности)
	var publicTasks []models.Task
	for _, task := range tasks {
		publicTasks = append(publicTasks, models.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Language:    task.Language,
			Difficulty:  task.Difficulty,
			Template:    task.Template,
			// Tests не включаем для безопасности
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publicTasks)
}
