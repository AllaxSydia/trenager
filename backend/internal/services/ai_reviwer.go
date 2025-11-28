package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// Добавь эти структуры в начало файла
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
	client       *openai.Client
	useSmartMock bool
}

func NewAIReviewer() *AIReviewer {
	// Проверяем доступность API ключей
	openaiKey := os.Getenv("OPENAI_API_KEY")
	geminiKey := os.Getenv("GEMINI_API_KEY")

	if openaiKey != "" {
		fmt.Println("🔑 AI Reviewer: Using OpenAI API")
		client := openai.NewClient(openaiKey)
		return &AIReviewer{client: client, useSmartMock: false}
	} else if geminiKey != "" {
		fmt.Println("🔑 AI Reviewer: Using Gemini API")
		// Здесь можно добавить клиент для Gemini
		return &AIReviewer{client: nil, useSmartMock: true}
	} else {
		fmt.Println("🔧 AI Reviewer: Using smart mock mode")
		return &AIReviewer{client: nil, useSmartMock: true}
	}
}

func (r *AIReviewer) ReviewCode(req CodeReviewRequest) (*CodeReviewResponse, error) {
	if r.useSmartMock || r.client == nil {
		fmt.Println("🔧 Using smart mock reviewer")
		return r.smartMockReview(req)
	}

	// Используем OpenAI если доступен
	fmt.Printf("🔑 AI Reviewer: Using OpenAI API for %s code\n", req.Language)

	prompt := r.buildPrompt(req.Code, req.Language, req.TaskContext)

	resp, err := r.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
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
			MaxTokens:   2000,
			Temperature: 0.3,
		},
	)

	if err != nil {
		fmt.Printf("❌ OpenAI API error: %v, falling back to smart mock\n", err)
		return r.smartMockReview(req)
	}

	aiResponse := resp.Choices[0].Message.Content
	return r.parseAIResponse(aiResponse, req)
}

// Добавляем метод smartMockReview в AIReviewer
func (r *AIReviewer) smartMockReview(req CodeReviewRequest) (*CodeReviewResponse, error) {
	score := r.analyzeCodeQuality(req.Code, req.Language)

	return &CodeReviewResponse{
		Score:         score,
		Comments:      r.generateSmartComments(req.Code, req.Language),
		Suggestions:   r.generateSmartSuggestions(req.Code, req.Language),
		BestPractices: r.getBestPractices(req.Language),
		Complexity:    r.analyzeComplexity(req.Code),
		AIResponse:    "Умный анализ кода. Для подключения ИИ требуется API ключ или сервер с 8+ GB RAM.",
	}, nil
}

func (r *AIReviewer) analyzeCodeQuality(code, language string) int {
	score := 6 // базовый балл

	// Анализ читаемости
	lines := strings.Count(code, "\n") + 1
	if lines > 3 {
		score += 1
	}

	// Проверка на наличие функций
	if strings.Contains(code, "def ") || strings.Contains(code, "function ") {
		score += 1
	}

	// Проверка комментариев
	if strings.Contains(code, "//") || strings.Contains(code, "#") || strings.Contains(code, "/*") {
		score += 1
	}

	// Проверка форматирования
	if strings.Contains(code, "\t") || strings.Contains(code, "    ") {
		score += 1
	}

	// Проверка на ошибки (базовые)
	if !r.hasSyntaxErrors(code, language) {
		score += 1
	}

	if score > 10 {
		score = 10
	}

	return score
}

func (r *AIReviewer) hasSyntaxErrors(code, language string) bool {
	// Простые проверки синтаксиса
	switch language {
	case "python":
		// Проверяем несбалансированные скобки
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
		// Проверяем точку с запятой в конце (базово)
		lines := strings.Split(code, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" && !strings.HasPrefix(trimmed, "//") &&
				!strings.HasPrefix(trimmed, "/*") && !strings.HasPrefix(trimmed, "*") &&
				!strings.Contains(trimmed, "{") && !strings.Contains(trimmed, "}") &&
				!strings.HasSuffix(trimmed, "{") && !strings.HasPrefix(trimmed, "}") {
				if !strings.HasSuffix(trimmed, ";") && !strings.HasSuffix(trimmed, "{") {
					return true
				}
			}
		}
	}
	return false
}

