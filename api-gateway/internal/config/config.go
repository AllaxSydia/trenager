package config

import (
	"os"
)

type Config struct {
	Port     string
	Services map[string]ServiceConfig
}

type ServiceConfig struct {
	Name string
	Host string
	Port string
	URL  string
}

func LoadConfig() *Config {
	return &Config{
		Port: getEnv("GATEWAY_PORT", "8080"),
		Services: map[string]ServiceConfig{
			"auth": {
				Name: "auth-service",
				Host: getEnv("AUTH_SERVICE_HOST", "localhost"),
				Port: getEnv("AUTH_SERVICE_PORT", "50051"),
				URL:  getEnv("AUTH_SERVICE_URL", "localhost:50051"),
			},
			"task": {
				Name: "task-service",
				Host: getEnv("TASK_SERVICE_HOST", "localhost"),
				Port: getEnv("TASK_SERVICE_PORT", "50052"),
				URL:  getEnv("TASK_SERVICE_URL", "localhost:50052"),
			},
			"grading": {
				Name: "grading-service",
				Host: getEnv("GRADING_SERVICE_HOST", "localhost"),
				Port: getEnv("GRADING_SERVICE_PORT", "50053"),
				URL:  getEnv("GRADING_SERVICE_URL", "localhost:50053"),
			},
			"execution": {
				Name: "execution-service",
				Host: getEnv("EXECUTION_SERVICE_HOST", "localhost"),
				Port: getEnv("EXECUTION_SERVICE_PORT", "50054"),
				URL:  getEnv("EXECUTION_SERVICE_URL", "localhost:50054"),
			},
			"ai": {
				Name: "ai-service",
				Host: getEnv("AI_SERVICE_HOST", "localhost"),
				Port: getEnv("AI_SERVICE_PORT", "50055"),
				URL:  getEnv("AI_SERVICE_URL", "localhost:50055"),
			},
			"analytics": {
				Name: "analytics-service",
				Host: getEnv("ANALYTICS_SERVICE_HOST", "localhost"),
				Port: getEnv("ANALYTICS_SERVICE_PORT", "50056"),
				URL:  getEnv("ANALYTICS_SERVICE_URL", "localhost:50056"),
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
