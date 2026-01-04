package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type CodeReviewRequest struct {
	Code        string `json:"code"`
	Language    string `json:"language"`
	TaskContext string `json:"task_context"`
	UserId      int64  `json:"user_id"`
}

type CodeReviewResponse struct {
	Score         int      `json:"score"`          // 1-10
	Comments      []string `json:"comments"`       // конкретные замечания
	Suggestions   []string `json:"suggestions"`    // предложения по улучшению
	BestPractices []string `json:"best_practices"` // лучшие практики
	Complexity    string   `json:"complexity"`     // низкая/средняя/высокая
	AIResponse    string   `json:"ai_response"`    // полный ответ ИИ
}

type AIReviewer struct {
	client       interface{}
	useSmartMock bool
	apiKey       string
	apiURL       string
	model        string
	apiType      string
}

// Вспомогательная функция для маскировки ключей в логах
func maskKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "***" + key[len(key)-4:]
}

func NewAIReviewer() *AIReviewer {
	// Сначала проверяем OpenRouter - это наш приоритет
	openrouterKey := strings.TrimSpace(os.Getenv("OPENROUTER_API_KEY"))

	// Получаем модель
	model := strings.TrimSpace(os.Getenv("OPENROUTER_MODEL"))
	if model == "" {
		model = strings.TrimSpace(os.Getenv("AI_MODEL"))
		if model == "" {
			model = "deepseek/deepseek-chat"
		}
	}

	// Подробное логирование для диагностики
	fmt.Println("🔍 ==================== AI REVIEWER INIT ====================")
	fmt.Println("📋 Checking environment variables...")

	// Проверяем ВСЕ переменные окружения
	allEnvVars := []string{
		"OPENROUTER_API_KEY",
		"OPENROUTER_MODEL",
		"AI_MODEL",
		"OPENAI_API_KEY",
		"DEEPSEEK_API_KEY",
		"DB_HOST", // для проверки что env вообще загружается
	}

	for _, envVar := range allEnvVars {
		val := os.Getenv(envVar)
		if val != "" {
			if strings.Contains(envVar, "KEY") {
				fmt.Printf("   ✅ %s: %s (length: %d)\n", envVar, maskKey(val), len(val))
			} else {
				fmt.Printf("   ✅ %s: %s\n", envVar, val)
			}
		} else {
			fmt.Printf("   ❌ %s: NOT SET\n", envVar)
		}
	}

	fmt.Println("📊 Decision making...")

	// Приоритет 1: OpenRouter
	if openrouterKey != "" && openrouterKey != "My-secret-key-openrouter-ai" {
		fmt.Printf("✅ Using OpenRouter API with model: %s\n", model)
		fmt.Printf("   API Key: %s\n", maskKey(openrouterKey))

		return &AIReviewer{
			client: &http.Client{
				Timeout: 45 * time.Second,
			},
			useSmartMock: false,
			apiKey:       openrouterKey,
			apiURL:       "https://openrouter.ai/api/v1/chat/completions",
			model:        model,
			apiType:      "openrouter",
		}
	} else if openrouterKey == "My-secret-key-openrouter-ai" {
		fmt.Println("⚠️  OpenRouter key is placeholder, checking other options...")
	}

	// Проверяем другие API как fallback
	openaiKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	deepseekKey := strings.TrimSpace(os.Getenv("DEEPSEEK_API_KEY"))

	// Приоритет 2: OpenAI
	if openaiKey != "" {
		fmt.Println("✅ Using OpenAI API (OpenRouter not properly configured)")
		fmt.Printf("   Model: gpt-3.5-turbo (default)\n")

		client := openai.NewClient(openaiKey)
		return &AIReviewer{
			client:       client,
			useSmartMock: false,
			apiKey:       openaiKey,
			apiURL:       "https://api.openai.com/v1/chat/completions",
			model:        "gpt-3.5-turbo",
			apiType:      "openai",
		}
	}

	// Приоритет 3: DeepSeek
	if deepseekKey != "" {
		fmt.Println("✅ Using DeepSeek API")

		return &AIReviewer{
			client: &http.Client{
				Timeout: 45 * time.Second,
			},
			useSmartMock: false,
			apiKey:       deepseekKey,
			apiURL:       "https://api.deepseek.com/v1/chat/completions",
			model:        "deepseek-chat",
			apiType:      "deepseek",
		}
	}

	// Приоритет 4: Smart mock
	fmt.Println("🔧 Using SMART MOCK MODE (no valid API keys found)")
	fmt.Println("   Add OPENROUTER_API_KEY to .env file for real AI")

	return &AIReviewer{
		client:       nil,
		useSmartMock: true,
		apiType:      "mock",
	}
}

