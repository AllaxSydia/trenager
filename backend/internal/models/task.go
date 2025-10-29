package models

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Template    string `json:"template"`
	Tests       []Test `json:"tests"`
}

type Test struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

type ExecutionRequest struct {
	TaskID   string `json:"task_id"`
	Code     string `json:"code"`
	Language string `json:"language"`
}

type ExecutionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Output  string `json:"output"`
}

// CheckRequest - запрос на проверку решения
type CheckRequest struct {
	TaskID   interface{} `json:"task_id"` // ← принимает и строки и числа
	Code     string      `json:"code"`
	Language string      `json:"language"`
	Tests    []Test      `json:"tests,omitempty"`
}

// CheckResponse - ответ проверки решения
type CheckResponse struct {
	Success  bool   `json:"success"`
	Passed   bool   `json:"passed"`
	Output   string `json:"output"`
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
	Message  string `json:"message"`
}

// CheckResult - результат проверки
type CheckResult struct {
	Passed   bool   `json:"passed"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Message  string `json:"message"`
}