func (r *AIReviewer) generateSmartComments(code, language string) []string {
	comments := []string{"Код решает поставленную задачу"}

	// Анализ структуры
	lines := strings.Count(code, "\n") + 1
	if lines < 5 {
		comments = append(comments, "Код слишком короткий, можно добавить больше функциональности")
	}

	// Проверка комментариев
	if !strings.Contains(code, "//") && !strings.Contains(code, "#") && !strings.Contains(code, "/*") {
		comments = append(comments, "Добавьте комментарии для лучшей читаемости")
	}

	// Проверка форматирования
	if !strings.Contains(code, "\n") {
		comments = append(comments, "Рекомендуется разбить код на несколько строк")
	}

	// Проверка на хардкод
	if strings.Count(code, "\"") > 10 || strings.Count(code, "'") > 10 {
		comments = append(comments, "Много хардкоженных значений,可以考虑 вынести в константы")
	}

	// Языко-специфичные проверки
	switch language {
	case "python":
		if strings.Contains(code, "print(") && !strings.Contains(code, "def ") {
			comments = append(comments, "Рекомендуется вынести логику в функции")
		}
	case "javascript":
		if strings.Contains(code, "var ") {
			comments = append(comments, "Используйте const/let вместо var")
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

	// Языко-специфичные предложения
	switch language {
	case "python":
		suggestions = append(suggestions,
			"Добавить type hints для функций",
			"Использовать docstring для документирования",
		)
	case "javascript":
		suggestions = append(suggestions,
			"Использовать стрелочные функции где возможно",
			"Добавить проверки типов с TypeScript",
		)
	case "java":
		suggestions = append(suggestions,
			"Следовать Java Code Conventions",
			"Использовать модификаторы доступа",
		)
	case "cpp":
		suggestions = append(suggestions,
			"Использовать smart pointers вместо raw pointers",
			"Добавить обработку исключений",
		)
	}

	// Предложения по структуре
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

	// Простой анализ сложности
	if strings.Contains(code, "for ") || strings.Contains(code, "while ") {
		complexity++
	}
	if strings.Contains(code, "if ") || strings.Contains(code, "switch ") {
		complexity++
	}
	if strings.Contains(code, "def ") || strings.Contains(code, "function ") {
		complexity++
	}

	if lines <= 5 && complexity == 0 {
		return "низкая"
	} else if lines <= 15 && complexity <= 2 {
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
		}
	case "javascript":
		return []string{
			"Использование const/let вместо var",
			"Стрелочные функции для сохранения контекста",
			"Template literals для строк",
			"Деструктуризация объектов и массивов",
		}
	case "java":
		return []string{
			"Следование Java Code Conventions",
			"Использование модификаторов доступа",
			"Обработка исключений",
			"Принципы ООП (инкапсуляция, наследование, полиморфизм)",
		}
	case "cpp":
		return []string{
			"Использование smart pointers",
			"Следование RAII",
			"Избегание raw pointers",
			"Использование STL контейнеров",
		}
	default:
		return []string{
			"Читаемость кода",
			"Разделение ответственности",
			"Обработка ошибок",
			"Тестирование кода",
		}
	}
}

func (r *AIReviewer) buildPrompt(code, language, taskContext string) string {
	return fmt.Sprintf(`
Проанализируй этот код на %s и дай развернутую оценку в JSON формате.

КОНТЕКСТ ЗАДАЧИ: %s

КОД ДЛЯ АНАЛИЗА:
%s

Проанализируй по следующим критериям:
1. Корректность решения задачи
2. Читаемость и стиль кода
3. Эффективность алгоритма
4. Следование best practices для языка %s
5. Безопасность (если применимо)
6. Оптимальность решения

Верни ответ в строгом JSON формате:
{
	"score": 8,
	"comments": ["конкретное замечание 1", "конкретное замечание 2"],
	"suggestions": ["конкретное предложение 1", "конкретное предложение 2"],
	"best_practices": ["практика 1", "практика 2"],
	"complexity": "низкая/средняя/высокая",
	"ai_response": "полный текстовый анализ"
}

Будь конкретен и давай практические советы для улучшения. Не добавляй никакого текста кроме JSON.`, language, taskContext, code, language)
}

func (r *AIReviewer) parseAIResponse(aiResponse string, req CodeReviewRequest) (*CodeReviewResponse, error) {
	var response CodeReviewResponse

	// Очищаем ответ от возможных markdown блоков кода
	cleanedResponse := r.cleanJSONResponse(aiResponse)

	// Пытаемся распарсить JSON ответ
	if err := json.Unmarshal([]byte(cleanedResponse), &response); err != nil {
		fmt.Printf("❌ Failed to parse AI response as JSON: %v\n", err)
		fmt.Printf("Raw response: %s\n", aiResponse)
		return r.createFallbackResponse(aiResponse, req)
	}

	// Добавляем полный ответ ИИ
	response.AIResponse = aiResponse
	fmt.Printf("✅ Successfully parsed AI response, score: %d\n", response.Score)

	return &response, nil
}

func (r *AIReviewer) cleanJSONResponse(response string) string {
	// Убираем markdown блоки кода если есть
	response = strings.TrimPrefix(response, "```json")
	response = strings.TrimPrefix(response, "```")
	response = strings.TrimSuffix(response, "```")
	response = strings.TrimSpace(response)

	// Иногда GPT добавляет пояснения перед JSON
	if idx := strings.Index(response, "{"); idx > 0 {
		response = response[idx:]
	}

	return response
}

func (r *AIReviewer) createFallbackResponse(aiResponse string, req CodeReviewRequest) (*CodeReviewResponse, error) {
	// Если ИИ не вернул JSON, анализируем текстовый ответ
	score := r.extractScore(aiResponse)
	comments := r.extractComments(aiResponse)

	return &CodeReviewResponse{
		Score:    score,
		Comments: comments,
		Suggestions: []string{
			"Рассмотреть альтернативные подходы",
			"Добавить обработку крайних случаев",
		},
		BestPractices: []string{
			"Использовать осмысленные имена переменных",
			"Разделять код на небольшие функции",
		},
		Complexity: "средняя",
		AIResponse: aiResponse,
	}, nil
}

func (r *AIReviewer) extractScore(text string) int {
	// Простая эвристика для извлечения оценки из текста
	if strings.Contains(text, "10/10") || strings.Contains(text, "отлично") {
		return 10
	} else if strings.Contains(text, "9/10") || strings.Contains(text, "превосходно") {
		return 9
	} else if strings.Contains(text, "8/10") || strings.Contains(text, "очень хорошо") {
		return 8
	} else if strings.Contains(text, "7/10") || strings.Contains(text, "хорошо") {
		return 7
	} else {
		return 6
	}
}

func (r *AIReviewer) extractComments(text string) []string {
	// Простая эвристика для извлечения комментариев
	var comments []string
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if (strings.HasPrefix(line, "-") || strings.HasPrefix(line, "•")) && len(line) > 10 {
			comment := strings.TrimPrefix(strings.TrimPrefix(line, "-"), "•")
			comment = strings.TrimSpace(comment)
			if comment != "" {
				comments = append(comments, comment)
			}
		}
	}

	if len(comments) == 0 {
		comments = []string{"Код решает задачу", "Есть возможности для улучшения"}
	}

	return comments
}