func (r *AIReviewer) ReviewCode(req CodeReviewRequest) (*CodeReviewResponse, error) {
	fmt.Printf("\n🎯 ==================== AI REVIEW REQUEST ====================\n")
	fmt.Printf("📝 Language: %s\n", req.Language)
	fmt.Printf("📏 Code length: %d characters\n", len(req.Code))
	fmt.Printf("📋 Task context: %s\n", req.TaskContext)
	fmt.Printf("🤖 AI Mode: %s\n", r.apiType)

	// Если нет реального API, используем smart mock
	if r.useSmartMock || r.client == nil {
		fmt.Println("🔧 Using smart mock reviewer")
		return r.smartMockReview(req)
	}

	fmt.Printf("🚀 Calling %s API with model: %s\n", r.apiType, r.model)

	// Формируем промпт
	prompt := r.buildPrompt(req.Code, req.Language, req.TaskContext)

	// Добавляем таймаут для безопасности
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var aiResponse string
	var err error

	// Выбираем метод вызова
	switch r.apiType {
	case "openai":
		aiResponse, err = r.callOpenAIDirect(ctx, prompt)
	case "openrouter":
		aiResponse, err = r.callOpenRouterAPI(ctx, prompt)
	case "deepseek":
		aiResponse, err = r.callDeepSeekAPI(ctx, prompt)
	default:
		err = fmt.Errorf("unknown API type: %s", r.apiType)
	}

	if err != nil {
		fmt.Printf("❌ %s API call failed: %v\n", r.apiType, err)
		fmt.Println("🔄 Falling back to smart mock...")
		return r.smartMockReview(req)
	}

	fmt.Printf("✅ %s API responded successfully\n", r.apiType)
	return r.parseAIResponse(aiResponse, req)
}

func (r *AIReviewer) callOpenAIDirect(ctx context.Context, prompt string) (string, error) {
	client, ok := r.client.(*openai.Client)
	if !ok {
		return "", fmt.Errorf("invalid OpenAI client")
	}

	fmt.Println("📡 Calling OpenAI API...")

	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: r.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `Ты - опытный код-ревьюер. Анализируй код студентов и давай конструктивную обратную связь. 
Отвечай ТОЛЬКО в валидном JSON формате без каких-либо дополнительных текстов.`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   1500,
			Temperature: 0.3,
		},
	)

	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %v", err)
	}

	fmt.Printf("📥 Received response from OpenAI, tokens: %d\n", resp.Usage.TotalTokens)
	return resp.Choices[0].Message.Content, nil
}

func (r *AIReviewer) callOpenRouterAPI(ctx context.Context, prompt string) (string, error) {
	client, ok := r.client.(*http.Client)
	if !ok {
		return "", fmt.Errorf("invalid HTTP client")
	}

	fmt.Println("📡 Calling OpenRouter API...")

	requestBody := map[string]interface{}{
		"model": r.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Ты - опытный код-ревьюер. Анализируй код студентов и давай конструктивную обратную связь. Отвечай ТОЛЬКО в валидном JSON формате без каких-либо дополнительных текстов.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  1500,
		"temperature": 0.3,
	}

	return r.makeHTTPRequest(ctx, client, r.apiURL, r.apiKey, requestBody, "openrouter")
}

