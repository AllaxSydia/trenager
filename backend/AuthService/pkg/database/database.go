package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config содержит параметры подключения к базе данных
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

// NewPostgresConnection создает новое подключение к PostgreSQL
func NewPostgresConnection(cfg Config) (*sqlx.DB, error) {
	// Формируем строку подключения
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	log.Printf("🔌 Connecting to PostgreSQL at %s:%d/%s", cfg.Host, cfg.Port, cfg.Name)

	// Подключаемся к базе данных
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Настраиваем пул соединений
	db.SetMaxOpenConns(25)                 // Максимальное количество открытых соединений
	db.SetMaxIdleConns(10)                 // Максимальное количество idle соединений
	db.SetConnMaxLifetime(5 * time.Minute) // Максимальное время жизни соединения
	db.SetConnMaxIdleTime(2 * time.Minute) // Максимальное время idle соединения

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	// Проверяем существование таблицы users
	if err := checkUsersTable(db); err != nil {
		log.Printf("⚠️ Users table not found or error: %v", err)
		log.Println("💡 Run migrations to create the table")
	}

	log.Println("✅ PostgreSQL connection established successfully")

	// Логируем статистику подключения
	go logConnectionStats(db)

	return db, nil
}

// checkUsersTable проверяет существование таблицы users
func checkUsersTable(db *sqlx.DB) error {
	var tableExists bool
	query := `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'users'
		)
	`

	err := db.Get(&tableExists, query)
	if err != nil {
		return fmt.Errorf("failed to check users table: %w", err)
	}

	if !tableExists {
		return fmt.Errorf("users table does not exist")
	}

	log.Println("✅ Users table exists")
	return nil
}

// logConnectionStats периодически логирует статистику подключений
func logConnectionStats(db *sqlx.DB) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		stats := db.Stats()
		log.Printf("📊 Database connection stats: "+
			"OpenConnections=%d, InUse=%d, Idle=%d, WaitCount=%d",
			stats.OpenConnections,
			stats.InUse,
			stats.Idle,
			stats.WaitCount,
		)
	}
}

// RunMigrations выполняет миграции базы данных
func RunMigrations(db *sqlx.DB, migrationDir string) error {
	log.Println("🔧 Checking for database migrations...")

	// Проверяем существование таблицы миграций
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Здесь можно добавить логику применения миграций из файлов
	// Пока просто возвращаем успех
	log.Println("✅ Migrations check complete")
	return nil
}

// HealthCheck проверяет здоровье базы данных
func HealthCheck(db *sqlx.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result int
	err := db.GetContext(ctx, &result, "SELECT 1")
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

// Transaction выполняет операции в транзакции
func Transaction(db *sqlx.DB, fn func(*sqlx.Tx) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction failed (rollback error: %v): %w", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// CloseConnection закрывает подключение к базе данных
func CloseConnection(db *sqlx.DB) error {
	if db == nil {
		return nil
	}

	log.Println("🔌 Closing database connection...")

	// Закрываем все соединения
	stats := db.Stats()
	log.Printf("📊 Final connection stats: OpenConnections=%d", stats.OpenConnections)

	if err := db.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	log.Println("✅ Database connection closed")
	return nil
}
