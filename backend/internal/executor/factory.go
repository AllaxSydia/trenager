package executor

import (
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
