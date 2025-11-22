package handlers

import (
	"backend/internal/database"
	"backend/internal/models"
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

	// Получаем задачу для проверки
	taskKey := req.Language + "_" + taskID
	task, exists := tasks[taskKey]
	if !exists {
		http.Error(w, `{"success": false, "message": "Task not found"}`, http.StatusNotFound)
		return
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
			inputs = []string{test.Input}
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