func (r *AIReviewer) callDeepSeekAPI(ctx context.Context, prompt string) (string, error) {
	client, ok := r.client.(*http.Client)
	if !ok {
		return "", fmt.Errorf("invalid HTTP client")
	}

	fmt.Println("📡 Calling DeepSeek API...")

	requestBody := map[string]interface{}{
		"model": r.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "Ты - опытный код-ревьюер. Анализируй код студентов и давай конструктивную обратную связь. Отвечай ТОЛЬКО в валидном JSON формате без каких-либо дополнительных текстов.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  1500,
		"temperature": 0.3,
	}

	return r.makeHTTPRequest(ctx, client, r.apiURL, r.apiKey, requestBody, "deepseek")
}

func (r *AIReviewer) makeHTTPRequest(ctx context.Context, client *http.Client, url, apiKey string, body map[string]interface{}, apiType string) (string, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	fmt.Printf("📤 Sending request to %s\n", url)
	fmt.Printf("   Body size: %d bytes\n", len(jsonBody))

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Дополнительные заголовки для OpenRouter
	if apiType == "openrouter" {
		req.Header.Set("HTTP-Referer", "http://localhost:8080")
		req.Header.Set("X-Title", "Code Review Platform")
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("📥 Response status: %d %s\n", resp.StatusCode, resp.Status)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != 200 {
		// Ограничиваем вывод ошибки
		errorText := string(responseBody)
		if len(errorText) > 200 {
			errorText = errorText[:200] + "..."
		}
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, errorText)
	}

	fmt.Printf("✅ API response size: %d bytes\n", len(responseBody))

	var result map[string]interface{}
	if err := json.Unmarshal(responseBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse API response: %v", err)
	}

	// Извлекаем текст ответа
	if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					fmt.Printf("📄 Extracted AI response, length: %d chars\n", len(content))
					return content, nil
				}
			}
		}
	}

	return "", fmt.Errorf("invalid API response structure")
}

// ==================== SMART MOCK METHODS ====================

func (r *AIReviewer) smartMockReview(req CodeReviewRequest) (*CodeReviewResponse, error) {
	fmt.Println("🤖 Generating smart mock review...")

	score := r.analyzeCodeQuality(req.Code, req.Language)

	// Для мока добавляем больше реализма
	comments := r.generateSmartComments(req.Code, req.Language)
	suggestions := r.generateSmartSuggestions(req.Code, req.Language)
	bestPractices := r.getBestPractices(req.Language)
	complexity := r.analyzeComplexity(req.Code)

	// Создаем реалистичный текстовый ответ
	aiResponse := fmt.Sprintf(
		"Код на %s успешно проанализирован. Оценка: %d/10. Сложность: %s. %s",
		req.Language,
		score,
		complexity,
		"Для получения детального анализа с использованием ИИ добавьте API ключ OpenRouter в настройки.",
	)

	fmt.Printf("📊 Mock review generated. Score: %d/10, Complexity: %s\n", score, complexity)

	return &CodeReviewResponse{
		Score:         score,
		Comments:      comments,
		Suggestions:   suggestions,
		BestPractices: bestPractices,
		Complexity:    complexity,
		AIResponse:    aiResponse,
	}, nil
}

func (r *AIReviewer) analyzeCodeQuality(code, language string) int {
	score := 6 // базовый балл
	lines := strings.Count(code, "\n") + 1

	// Анализ качества кода
	if lines > 10 {
		score += 1 // Больше кода = обычно лучше структурирован
	}
	if strings.Contains(code, "def ") || strings.Contains(code, "function ") {
		score += 2 // Использование функций - хорошо
	}
	if strings.Contains(code, "//") || strings.Contains(code, "#") || strings.Contains(code, "/*") {
		score += 1 // Комментарии
	}
	if strings.Contains(code, "\t") || strings.Contains(code, "    ") {
		score += 1 // Правильное форматирование
	}
	if !r.hasSyntaxErrors(code, language) {
		score += 1 // Нет синтаксических ошибок
	}

	// Проверка на наличие обработки ошибок
	if strings.Contains(code, "try") || strings.Contains(code, "except") ||
		strings.Contains(code, "catch") || strings.Contains(code, "if err") {
		score += 1
	}

	// Ограничиваем максимальный балл
	if score > 10 {
		score = 10
	}
	if score < 1 {
		score = 1
	}

	return score
}

