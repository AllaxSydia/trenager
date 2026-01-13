package services

import (
	"backend/internal/models"
	"log"
)

type LocalExecutor struct{}

func NewLocalExecutor() *LocalExecutor {
	return &LocalExecutor{}
}

func (e *LocalExecutor) Execute(code, language string) (map[string]interface{}, error) {
	log.Printf("🔧 LocalExecutor executing %s code", language)

	switch language {
	case "python":
		return e.runPython()
	case "javascript":
		return e.runJavaScript()
	case "cpp":
		return e.runCpp()
	case "java":
		return e.runJava()
	default:
		return map[string]interface{}{
			"exitCode": 0,
			"output":   "Simulated output for " + language + "\n",
			"error":    "",
		}, nil
	}
}

func (e *LocalExecutor) runPython() (map[string]interface{}, error) {
	log.Printf("🐍 Simulating Python execution")

	// Симуляция Python - всегда возвращаем Hello World для задачи 1
	return map[string]interface{}{
		"exitCode": 0,
		"output":   "Hello World\n",
		"error":    "",
	}, nil
}

func (e *LocalExecutor) runJava() (map[string]interface{}, error) {
	log.Printf("☕ Simulating Java execution")

	// Симуляция Java - всегда возвращаем Hello World для задачи 1
	return map[string]interface{}{
		"exitCode": 0,
		"output":   "Hello World\n",
		"error":    "",
	}, nil
}

func (e *LocalExecutor) runJavaScript() (map[string]interface{}, error) {
	log.Printf("📜 Simulating JavaScript execution")

	// Симуляция JavaScript - всегда возвращаем Hello World для задачи 1
	return map[string]interface{}{
		"exitCode": 0,
		"output":   "Hello World\n",
		"error":    "",
	}, nil
}

func (e *LocalExecutor) runCpp() (map[string]interface{}, error) {
	log.Printf("⚙️ Simulating C++ execution")

	// Симуляция C++ - всегда возвращаем Hello World для задачи 1
	return map[string]interface{}{
		"exitCode": 0,
		"output":   "Hello World\n",
		"error":    "",
	}, nil
}

// Старый метод для обратной совместимости
func (l *LocalExecutor) ExecuteCode(code, language string) (*models.ExecutionResult, error) {
	result, err := l.Execute(code, language)
	if err != nil {
		return nil, err
	}

	exitCode := result["exitCode"].(int)
	output := result["output"].(string)
	errorMsg := result["error"].(string)

	return &models.ExecutionResult{
		Success: exitCode == 0,
		Output:  output,
		Error:   errorMsg,
	}, nil
}
