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
		Description: "Напишите программу которая принимает два числа через input() и выводит их сумму",
		Language:    "python",
		Difficulty:  "easy",
		Template: `# Введите два числа и выведите их сумму
		num1 = int(input())
		num2 = int(input())
		print(num1 + num2)`,
		Tests: []models.Test{
			{
				Input:          "5\n3",
				ExpectedOutput: "8",
			},
			{
				Input:          "10\n20",
				ExpectedOutput: "30",
			},
			{
				Input:          "-5\n8",
				ExpectedOutput: "3",
			},
		},
	},
	"python_3": {
		ID:          "3",
		Title:       "Факториал числа",
		Description: "Напишите программу которая принимает число через input() и вычисляет его факториал",
		Language:    "python",
		Difficulty:  "medium",
		Template: `# Вычислите факториал введенного числа
		n = int(input())
		factorial = 1
		for i in range(1, n + 1):
			factorial *= i
		print(factorial)`,
		Tests: []models.Test{
			{
				Input:          "5",
				ExpectedOutput: "120",
			},
			{
				Input:          "3",
				ExpectedOutput: "6",
			},
			{
				Input:          "1",
				ExpectedOutput: "1",
			},
		},
	},
	"python_4": {
		ID:          "4",
		Title:       "Проверка на чётность",
		Description: "Напишите программу которая принимает число и выводит 'чётное' если число чётное, 'нечётное' если нечётное",
		Language:    "python",
		Difficulty:  "easy",
		Template: `# Проверьте чётность числа
		num = int(input())
		if num % 2 == 0:
			print("чётное")
		else:
			print("нечётное")`,
		Tests: []models.Test{
			{
				Input:          "4",
				ExpectedOutput: "чётное",
			},
			{
				Input:          "7",
				ExpectedOutput: "нечётное",
			},
			{
				Input:          "0",
				ExpectedOutput: "чётное",
			},
			{
				Input:          "-3",
				ExpectedOutput: "нечётное",
			},
		},
	},
	"python_5": {
		ID:          "5",
		Title:       "Поиск максимума",
		Description: "Напишите программу которая принимает три числа и выводит наибольшее из них",
		Language:    "python",
		Difficulty:  "medium",
		Template: `# Найдите максимальное из трёх чисел
		a = int(input())
		b = int(input()) 
		c = int(input())
		max_num = a
		if b > max_num:
			max_num = b
		if c > max_num:
			max_num = c
		print(max_num)`,
		Tests: []models.Test{
			{
				Input:          "1\n2\n3",
				ExpectedOutput: "3",
			},
			{
				Input:          "10\n5\n8",
				ExpectedOutput: "10",
			},
			{
				Input:          "-5\n-2\n-10",
				ExpectedOutput: "-2",
			},
			{
				Input:          "7\n7\n7",
				ExpectedOutput: "7",
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

	// Если указаны язык и ID - возвращаем конкретную задачу с тестами
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

	// Если указан только язык - возвращаем все задачи для этого языка с тестами
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

	// Если параметров нет - возвращаем все задачи (без тестов или с тестами)
	var publicTasks []models.Task
	for _, task := range tasks {
		publicTasks = append(publicTasks, models.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Language:    task.Language,
			Difficulty:  task.Difficulty,
			Template:    task.Template,
			// Tests: task.Tests, // Раскомментируйте если хотите возвращать тесты
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(publicTasks)
}
