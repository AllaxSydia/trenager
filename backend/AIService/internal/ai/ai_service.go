package ai

import (
	"context"
	"log"
	"strings"
)

type AIService struct {
	// Здесь будет клиент к OpenAI/Anthropic/локальной модели
}

func NewAIService() *AIService {
	return &AIService{}
}

// GetHint - генерирует подсказку для задачи
func (s *AIService) GetHint(ctx context.Context, taskID, userCode, language string, hintLevel int) (string, string, error) {
	log.Printf("🤖 Generating hint for task %s, level %d", taskID, hintLevel)

	hints := map[int]string{
		1: "Попробуйте разбить задачу на более мелкие части.",
		2: "Обратите внимание на входные данные и ожидаемый результат.",
		3: "Вам может помочь использование цикла или рекурсии.",
		4: "Попробуйте использовать структуру данных, такую как словарь или множество.",
		5: "Вот пример решения: функция принимает параметры и возвращает результат.",
	}

	hint := hints[hintLevel]
	if hintLevel < 1 || hintLevel > 5 {
		hint = hints[3]
	}

	exampleCode := `# Пример решения
def solution(data):
    # Ваш код здесь
    return result`

	if language == "go" {
		exampleCode = `// Пример решения
func solution(data interface{}) interface{} {
    // Ваш код здесь
    return result
}`
	} else if language == "javascript" {
		exampleCode = `// Пример решения
function solution(data) {
    // Ваш код здесь
    return result;
}`
	}

	return hint, exampleCode, nil
}

// ReviewCode - проверяет качество кода
func (s *AIService) ReviewCode(ctx context.Context, code, language, taskDescription string) (int, []string, []string, string, error) {
	log.Printf("🔍 Reviewing code for language %s", language)

	qualityScore := 85
	issues := []string{}
	suggestions := []string{}

	// Проверка на общие проблемы
	if len(code) < 10 {
		issues = append(issues, "Код слишком короткий")
		suggestions = append(suggestions, "Добавьте больше логики в решение")
		qualityScore -= 20
	}

	if strings.Contains(code, "TODO") {
		issues = append(issues, "Обнаружен TODO комментарий")
		suggestions = append(suggestions, "Завершите реализацию функции")
		qualityScore -= 10
	}

	if !strings.Contains(code, "return") && !strings.Contains(code, "return ") {
		suggestions = append(suggestions, "Убедитесь, что функция возвращает результат")
		qualityScore -= 15
	}

	overallFeedback := "Хорошая работа! Код в основном правильный, но есть небольшие улучшения."

	if qualityScore >= 90 {
		overallFeedback = "Отличный код! Хорошая структура и читаемость."
	} else if qualityScore >= 70 {
		overallFeedback = "Хорошая попытка! Есть несколько моментов для улучшения."
	} else {
		overallFeedback = "Код требует доработки. Обратите внимание на рекомендации."
	}

	if qualityScore < 0 {
		qualityScore = 0
	}

	return qualityScore, issues, suggestions, overallFeedback, nil
}

// GetRecommendations - рекомендует следующие задачи
func (s *AIService) GetRecommendations(ctx context.Context, userID string, completedTasks []string, weakTopics []string) ([]map[string]interface{}, error) {
	log.Printf("📚 Generating recommendations for user %s", userID)

	recommendations := []map[string]interface{}{
		{
			"task_id":          "rec-001",
			"title":            "Алгоритмы сортировки",
			"reason":           "На основе вашего прогресса",
			"difficulty_score": 0.3,
		},
		{
			"task_id":          "rec-002",
			"title":            "Работа со строками",
			"reason":           "Рекомендуется для укрепления навыков",
			"difficulty_score": 0.5,
		},
		{
			"task_id":          "rec-003",
			"title":            "Динамическое программирование",
			"reason":           "Вы готовы к более сложным задачам",
			"difficulty_score": 0.8,
		},
	}

	return recommendations, nil
}

// AskQuestion - отвечает на вопросы по коду
func (s *AIService) AskQuestion(ctx context.Context, question, codeContext, language string) (string, []string, error) {
	log.Printf("💬 Answering question: %s", question[:min(50, len(question))])

	// Простые ответы на частые вопросы
	answer := ""
	codeExamples := []string{}

	if strings.Contains(strings.ToLower(question), "error") || strings.Contains(strings.ToLower(question), "ошибк") {
		answer = "Ошибки часто возникают из-за несоответствия типов данных. Проверьте, что функция возвращает ожидаемый тип."
		codeExamples = append(codeExamples, "result = int(input)  # Преобразование к числу")
	} else if strings.Contains(strings.ToLower(question), "for") || strings.Contains(strings.ToLower(question), "цикл") {
		answer = "Циклы используются для повторения операций. for i in range(10): выполняет блок кода 10 раз."
		codeExamples = append(codeExamples, "for i in range(len(arr)):\n    print(arr[i])")
	} else if strings.Contains(strings.ToLower(question), "function") || strings.Contains(strings.ToLower(question), "функци") {
		answer = "Функции помогают организовать код. Используйте def для определения функции."
		codeExamples = append(codeExamples, "def my_function(param):\n    return param * 2")
	} else {
		answer = "Попробуйте разбить задачу на шаги и проверить каждый шаг отдельно. Убедитесь, что вы правильно обрабатываете входные данные."
	}

	return answer, codeExamples, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
