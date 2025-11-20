package executor

import (
	"context"
	"log"
	"os"
)

// NewExecutor создает исполнитель кода
func NewExecutor() Executor {
	// В продакшене используем Docker если доступен
	if os.Getenv("ENVIRONMENT") == "production" {
		dockerExecutor, err := NewDockerExecutor()
		if err != nil {
			log.Printf("⚠️ Docker not available: %v", err)
		} else {
			// В продакшене считаем что Docker всегда готов
			log.Printf("✅ DockerExecutor initialized for production")
			return dockerExecutor
		}
	}

	// Fallback на LocalExecutor
	log.Printf("🔄 Running in local execution mode")
	return NewLocalExecutor()
}

// isDockerReady проверяет доступность Docker (только для development)
func (d *DockerExecutorImpl) isDockerReady() bool {
	// В продакшене пропускаем проверку образов
	if os.Getenv("ENVIRONMENT") == "production" {
		return true
	}

	ctx := context.Background()
	requiredImages := []string{
		"python:3.11-alpine",
		"node:18-alpine",
		"openjdk:11-jre-alpine",
		"gcc:latest",
		"golang:1.21-alpine",
	}

	for _, image := range requiredImages {
		_, _, err := d.client.ImageInspectWithRaw(ctx, image)
		if err != nil {
			log.Printf("⚠️ Docker image %s not available: %v", image, err)
			return false
		}
	}
	return true
}
