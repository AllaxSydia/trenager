package config

import (
	"os"
	"strings"
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

// Добавь это если нет:
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Создаёт и заполняет конфиг при запуске приложения
func Load() *Config {
	var cfg Config

	// Railway использует DATABASE_URL или отдельные переменные
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		// Парсим DATABASE_URL
		cfg.Database = parseDatabaseURL(dbURL)
	} else {
		// Используем отдельные переменные
		cfg.Database.Host = getEnv("DB_HOST", "localhost")
		cfg.Database.Port = getEnv("DB_PORT", "5432")
		cfg.Database.User = getEnv("DB_USER", "postgres")
		cfg.Database.Password = getEnv("DB_PASSWORD", "password")
		cfg.Database.DBName = getEnv("DB_NAME", "codeplatform")
	}

	cfg.Server.Port = getEnv("PORT", "8080")

	// SSLMode: require для продакшна, disable для разработки
	if os.Getenv("RAILWAY_ENVIRONMENT") != "" {
		cfg.Database.SSLMode = "require"
	} else {
		cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")
	}

	return &cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// parseDatabaseURL парсит DATABASE_URL формата:
// postgresql://username:password@host:port/database
func parseDatabaseURL(dbURL string) DatabaseConfig {
	var cfg DatabaseConfig

	// Упрощенный парсинг DATABASE_URL
	if strings.HasPrefix(dbURL, "postgresql://") {
		// Убираем префикс
		urlWithoutPrefix := strings.TrimPrefix(dbURL, "postgresql://")

		// Разделяем на части: user:pass@host:port/dbname
		parts := strings.Split(urlWithoutPrefix, "@")
		if len(parts) == 2 {
			// Парсим user:password
			userPass := strings.Split(parts[0], ":")
			if len(userPass) == 2 {
				cfg.User = userPass[0]
				cfg.Password = userPass[1]
			}

			// Парсим host:port/dbname
			hostPortDb := strings.Split(parts[1], "/")
			if len(hostPortDb) == 2 {
				cfg.DBName = hostPortDb[1]

				// Парсим host:port
				hostPort := strings.Split(hostPortDb[0], ":")
				if len(hostPort) == 2 {
					cfg.Host = hostPort[0]
					cfg.Port = hostPort[1]
				} else {
					cfg.Host = hostPort[0]
					cfg.Port = "5432"
				}
			}
		}
	}

	return cfg
}
