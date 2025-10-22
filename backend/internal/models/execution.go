package models

import "time"

// ExecutionResult представляет результат выполнения кода
type ExecutionResult struct {
	ID            string        `json:"id"`              // Уникальный номер попытки
	TaskID        string        `json:"task_id"`         // ID задания
	UserID        string        `json:"user_id"`         // ID пользователя
	Code          string        `json:"code"`            // Сам код который выполнили
	Language      string        `json:"language"`        // Язык программирования
	Output        string        `json:"output"`          // Что программа напечатала
	Success       bool          `json:"success"`         // Успешно или с ошибкой
	Error         string        `json:"error,omitempty"` // Текст ошибки (если была)
	ExecutionTime time.Duration `json:"execution_time"`
	CreatedAt     time.Time     `json:"created_at"`
}

// DockerExecutionConfig конфигурация для Docker контейнера
type DockerExecutionConfig struct {
	Image     string            `json:"image"` // Какой образ использовать
	Cmd       []string          `json:"cmd"`
	Memory    int64             `json:"memory"`     // Сколько памяти дать
	CPUShares int64             `json:"cpu_shares"` // Сколько процессора
	Timeout   time.Duration     `json:"timeout"`    // Максимальное время работы
	Env       map[string]string `json:"env"`
}

// LanguageConfig конфигурация для разных языков программирования
type LanguageConfig struct {
	DockerImage string        `json:"docker_image"`          // Какой образ использовать
	CompileCmd  []string      `json:"compile_cmd,omitempty"` // Сколько памяти дать
	RunCmd      []string      `json:"run_cmd"`               // Сколько процессора
	FileName    string        `json:"file_name"`             // Максимальное время работы
	Timeout     time.Duration `json:"timeout"`
}
