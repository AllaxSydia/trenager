package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var DB *sql.DB

func Init() {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "db"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "postgres"
	}
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "trenager"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Ошибка открытия БД: %v", err)
	}

	for i := 0; i < 10; i++ {
		if err = DB.Ping(); err == nil {
			break
		}
		log.Printf("⏳ Ожидание PostgreSQL... (%d/10)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("❌ Не удалось подключиться к PostgreSQL: %v", err)
	}

	log.Println("✅ Соединение с PostgreSQL успешно")
	createUsersTable()
	createTaskSolutionsTable()
	createDefaultUsers()
}

func createUsersTable() {
	// Сначала добавляем колонку role, если её нет
	alterQuery := `
	DO $$ 
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
			WHERE table_name='users' AND column_name='role') THEN
			ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'student';
		END IF;
	END $$;
	`
	_, err := DB.Exec(alterQuery)
	if err != nil {
		log.Printf("⚠️ Предупреждение при добавлении колонки role: %v", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(100) UNIQUE NOT NULL,
		email VARCHAR(150) UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		role VARCHAR(20) DEFAULT 'student',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(query)
	if err != nil {
		log.Fatalf("❌ Ошибка при создании таблицы users: %v", err)
	}
	log.Println("✅ Таблица users готова")
}

func createTaskSolutionsTable() {
	query := `
	CREATE TABLE IF NOT EXISTS task_solutions (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		task_id VARCHAR(255) NOT NULL,
		language VARCHAR(50) NOT NULL,
		code TEXT NOT NULL,
		success BOOLEAN DEFAULT FALSE,
		passed_tests INTEGER DEFAULT 0,
		total_tests INTEGER DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, task_id, language)
	);
	CREATE INDEX IF NOT EXISTS idx_task_solutions_user_id ON task_solutions(user_id);
	CREATE INDEX IF NOT EXISTS idx_task_solutions_task_id ON task_solutions(task_id);
	CREATE INDEX IF NOT EXISTS idx_task_solutions_created_at ON task_solutions(created_at);
	`
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("❌ Ошибка при создании таблицы task_solutions: %v", err)
	}
	log.Println("✅ Таблица task_solutions готова")
}

func createDefaultUsers() {
	// Создаём тестовых пользователей, если их ещё нет
	defaultUsers := []struct {
		username string
		email    string
		password string
		role     string
	}{
		{
			username: "teacher",
			email:    "teacher@mail.com",
			password: "123456789",
			role:     "teacher",
		},
		{
			username: "student",
			email:    "student@mail.com",
			password: "123456789",
			role:     "student",
		},
	}

	for _, user := range defaultUsers {
		// Проверяем, существует ли пользователь
		var exists bool
		err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.email).Scan(&exists)
		if err != nil {
			log.Printf("⚠️ Ошибка при проверке пользователя %s: %v", user.email, err)
			continue
		}

		if exists {
			// Обновляем роль, если пользователь существует, но роль неправильная
			_, err = DB.Exec("UPDATE users SET role = $1 WHERE email = $2", user.role, user.email)
			if err != nil {
				log.Printf("⚠️ Ошибка при обновлении роли для %s: %v", user.email, err)
			}
			continue
		}

		// Хэшируем пароль
		hash, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("⚠️ Ошибка при хэшировании пароля для %s: %v", user.email, err)
			continue
		}

		// Создаём пользователя
		query := `
			INSERT INTO users (username, email, password_hash, role, created_at, updated_at)
			VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			ON CONFLICT (email) DO UPDATE SET
				role = EXCLUDED.role,
				password_hash = EXCLUDED.password_hash,
				updated_at = CURRENT_TIMESTAMP
		`
		_, err = DB.Exec(query, user.username, user.email, string(hash), user.role)
		if err != nil {
			log.Printf("⚠️ Ошибка при создании пользователя %s: %v", user.email, err)
		} else {
			log.Printf("✅ Тестовый пользователь создан/обновлён: %s (%s)", user.email, user.role)
		}
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
