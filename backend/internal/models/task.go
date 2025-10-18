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
