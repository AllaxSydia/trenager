package services

import (
	"backend/internal/models"
)

// TaskService сервис для работы с задачами
type TaskService struct {
	tasks map[string]*models.Task
}

func NewTaskService() *TaskService {
	service := &TaskService{
		tasks: make(map[string]*models.Task),
	}

	// Инициализируем тестовые задачи
	service.initializeTasks()

	return service
}

// GetTasks возвращает все задачи
func (s *TaskService) GetTasks() []*models.Task {
	var tasks []*models.Task
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetTaskByID возвращает задачу по ID
func (s *TaskService) GetTaskByID(id string) (*models.Task, error) {
	task, exists := s.tasks[id]
	if !exists {
		return nil, nil
	}
	return task, nil
}

// GetTemplateForLanguage возвращает шаблон кода для конкретного языка
func (s *TaskService) GetTemplateForLanguage(taskID, language string) string {
	task, exists := s.tasks[taskID]
	if !exists {
		return s.getDefaultTemplate(language)
	}

	// Если задача уже содержит шаблон на Python, адаптируем его для других языков
	if language == "python" {
		return task.Template
	}

	// Адаптируем шаблон для других языков
	switch taskID {
	case "1": // Hello World
		switch language {
		case "javascript":
			return `console.log("Hello, World!");`
		case "java":
			return `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}`
		case "cpp":
			return `#include <iostream>
using namespace std;

int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}`
		case "go":
			return `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`
		default:
			return task.Template
		}

	case "2": // Сумма двух чисел
		switch language {
		case "javascript":
			return `function sum(a, b) {
    return a + b;
}

// Тестирование
console.log(sum(2, 3));`
		case "java":
			return `public class Main {
    public static int sum(int a, int b) {
        return a + b;
    }
    
    public static void main(String[] args) {
        System.out.println(sum(2, 3));
    }
}`
		case "cpp":
			return `#include <iostream>
using namespace std;

int sum(int a, int b) {
    return a + b;
}

int main() {
    std::cout << sum(2, 3) << std::endl;
    return 0;
}`
		case "go":
			return `package main

import "fmt"

func sum(a, b int) int {
    return a + b
}

func main() {
    fmt.Println(sum(2, 3))
}`
		default:
			return task.Template
		}

	case "3": // Факториал
		switch language {
		case "javascript":
			return `function factorial(n) {
    if (n === 0) return 1;
    let result = 1;
    for (let i = 1; i <= n; i++) {
        result *= i;
    }
    return result;
}

// Тестирование
console.log(factorial(5));`
		case "java":
			return `public class Main {
    public static int factorial(int n) {
        if (n == 0) return 1;
        int result = 1;
        for (int i = 1; i <= n; i++) {
            result *= i;
        }
        return result;
    }
    
    public static void main(String[] args) {
        System.out.println(factorial(5));
    }
}`
		case "cpp":
			return `#include <iostream>
using namespace std;

int factorial(int n) {
    if (n == 0) return 1;
    int result = 1;
    for (int i = 1; i <= n; i++) {
        result *= i;
    }
    return result;
}

int main() {
    std::cout << factorial(5) << std::endl;
    return 0;
}`
		case "go":
			return `package main

import "fmt"

func factorial(n int) int {
    if n == 0 {
        return 1
    }
    result := 1
    for i := 1; i <= n; i++ {
        result *= i
    }
    return result
}

func main() {
    fmt.Println(factorial(5))
}`
		default:
			return task.Template
		}
	}

	return s.getDefaultTemplate(language)
}

func (s *TaskService) getDefaultTemplate(language string) string {
	switch language {
	case "python":
		return `print("Hello, World!")`
	case "javascript":
		return `console.log("Hello, World!");`
	case "java":
		return `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}`
	case "cpp":
		return `#include <iostream>
using namespace std;

int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}`
	case "go":
		return `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`
	default:
		return `print("Hello, World!")`
	}
}

// initializeTasks инициализирует тестовые задачи
func (s *TaskService) initializeTasks() {
	// Создаем задачи напрямую в мапе
	s.tasks["1"] = &models.Task{
		ID:          "1",
		Title:       "Hello World",
		Description: "Напишите программу которая выводит 'Hello, World!'",
		Template:    `print("Hello, World!")`,
		Tests: []models.Test{
			{
				Input:          "",
				ExpectedOutput: "Hello, World!",
			},
		},
	}

	s.tasks["2"] = &models.Task{
		ID:          "2",
		Title:       "Сумма двух чисел",
		Description: "Напишите функцию sum(a, b) которая возвращает сумму двух чисел",
		Template: `def sum(a, b):
    return a + b

# Тестирование
result = sum(2, 3)
print(result)`,
		Tests: []models.Test{
			{
				Input:          "2, 3",
				ExpectedOutput: "5",
			},
			{
				Input:          "10, -5",
				ExpectedOutput: "5",
			},
			{
				Input:          "0, 0",
				ExpectedOutput: "0",
			},
		},
	}

	s.tasks["3"] = &models.Task{
		ID:          "3",
		Title:       "Факториал",
		Description: "Напишите функцию для вычисления факториала числа",
		Template: `def factorial(n):
    if n == 0:
        return 1
    result = 1
    for i in range(1, n + 1):
        result *= i
    return result

# Тестирование
print(factorial(5))`,
		Tests: []models.Test{
			{
				Input:          "5",
				ExpectedOutput: "120",
			},
			{
				Input:          "0",
				ExpectedOutput: "1",
			},
			{
				Input:          "1",
				ExpectedOutput: "1",
			},
		},
	}
}
