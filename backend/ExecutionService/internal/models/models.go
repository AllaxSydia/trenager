package models

type ExecutionRequest struct {
	ID            string
	Code          string
	Language      string
	Input         string
	TimeLimitMs   int
	MemoryLimitMB int
	TestCases     []TestCase
}

type ExecutionResult struct {
	ID              string
	Success         bool
	Output          string
	Error           string
	ExecutionTimeMs int64
	MemoryUsedBytes int64
	Status          string // completed, timeout, memory_error, runtime_error
	TestResults     []TestResult
}

type TestCase struct {
	Input          string
	ExpectedOutput string
	IsHidden       bool
}

type TestResult struct {
	ID              string
	TestID          string
	Passed          bool
	ActualOutput    string
	ExpectedOutput  string
	Error           string
	ExecutionTimeMs int64
}

type ExecutionStatus struct {
	ID     string
	Status string // pending, running, completed, failed
	Result *ExecutionResult
}
