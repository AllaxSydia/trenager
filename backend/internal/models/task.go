package models

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Language    string `json:"language,omitempty"`
	Difficulty  string `json:"difficulty,omitempty"`
	Template    string `json:"template"`
	Tests       []Test `json:"tests"`
}

type Test struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
}

type ExecutionRequest struct {
	TaskID   string   `json:"task_id"`
	Code     string   `json:"code"`
	Language string   `json:"language"`
	Inputs   []string `json:"inputs,omitempty"`
}

// ExecutionResponse представляет ответ от выполнения кода
type ExecutionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Output  string `json:"output"`
}

// CheckRequest - запрос на проверку решения
type CheckRequest struct {
	TaskID   interface{} `json:"task_id"` // принимает и строки и числа
	Code     string      `json:"code"`
	Language string      `json:"language"`
	Tests    []Test      `json:"tests,omitempty"`
}

// CheckResponse - ответ проверки решения
type CheckResponse struct {
	Success     bool         `json:"success"`
	Message     string       `json:"message"`
	TestResults []TestResult `json:"test_results"`
	TotalTests  int          `json:"total_tests"`
	PassedTests int          `json:"passed_tests"`
}

// TestResult - результат выполнения одного теста
type TestResult struct {
	TestNumber int    `json:"test_number"`
	Passed     bool   `json:"passed"`
	Output     string `json:"output"`
	Error      string `json:"error,omitempty"`
	Expected   string `json:"expected"`
	Actual     string `json:"actual"`
}

// CheckResult - устаревшая структура (оставлена для обратной совместимости)
type CheckResult struct {
	Passed   bool   `json:"passed"`
	Expected string `json:"expected"`
	Actual   string `json:"actual"`
	Message  string `json:"message"`
}