func (r *AIReviewer) hasSyntaxErrors(code, language string) bool {
	// Базовая проверка синтаксиса
	switch language {
	case "python":
		// Проверяем баланс скобок
		if strings.Count(code, "(") != strings.Count(code, ")") {
			return true
		}
		if strings.Count(code, "[") != strings.Count(code, "]") {
			return true
		}
		if strings.Count(code, "{") != strings.Count(code, "}") {
			return true
		}
	case "javascript", "java", "cpp":
		// Для C-подобных языков проверяем базовый синтаксис
		lines := strings.Split(code, "\n")
		for i, line := range lines {
			trimmed := strings.TrimSpace(line)
			// Пропускаем пустые строки и комментарии
			if trimmed == "" || strings.HasPrefix(trimmed, "//") ||
				strings.HasPrefix(trimmed, "/*") || strings.HasPrefix(trimmed, "*") {
				continue
			}

			// Проверяем непарные скобки
			if strings.Count(line, "(") != strings.Count(line, ")") &&
				!strings.Contains(line, "if") && !strings.Contains(line, "while") &&
				!strings.Contains(line, "for") {
				fmt.Printf("⚠️ Possible syntax error line %d: %s\n", i+1, line)
				return true
			}
		}
	}

	return false
}

func (r *AIReviewer) generateSmartComments(code, language string) []string {
	comments := []string{"Код решает поставленную задачу"}
	lines := strings.Count(code, "\n") + 1

	if lines < 5 {
		comments = append(comments, "Код слишком короткий, можно добавить больше функциональности")
	} else if lines > 50 {
		comments = append(comments, "Код довольно длинный, можно разделить на функции/модули")
	}

	if !strings.Contains(code, "//") && !strings.Contains(code, "#") && !strings.Contains(code, "/*") {
		comments = append(comments, "Добавьте комментарии для лучшей читаемости")
	}

	if strings.Count(code, "\"") > 10 || strings.Count(code, "'") > 10 {
		comments = append(comments, "Много хардкоженных значений, рекомендуется вынести в константы")
	}

	// Языко-специфичные комментарии
	switch language {
	case "python":
		if strings.Contains(code, "print(") && !strings.Contains(code, "def ") {
			comments = append(comments, "Рекомендуется вынести логику в функции")
		}
		if strings.Contains(code, "eval(") || strings.Contains(code, "exec(") {
			comments = append(comments, "Использование eval/exec может быть небезопасным")
		}
	case "javascript":
		if strings.Contains(code, "var ") {
			comments = append(comments, "Используйте const/let вместо var")
		}
		if strings.Contains(code, "==") && !strings.Contains(code, "===") {
			comments = append(comments, "Для сравнения используйте === вместо ==")
		}
	case "java":
		if strings.Contains(code, "System.out.println") && !strings.Contains(code, "class ") {
			comments = append(comments, "Рекомендуется создать класс для структурирования кода")
		}
	}

	return comments
}

