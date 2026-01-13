package handlers

import (
	"backend/internal/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// TaskHandler обрабатывает запросы связанные с задачами
type TaskHandler struct {
	DB *sql.DB
}

// NewTaskHandler создает новый экземпляр TaskHandler
func NewTaskHandler(db *sql.DB) *TaskHandler {
	return &TaskHandler{DB: db}
}

// GetTasksHandler возвращает задачи (публичный доступ)
func (h *TaskHandler) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры языка и ID из query string
	language := r.URL.Query().Get("language") // Изменил с "lang" на "language"
	taskID := r.URL.Query().Get("id")

	// Если указаны язык и ID - возвращаем конкретную задачу
	if language != "" && taskID != "" {
		h.getTaskByLanguageAndID(w, language, taskID)
		return
	}

	// Если указан только язык - возвращаем все задачи для этого языка
	if language != "" {
		h.getTasksByLanguage(w, language)
		return
	}

	// Если параметров нет - возвращаем все задачи
	h.getAllTasks(w)
}

// getTaskByLanguageAndID возвращает конкретную задачу по языку и ID
func (h *TaskHandler) getTaskByLanguageAndID(w http.ResponseWriter, language, taskID string) {
	// Ищем задачу в БД
	query := `
		SELECT id::text, title, description, language,
            COALESCE(template, starter_code) as template,
            starter_code, tests, created_at, updated_at
    	FROM tasks
    	WHERE language = $1 AND id::text = $2 AND is_published = true
	`

	var task models.Task
	var testsJSON []byte
	var createdAt, updatedAt time.Time
	var starterCode, template sql.NullString

	err := h.DB.QueryRow(query, language, taskID).Scan(
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
		// Если не найдено в БД, возвращаем встроенные задачи
		h.getBuiltInTask(w, language, taskID)
		return
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
		http.Error(w, "Error parsing tests", http.StatusInternalServerError)
		return
	}

	// Добавляем метаданные
	task.CreatedAt = createdAt
	task.UpdatedAt = updatedAt

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// getBuiltInTask возвращает встроенную задачу
func (h *TaskHandler) getBuiltInTask(w http.ResponseWriter, language, taskID string) {
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
		return
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

// getTasksByLanguage возвращает задачи по языку
func (h *TaskHandler) getTasksByLanguage(w http.ResponseWriter, language string) {
	var allTasks []models.Task

	// Сначала берем задачи из БД
	query := `
        SELECT id::text, title, description, language, 
               COALESCE(template, starter_code) as template,
               starter_code, tests, created_at, updated_at,
               is_published, created_by
        FROM tasks 
        WHERE language = $1 AND is_published = true
        ORDER BY created_at DESC
    `

	rows, err := h.DB.Query(query, language)
	if err != nil {
		log.Printf("❌ Ошибка запроса задач из БД: %v", err)
		// Продолжаем - вернем только встроенные задачи
	} else {
		defer rows.Close()

		for rows.Next() {
			var task models.Task
			var testsJSON []byte
			var createdAt, updatedAt time.Time
			var starterCode, template string
			var isPublished bool
			var createdBy sql.NullInt64

			err := rows.Scan(
				&task.ID,
				&task.Title,
				&task.Description,
				&task.Language,
				&template,
				&starterCode,
				&testsJSON,
				&createdAt,
				&updatedAt,
				&isPublished,
				&createdBy,
			)

			if err != nil {
				log.Printf("⚠️ Ошибка сканирования задачи: %v", err)
				continue
			}

			task.Template = template
			task.StarterCode = starterCode

			// Парсим тесты
			if len(testsJSON) > 0 {
				if err := json.Unmarshal(testsJSON, &task.Tests); err != nil {
					log.Printf("⚠️ Ошибка парсинга тестов: %v", err)
					task.Tests = []models.Test{}
				}
			}

			task.CreatedAt = createdAt
			task.UpdatedAt = updatedAt
			allTasks = append(allTasks, task)
		}

		log.Printf("✅ Загружено %d задач из БД для языка %s", len(allTasks), language)
	}

	// Дополняем встроенными задачами если БД пустая
	if len(allTasks) == 0 {
		log.Printf("⚠️ В БД нет задач для языка %s, используем встроенные", language)
		allTasks = h.getBuiltInTasksByLanguage(language)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allTasks); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// getBuiltInTasksByLanguage возвращает встроенные задачи по языку
func (h *TaskHandler) getBuiltInTasksByLanguage(language string) []models.Task {
	var tasks []models.Task

	builtInTasks := []models.Task{
		{
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
		{
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
		{
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

	for _, task := range builtInTasks {
		if task.Language == language {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

// getAllTasks возвращает все задачи
func (h *TaskHandler) getAllTasks(w http.ResponseWriter) {
	// Сначала собираем встроенные задачи
	var allTasks []models.Task
	allTasks = append(allTasks, h.getBuiltInTasksByLanguage("python")...)

	// Дополняем задачами из БД
	query := `
		SELECT id::text, title, description, language, template, 
		       starter_code, tests, created_at, updated_at
		FROM tasks 
		WHERE is_published = true
		ORDER BY language, created_at DESC
	`

	rows, err := h.DB.Query(query)
	if err != nil {
		// Если ошибка БД, возвращаем только встроенные задачи
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allTasks)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		var testsJSON []byte
		var createdAt, updatedAt time.Time
		var starterCode, template sql.NullString

		err := rows.Scan(
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
			continue
		}

		// Заполняем опциональные поля
		if template.Valid {
			task.Template = template.String
		}
		if starterCode.Valid {
			task.StarterCode = starterCode.String
		}

		// Парсим тесты
		if err := json.Unmarshal(testsJSON, &task.Tests); err == nil {
			task.CreatedAt = createdAt
			task.UpdatedAt = updatedAt
			allTasks = append(allTasks, task)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allTasks)
}

// CreateTaskHandler создает новую задачу (только для учителей)
func (h *TaskHandler) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем авторизацию и роль учителя
	userID, role, err := h.getUserFromRequest(r)
	if err != nil || role != "teacher" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	var taskReq models.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if taskReq.Title == "" || taskReq.Description == "" || taskReq.Language == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if len(taskReq.Tests) == 0 {
		http.Error(w, "At least one test is required", http.StatusBadRequest)
		return
	}

	// Конвертируем тесты в JSON
	testsJSON, err := json.Marshal(taskReq.Tests)
	if err != nil {
		http.Error(w, "Error processing tests", http.StatusInternalServerError)
		return
	}

	// Вставляем в БД
	query := `
		INSERT INTO tasks (
			title, description, language, difficulty, template, starter_code,
			tests, created_by, created_at, updated_at, is_published
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	now := time.Now()
	var taskID int
	err = h.DB.QueryRow(
		query,
		taskReq.Title,
		taskReq.Description,
		taskReq.Language,
		taskReq.Difficulty,
		taskReq.Template,
		taskReq.StarterCode,
		testsJSON,
		userID,
		now,
		now,
		true, // is_published
	).Scan(&taskID)

	if err != nil {
		http.Error(w, "Error creating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем созданную задачу
	response := map[string]interface{}{
		"id":      strconv.Itoa(taskID),
		"message": "Task created successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetTeacherTasksHandler возвращает задачи созданные учителем
func (h *TaskHandler) GetTeacherTasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем авторизацию
	userID, role, err := h.getUserFromRequest(r)
	if err != nil || role != "teacher" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	log.Printf("🔍 Загрузка задач для учителя ID: %d", userID)

	query := `
        SELECT id::text, title, description, language, 
               COALESCE(template, starter_code) as template,
               starter_code, tests, created_at, updated_at, 
               is_published
        FROM tasks 
        WHERE created_by = $1
        ORDER BY created_at DESC
    `

	rows, err := h.DB.Query(query, userID)
	if err != nil {
		log.Printf("❌ Ошибка запроса задач учителя: %v", err)
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var teacherTasks []models.Task
	for rows.Next() {
		var task models.Task
		var testsJSON []byte
		var createdAt, updatedAt time.Time
		var starterCode, template string
		var isPublished bool

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Language,
			&template,
			&starterCode,
			&testsJSON,
			&createdAt,
			&updatedAt,
			&isPublished,
		)

		if err != nil {
			log.Printf("⚠️ Ошибка сканирования задачи учителя: %v", err)
			continue
		}

		task.Template = template
		task.StarterCode = starterCode

		// Парсим тесты
		if len(testsJSON) > 0 {
			if err := json.Unmarshal(testsJSON, &task.Tests); err != nil {
				log.Printf("⚠️ Ошибка парсинга тестов учителя: %v", err)
				task.Tests = []models.Test{}
			}
		}

		task.CreatedAt = createdAt
		task.UpdatedAt = updatedAt
		teacherTasks = append(teacherTasks, task)
	}

	log.Printf("✅ Найдено %d задач для учителя ID: %d", len(teacherTasks), userID)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(teacherTasks); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// getUserFromRequest извлекает данные пользователя из запроса
func (h *TaskHandler) getUserFromRequest(r *http.Request) (int, string, error) {
	// Здесь должна быть проверка JWT токена
	// Пока возвращаем тестовые данные
	return 1, "teacher", nil
}

// UpdateTaskHandler обновляет существующую задачу
func (h *TaskHandler) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем авторизацию и роль учителя
	userID, role, err := h.getUserFromRequest(r)
	if err != nil || role != "teacher" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Получаем ID задачи из query параметров
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	var taskReq models.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Валидация
	if taskReq.Title == "" || taskReq.Description == "" || taskReq.Language == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if len(taskReq.Tests) == 0 {
		http.Error(w, "At least one test is required", http.StatusBadRequest)
		return
	}

	// Проверяем, принадлежит ли задача этому учителю
	var createdBy int
	err = h.DB.QueryRow(
		"SELECT created_by FROM tasks WHERE id::text = $1",
		taskID,
	).Scan(&createdBy)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if createdBy != userID {
		http.Error(w, "You can only update your own tasks", http.StatusForbidden)
		return
	}

	// Конвертируем тесты в JSON
	testsJSON, err := json.Marshal(taskReq.Tests)
	if err != nil {
		http.Error(w, "Error processing tests", http.StatusInternalServerError)
		return
	}

	// Обновляем задачу в БД
	query := `
		UPDATE tasks 
		SET 
			title = $1,
			description = $2,
			language = $3,
			difficulty = $4,
			template = $5,
			starter_code = $6,
			tests = $7,
			updated_at = $8,
			is_published = $9
		WHERE id::text = $10 AND created_by = $11
		RETURNING id
	`

	now := time.Now()
	var updatedID int
	err = h.DB.QueryRow(
		query,
		taskReq.Title,
		taskReq.Description,
		taskReq.Language,
		taskReq.Difficulty,
		taskReq.Template,
		taskReq.StarterCode,
		testsJSON,
		now,
		taskReq.IsPublished,
		taskID,
		userID,
	).Scan(&updatedID)

	if err != nil {
		http.Error(w, "Error updating task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем успешный ответ
	response := map[string]interface{}{
		"id":      strconv.Itoa(updatedID),
		"message": "Task updated successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteTaskHandler удаляет задачу
func (h *TaskHandler) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем авторизацию и роль учителя
	userID, role, err := h.getUserFromRequest(r)
	if err != nil || role != "teacher" {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	// Получаем ID задачи из query параметров
	taskID := r.URL.Query().Get("id")
	if taskID == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	// Проверяем, принадлежит ли задача этому учителю
	var createdBy int
	err = h.DB.QueryRow(
		"SELECT created_by FROM tasks WHERE id::text = $1",
		taskID,
	).Scan(&createdBy)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error", http.StatusInternalServerError)
		}
		return
	}

	if createdBy != userID {
		http.Error(w, "You can only delete your own tasks", http.StatusForbidden)
		return
	}

	// Удаляем задачу из БД
	result, err := h.DB.Exec(
		"DELETE FROM tasks WHERE id::text = $1 AND created_by = $2",
		taskID,
		userID,
	)

	if err != nil {
		http.Error(w, "Error deleting task: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Проверяем, была ли удалена хотя бы одна строка
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking deletion result", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Task not found or already deleted", http.StatusNotFound)
		return
	}

	// Возвращаем успешный ответ
	response := map[string]interface{}{
		"message": "Task deleted successfully",
		"task_id": taskID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
