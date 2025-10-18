package config

import (
	"os"
	"strconv"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	Docker struct {
		Host string
	}
}

func Load() *Config {
	var cfg Config

	// Server config
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")

	// Database config
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "password")
	cfg.Database.DBName = getEnv("DB_NAME", "trenager")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// Docker config
	cfg.Docker.Host = getEnv("DOCKER_HOST", "unix:///var/run/docker.sock")

	return &cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
