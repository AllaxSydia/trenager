package handlers

import (
	"backend/internal/database"
	"backend/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// CheckHandler обрабатывает проверку кода на соответствие тестам
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"success": false, "message": "Only POST method allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Парсинг JSON запроса
	var req models.CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("❌ Failed to parse check request: %v", err)
		http.Error(w, `{"success": false, "message": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Валидация
	if req.Code == "" {
		http.Error(w, `{"success": false, "message": "Code is required"}`, http.StatusBadRequest)
		return
	}
	if req.Language == "" {
		http.Error(w, `{"success": false, "message": "Language is required"}`, http.StatusBadRequest)
		return
	}
	if req.TaskID == nil {
		http.Error(w, `{"success": false, "message": "Task ID is required"}`, http.StatusBadRequest)
		return
	}

	// Конвертируем TaskID в строку
	taskID := convertTaskIDToString(req.TaskID)

	log.Printf("🔍 Checking code for task: %s, language: %s, code length: %d", taskID, req.Language, len(req.Code))

	// Получаем задачу из БД
	task, err := getTaskFromDB(req.Language, taskID)
	if err != nil {
		// Если не найдено в БД, пробуем получить встроенную задачу
		task = getBuiltInTask(req.Language, taskID)
		if task.ID == "" {
			http.Error(w, `{"success": false, "message": "Task not found"}`, http.StatusNotFound)
			return
		}
	}

	// Используем тесты из задачи, если не предоставлены в запросе
	testsToRun := task.Tests
	if len(req.Tests) > 0 {
		testsToRun = req.Tests
	}

	if len(testsToRun) == 0 {
		http.Error(w, `{"success": false, "message": "No tests available for this task"}`, http.StatusBadRequest)
		return
	}

	// Выполняем код с тестами
	allTestsPassed := true
	var testResults []models.TestResult

	for i, test := range testsToRun {
		// Подготавливаем входные данные если есть
		var inputs []string
		if test.Input != "" {
			// Разделяем многострочный ввод
			inputs = strings.Split(test.Input, "\n")
		}

		// Выполняем код с текущим тестом
		result, err := codeExecutor.Execute(req.Code, req.Language, inputs)
		if err != nil {
			log.Printf("❌ Test %d execution error: %v", i+1, err)
			allTestsPassed = false
			testResults = append(testResults, models.TestResult{
				TestNumber: i + 1,
				Passed:     false,
				Output:     "",
				Error:      err.Error(),
				Expected:   test.ExpectedOutput,
				Actual:     "",
			})
			continue
		}

		// Получаем вывод
		output := ""
		if out, ok := result["output"].(string); ok {
			output = out
		}

		// Нормализуем вывод для сравнения
		normalizedOutput := normalizeOutput(output)
		normalizedExpected := normalizeOutput(test.ExpectedOutput)

		// Проверяем соответствие ожидаемому результату
		passed := normalizedOutput == normalizedExpected

		if !passed {
			allTestsPassed = false
		}

		testResults = append(testResults, models.TestResult{
			TestNumber: i + 1,
			Passed:     passed,
			Output:     output,
			Expected:   test.ExpectedOutput,
			Actual:     normalizedOutput,
		})

		log.Printf("🧪 Test %d: passed=%t, output='%s', expected='%s'",
			i+1, passed, normalizedOutput, normalizedExpected)
	}

	// Формируем ответ
	message := "✅ Все тесты пройдены!"
	if !allTestsPassed {
		message = "❌ Некоторые тесты не пройдены"
	}

	response := models.CheckResponse{
		Success:     allTestsPassed,
		Message:     message,
		TestResults: testResults,
		TotalTests:  len(testsToRun),
		PassedTests: countPassedTests(testResults),
	}

	log.Printf("📊 Check completed - Success: %t, Passed: %d/%d",
		allTestsPassed, response.PassedTests, response.TotalTests)

	// Сохраняем решение в БД, если пользователь авторизован
	auth := r.Header.Get("Authorization")
	if auth != "" {
		parts := strings.Fields(auth)
		if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
			// Пытаемся получить user_id из токена
			claims, err := ParseTokenFromRequest(r)
			if err == nil {
				if userIDFloat, ok := claims["sub"].(float64); ok {
					userID := int64(userIDFloat)
					saveTaskSolution(userID, taskID, req.Language, req.Code, allTestsPassed, response.PassedTests, response.TotalTests)
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("❌ Failed to encode check response: %v", err)
		http.Error(w, `{"success": false, "message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}
}

// getTaskFromDB получает задачу из базы данных
func getTaskFromDB(language, taskID string) (models.Task, error) {
	var task models.Task

	query := `
		SELECT id::text, title, description, language, template, 
		       starter_code, tests, created_at, updated_at
		FROM tasks 
		WHERE language = $1 AND id::text = $2 AND is_published = true
	`

	var testsJSON []byte
	var createdAt, updatedAt string // Используем string для временных меток
	var starterCode, template sql.NullString

	err := database.DB.QueryRow(query, language, taskID).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Language,
		&template,
		&starterCode,
		&testsJSON,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return task, err
	}

	// Заполняем опциональные поля
	if template.Valid {
		task.Template = template.String
	}
	if starterCode.Valid {
		task.StarterCode = starterCode.String
	}

	// Парсим тесты
	if err := json.Unmarshal(testsJSON, &task.Tests); err != nil {
		log.Printf("Error parsing tests JSON: %v", err)
		// Возвращаем задачу без тестов
		task.Tests = []models.Test{}
	}

	return task, nil
}

// getBuiltInTask возвращает встроенную задачу
func getBuiltInTask(language, taskID string) models.Task {
	// Встроенные задачи для обратной совместимости
	builtInTasks := map[string]models.Task{
		"python_1": {
			ID:          "1",
			Title:       "Hello World",
			Description: "Напишите программу которая выводит 'Hello, World!'",
			Language:    "python",
			Template:    `print("Hello, World!")`,
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
			Template: `num1 = int(input())
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
			},
		},
		"python_3": {
			ID:          "3",
			Title:       "Факториал",
			Description: "Напишите функцию для вычисления факториала числа",
			Language:    "python",
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
			},
		},
	}

	key := language + "_" + taskID
	if task, exists := builtInTasks[key]; exists {
		return task
	}

	return models.Task{} // Пустая задача если не найдено
}

// saveTaskSolution сохраняет решение задачи в БД
func saveTaskSolution(userID int64, taskID, language, code string, success bool, passedTests, totalTests int) {
	query := `
	INSERT INTO task_solutions (user_id, task_id, language, code, success, passed_tests, total_tests)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (user_id, task_id, language) 
	DO UPDATE SET 
		code = EXCLUDED.code,
		success = EXCLUDED.success,
		passed_tests = EXCLUDED.passed_tests,
		total_tests = EXCLUDED.total_tests,
		created_at = CURRENT_TIMESTAMP
	`
	_, err := database.DB.Exec(query, userID, taskID, language, code, success, passedTests, totalTests)
	if err != nil {
		log.Printf("⚠️ Ошибка при сохранении решения задачи: %v", err)
	} else {
		log.Printf("✅ Решение задачи сохранено: user_id=%d, task_id=%s, language=%s", userID, taskID, language)
	}
}

// convertTaskIDToString конвертирует TaskID в строку
func convertTaskIDToString(taskID interface{}) string {
	switch v := taskID.(type) {
	case string:
		return v
	case float64:
		return strconv.Itoa(int(v))
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

// normalizeOutput нормализует вывод для сравнения
func normalizeOutput(output string) string {
	// Убираем начальные и конечные пробелы, переводы строк
	return strings.TrimSpace(output)
}

// countPassedTests подсчитывает количество пройденных тестов
func countPassedTests(results []models.TestResult) int {
	count := 0
	for _, result := range results {
		if result.Passed {
			count++
		}
	}
	return count
}
