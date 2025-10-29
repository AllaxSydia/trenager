package handlers

import (
	"backend/internal/executor"
	"backend/internal/models"
	"backend/internal/services"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Короче, тут выбираем стратегию выполнения, либо Docker либо Локально
// потом преобразование результатов в единый формат ответа

var dockerService *services.DockerService // Изоляция
var localExecutor *executor.LocalExecutor // Быстро

func init() {
	var err error
	dockerService, err = services.NewDockerService()
	if err != nil {
		log.Printf("Warning: Docker service not available: %v", err)
		log.Println("Running in local execution mode")
	} else {
		log.Println("✅ Docker service initialized successfully")
	}

	// Инициализируем локальный исполнитель
	localExecutor = executor.NewLocalExecutor()
	// Если не получается создать Docker сервис для изолированного выполнения, то
	// Переходим в локальный режим
	// Локалка создаётся всегда
}

func ExecuteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, `{"success": false, "message": "Only POST method allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	// Парсинг JSON
	var req models.ExecutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success": false, "message": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	log.Printf("🔧 Executing code for language: %s", req.Language)

	var response models.ExecutionResponse

	// Пробуем Docker сначала
	if dockerService != nil {
		log.Println("🐳 Attempting Docker execution...")
		result, err := dockerService.ExecuteCode(req.Code, req.Language)
		if err != nil {
			log.Printf("❌ Docker execution failed: %v", err)
			log.Println("🔄 Falling back to local execution...")
			response = executeCodeWithLocalExecutor(req.Code, req.Language)
		} else {
			log.Printf("✅ Docker execution successful, output: %s", result.Output)
			response = models.ExecutionResponse{
				Success: result.Success,
				Message: "Code executed successfully via Docker",
				Output:  result.Output,
			}
			if !result.Success {
				response.Message = "Code execution failed in Docker"
			}
		}
	} else {
		log.Println("🔄 Docker not available, using local execution...")
		response = executeCodeWithLocalExecutor(req.Code, req.Language)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Новая функция использующая LocalExecutor для всех языков
func executeCodeWithLocalExecutor(code, language string) models.ExecutionResponse {
	log.Printf("🔧 Executing %s code with local executor", language)

	// УДАЛИ ЭТУ ПРОВЕРКУ - она блокирует Java!
	// if language == "java" {
	//     // Проверяем установлен ли Java компилятор
	//     if _, err := exec.LookPath("javac"); err != nil {
	//         return models.ExecutionResponse{
	//             Success: false,
	//             Message: "Java execution is not available. Please install Java JDK.",
	//             Output:  "",
	//         }
	//     }
	// }

	result, err := localExecutor.Execute(code, language)

	if err != nil {
		log.Printf("❌ Local execution error: %v", err)
		return models.ExecutionResponse{
			Success: false,
			Message: "Execution failed: " + err.Error(),
			Output:  "",
		}
	}

	// Преобразуем результат из LocalExecutor в models.ExecutionResponse
	exitCode := result["exitCode"].(int)
	output := result["output"].(string)
	errorMsg := result["error"].(string)

	success := exitCode == 0
	finalOutput := output
	if errorMsg != "" {
		finalOutput = errorMsg
		if output != "" {
			finalOutput = output + "\n" + errorMsg
		}
	}

	message := "Код выполнен успешно (локально)"
	if !success {
		message = "Ошибка выполнения кода"
	}

	log.Printf("✅ Local execution completed, success: %t, output length: %d", success, len(finalOutput))

	return models.ExecutionResponse{
		Success: success,
		Message: message,
		Output:  finalOutput,
	}
}

// CheckHandler - проверка решений задач
func CheckHandler(w http.ResponseWriter, r *http.Request) {
	// Логирование входящего запроса
	log.Printf("🎯 CheckHandler called: %s %s", r.Method, r.URL.Path)
	log.Printf("📋 Headers: Content-Type=%s, Content-Length=%s",
		r.Header.Get("Content-Type"),
		r.Header.Get("Content-Length"))

	// Проверка метода
	if r.Method != "POST" {
		log.Printf("❌ Method not allowed: %s", r.Method)
		http.Error(w, `{"success": false, "message": "Only POST method allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Чтение и логирование тела запроса
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("❌ Failed to read request body: %v", err)
		http.Error(w, `{"success": false, "message": "Failed to read request body"}`, http.StatusBadRequest)
		return
	}

	// Восстановление body для парсинга
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	log.Printf("📦 Raw request body (%d bytes): %s", len(bodyBytes), string(bodyBytes))

	// Парсинг JSON в интерфейс для гибкой обработки task_id
	var rawReq map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&rawReq); err != nil {
		log.Printf("❌ JSON parse error: %v", err)
		log.Printf("❌ Request body that failed parsing: %s", string(bodyBytes))
		http.Error(w, `{"success": false, "message": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Преобразуем task_id в строку независимо от типа
	taskID := ""
	if taskIDVal, exists := rawReq["task_id"]; exists {
		switch v := taskIDVal.(type) {
		case string:
			taskID = v
		case float64: // JSON числа становятся float64
			taskID = strconv.Itoa(int(v))
		case int:
			taskID = strconv.Itoa(v)
		default:
			log.Printf("❌ Unknown task_id type: %T", taskIDVal)
			http.Error(w, `{"success": false, "message": "Invalid task_id format"}`, http.StatusBadRequest)
			return
		}
	} else {
		log.Printf("❌ task_id is missing in request")
		http.Error(w, `{"success": false, "message": "task_id is required"}`, http.StatusBadRequest)
		return
	}

	// Извлекаем остальные поля
	language, _ := rawReq["language"].(string)
	code, _ := rawReq["code"].(string)

	log.Printf("🔍 Parsed request: task_id=%s, language=%s, code_length=%d",
		taskID, language, len(code))

	// Выполнение кода
	log.Printf("🚀 Starting code execution for task %s", taskID)
	var executionResult models.ExecutionResponse

	if dockerService != nil {
		log.Println("🐳 Attempting Docker execution...")
		result, err := dockerService.ExecuteCode(code, language)
		if err != nil {
			log.Printf("❌ Docker execution failed: %v", err)
			log.Println("🔄 Falling back to local execution...")
			executionResult = executeCodeWithLocalExecutor(code, language)
		} else {
			log.Printf("✅ Docker execution successful")
			executionResult = models.ExecutionResponse{
				Success: result.Success,
				Output:  result.Output,
			}
		}
	} else {
		log.Println("🔄 Docker not available, using local execution...")
		executionResult = executeCodeWithLocalExecutor(code, language)
	}

	log.Printf("📊 Execution result: success=%t, output_length=%d",
		executionResult.Success, len(executionResult.Output))

	// Проверка решения
	log.Printf("🧪 Checking solution against test cases")
	checkResult := checkSolution(taskID, executionResult.Output, language)

	// Формирование ответа
	response := models.CheckResponse{
		Success:  executionResult.Success && checkResult.Passed,
		Passed:   checkResult.Passed,
		Output:   executionResult.Output,
		Expected: checkResult.Expected,
		Actual:   checkResult.Actual,
		Message:  checkResult.Message,
	}

	log.Printf("✅ Check completed: passed=%t, message=%s", response.Passed, response.Message)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("❌ Failed to encode response: %v", err)
		http.Error(w, `{"success": false, "message": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("📤 Response sent successfully")
}

// checkSolution - проверяет вывод кода против ожидаемого результата
func checkSolution(taskID, actualOutput, language string) models.CheckResult {
	log.Printf("🔎 Checking solution for task=%s, language=%s", taskID, language)

	// Простая заглушка - в реальности брать из БД или из models.Task
	testCases := map[string]map[string]string{
		"1": {
			"python":     "Hello World",
			"javascript": "Hello World",
			"cpp":        "Hello World",
			"java":       "Hello World",
		},
		"hello_world": {
			"python":     "Hello World",
			"javascript": "Hello World",
			"cpp":        "Hello World",
			"java":       "Hello World",
		},
	}

	expected := "Hello World"
	if cases, exists := testCases[taskID]; exists {
		if exp, exists := cases[language]; exists {
			expected = exp
		}
	}

	// Улучшенная проверка - ищем ожидаемую строку в выводе
	actualTrimmed := strings.TrimSpace(actualOutput)
	passed := false

	// Для задачи "Hello World" проверяем что вывод содержит ожидаемую строку
	if taskID == "1" || taskID == "hello_world" {
		passed = strings.Contains(actualTrimmed, expected)
	} else {
		// Для других задач строгое сравнение
		passed = actualTrimmed == expected
	}

	log.Printf("📊 Test comparison: expected='%s', actual='%s', passed=%t",
		expected, actualTrimmed, passed)

	message := "✅ Тест пройден!"
	if !passed {
		message = "❌ Тест не пройден"
	}

	return models.CheckResult{
		Passed:   passed,
		Expected: expected,
		Actual:   actualTrimmed,
		Message:  message,
	}
}