func (r *AIReviewer) generateSmartSuggestions(code, language string) []string {
	suggestions := []string{
		"Добавить обработку ошибок и крайних случаев",
		"Протестировать код на различных входных данных",
	}

	switch language {
	case "python":
		suggestions = append(suggestions,
			"Добавить type hints для функций",
			"Использовать docstring для документирования",
			"Рассмотреть использование list/dict comprehensions",
		)
	case "javascript":
		suggestions = append(suggestions,
			"Использовать стрелочные функции где возможно",
			"Добавить проверки типов с TypeScript",
			"Использовать async/await вместо callbacks",
		)
	case "java":
		suggestions = append(suggestions,
			"Следовать Java Code Conventions",
			"Использовать модификаторы доступа",
			"Рассмотреть использование Stream API",
		)
	case "cpp":
		suggestions = append(suggestions,
			"Использовать smart pointers вместо raw pointers",
			"Добавить обработку исключений",
			"Использовать const-correctness",
		)
	}

	if !strings.Contains(code, "def ") && !strings.Contains(code, "function ") {
		suggestions = append(suggestions, "Вынести код в отдельные функции")
	}

	if strings.Count(code, "\n") > 20 {
		suggestions = append(suggestions, "Разбить код на несколько модулей/файлов")
	}

	return suggestions
}

func (r *AIReviewer) analyzeComplexity(code string) string {
	lines := strings.Count(code, "\n") + 1
	complexity := 0

	// Анализ сложности
	if strings.Contains(code, "for ") || strings.Contains(code, "while ") {
		complexity++
	}
	if strings.Contains(code, "if ") || strings.Contains(code, "switch ") {
		complexity++
	}
	if strings.Contains(code, "def ") || strings.Contains(code, "function ") {
		complexity++
	}
	if strings.Contains(code, "class ") || strings.Contains(code, "struct ") {
		complexity++
	}

	// Определяем уровень сложности
	if lines <= 10 && complexity <= 1 {
		return "низкая"
	} else if lines <= 30 && complexity <= 3 {
		return "средняя"
	} else {
		return "высокая"
	}
}

func (r *AIReviewer) getBestPractices(language string) []string {
	switch language {
	case "python":
		return []string{
			"Следование PEP8",
			"Использование type hints",
			"Правильное именование переменных (snake_case)",
			"Использование виртуального окружения",
			"Документирование с помощью docstrings",
		}
	case "javascript":
		return []string{
			"Использование const/let вместо var",
			"Стрелочные функции для сохранения контекста",
			"Template literals для строк",
			"Деструктуризация объектов и массивов",
			"Использование модулей (import/export)",
		}
	case "java":
		return []string{
			"Следование Java Code Conventions",
			"Использование модификаторов доступа",
			"Обработка исключений",
			"Принципы ООП (инкапсуляция, наследование, полиморфизм)",
			"Использование интерфейсов",
		}
	case "cpp":
		return []string{
			"Использование smart pointers",
			"Следование RAII",
			"Избегание raw pointers",
			"Использование STL контейнеров",
			"const-correctness",
		}
	default:
		return []string{
			"Читаемость кода",
			"Разделение ответственности",
			"Обработка ошибок",
			"Тестирование кода",
			"Документирование",
		}
	}
}

func (r *AIReviewer) parseAIResponse(aiResponse string, req CodeReviewRequest) (*CodeReviewResponse, error) {
	fmt.Println("🔍 Parsing AI response...")

	var response CodeReviewResponse
	cleanedResponse := r.cleanJSONResponse(aiResponse)

	if err := json.Unmarshal([]byte(cleanedResponse), &response); err != nil {
		fmt.Printf("❌ Failed to parse AI response as JSON: %v\n", err)
		fmt.Printf("Raw response (first 500 chars):\n%s\n",
			func() string {
				if len(aiResponse) > 500 {
					return aiResponse[:500] + "..."
				}
				return aiResponse
			}())

		// Создаем fallback ответ
		return r.createFallbackResponse(aiResponse, req)
	}

	response.AIResponse = aiResponse
	fmt.Printf("✅ Successfully parsed AI response. Score: %d/10, Complexity: %s\n",
		response.Score, response.Complexity)

	return &response, nil
}

