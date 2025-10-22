package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// База данных
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConnection(cfg Config) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ PostgreSQL connected successfully")
	return db, nil
}

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		// Таблица пользователей
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(36) PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Таблица задач
		`CREATE TABLE IF NOT EXISTS tasks (
			id VARCHAR(36) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			template TEXT,
			difficulty INTEGER DEFAULT 1,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Таблица тестов для задач
		`CREATE TABLE IF NOT EXISTS task_tests (
			id VARCHAR(36) PRIMARY KEY,
			task_id VARCHAR(36) REFERENCES tasks(id) ON DELETE CASCADE,
			input TEXT,
			expected_output TEXT NOT NULL,
			is_hidden BOOLEAN DEFAULT false
		)`,

		// Таблица прогресса пользователей
		`CREATE TABLE IF NOT EXISTS user_progress (
			user_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
			task_id VARCHAR(36) REFERENCES tasks(id) ON DELETE CASCADE,
			completed BOOLEAN DEFAULT false,
			attempts INTEGER DEFAULT 0,
			best_score FLOAT DEFAULT 0,
			last_attempt TIMESTAMP,
			PRIMARY KEY (user_id, task_id)
		)`,

		// Таблица выполненных заданий
		`CREATE TABLE IF NOT EXISTS code_executions (
			id VARCHAR(36) PRIMARY KEY,
			user_id VARCHAR(36) REFERENCES users(id) ON DELETE CASCADE,
			task_id VARCHAR(36) REFERENCES tasks(id) ON DELETE CASCADE,
			code TEXT NOT NULL,
			language VARCHAR(20) NOT NULL,
			output TEXT,
			success BOOLEAN DEFAULT false,
			execution_time INTEGER,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	log.Println("✅ Database migrations completed successfully")
	return nil
}
