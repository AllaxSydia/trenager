package models

import "time"

// ExecutionResult представляет результат выполнения кода
type ExecutionResult struct {
	ID            string        `json:"id"`
	TaskID        string        `json:"task_id"`
	UserID        string        `json:"user_id"`
	Code          string        `json:"code"`
	Language      string        `json:"language"`
	Output        string        `json:"output"`
	Success       bool          `json:"success"`
	Error         string        `json:"error,omitempty"`
	ExecutionTime time.Duration `json:"execution_time"`
	CreatedAt     time.Time     `json:"created_at"`
}

// DockerExecutionConfig конфигурация для Docker контейнера
type DockerExecutionConfig struct {
	Image     string            `json:"image"`
	Cmd       []string          `json:"cmd"`
	Memory    int64             `json:"memory"`     // в байтах
	CPUShares int64             `json:"cpu_shares"` // в %
	Timeout   time.Duration     `json:"timeout"`
	Env       map[string]string `json:"env"`
}

// LanguageConfig конфигурация для разных языков программирования
type LanguageConfig struct {
	DockerImage string        `json:"docker_image"`
	CompileCmd  []string      `json:"compile_cmd,omitempty"`
	RunCmd      []string      `json:"run_cmd"`
	FileName    string        `json:"file_name"`
	Timeout     time.Duration `json:"timeout"`
}