func (r *AIReviewer) cleanJSONResponse(response string) string {
	// Убираем markdown блоки кода
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	// Иногда AI добавляет пояснения перед JSON
	startIdx := strings.Index(response, "{")
	if startIdx > 0 {
		response = response[startIdx:]
	}

	// Ищем конец JSON
	endIdx := strings.LastIndex(response, "}")
	if endIdx != -1 {
		response = response[:endIdx+1]
	}

	return response
}

func (r *AIReviewer) createFallbackResponse(aiResponse string, req CodeReviewRequest) (*CodeReviewResponse, error) {
	fmt.Println("🔄 Creating fallback response from AI text...")

	score := r.extractScore(aiResponse)
	comments := r.extractComments(aiResponse)

	if len(comments) == 0 {
		comments = []string{
			"Код решает поставленную задачу",
			"Есть возможности для улучшения",
		}
	}

	return &CodeReviewResponse{
		Score:    score,
		Comments: comments,
		Suggestions: []string{
			"Рассмотреть альтернативные подходы",
			"Добавить обработку крайних случаев",
			"Протестировать с разными входными данными",
		},
		BestPractices: r.getBestPractices(req.Language),
		Complexity:    r.analyzeComplexity(req.Code),
		AIResponse:    aiResponse,
	}, nil
}

func (r *AIReviewer) extractScore(text string) int {
	// Ищем оценку в тексте
	text = strings.ToLower(text)

	if strings.Contains(text, "10/10") || strings.Contains(text, "отлично") {
		return 10
	} else if strings.Contains(text, "9/10") || strings.Contains(text, "превосходно") {
		return 9
	} else if strings.Contains(text, "8/10") || strings.Contains(text, "очень хорошо") {
		return 8
	} else if strings.Contains(text, "7/10") || strings.Contains(text, "хорошо") {
		return 7
	} else if strings.Contains(text, "6/10") || strings.Contains(text, "удовлетворительно") {
		return 6
	} else {
		return 6 // Минимальная оценка
	}
}

func (r *AIReviewer) extractComments(text string) []string {
	var comments []string
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) < 10 {
			continue
		}

		// Ищем пункты списка
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "•") ||
			strings.HasPrefix(line, "*") || (len(line) > 2 && line[0] >= '0' && line[0] <= '9' && line[1] == '.') {

			// Очищаем маркеры
			comment := strings.TrimPrefix(line, "-")
			comment = strings.TrimPrefix(comment, "•")
			comment = strings.TrimPrefix(comment, "*")
			comment = strings.TrimSpace(comment)

			// Убираем нумерацию "1. ", "2. " и т.д.
			if len(comment) > 2 && comment[0] >= '0' && comment[0] <= '9' && comment[1] == '.' {
				comment = comment[2:]
				comment = strings.TrimSpace(comment)
			}

			if comment != "" && !strings.Contains(comment, "{") && !strings.Contains(comment, "}") {
				comments = append(comments, comment)
			}
		}
	}

	return comments
}

func (r *AIReviewer) buildPrompt(code, language, taskContext string) string {
	return `Проанализируй этот код на языке программирования ` + language + ` и дай развернутую оценку.

КОНТЕКСТ ЗАДАЧИ: ` + taskContext + `

КОД СТУДЕНТА:
` + "```" + language + "\n" + code + "\n" + "```" + `

Пожалуйста, проанализируй код и верни ответ СТРОГО в следующем JSON формате:
{
  "score": 7,
  "comments": ["комментарий 1", "комментарий 2"],
  "suggestions": ["предложение 1", "предложение 2"],
  "best_practices": ["практика 1", "практика 2"],
  "complexity": "низкая",
  "ai_response": "Полный текстовый анализ на русском языке"
}

Критерии для анализа:
1. Корректность решения задачи
2. Читаемость и стиль кода
3. Эффективность алгоритма
4. Следование best practices для ` + language + `
5. Безопасность кода
6. Оптимальность решения

Будь конкретен и давай практические советы для улучшения. Не добавляй никакого текста кроме JSON.`
}
