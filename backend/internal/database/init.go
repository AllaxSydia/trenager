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

	// Настройка пула соединений
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	for i := 0; i < 15; i++ {
		if err = DB.Ping(); err == nil {
			break
		}
		log.Printf("⏳ Ожидание PostgreSQL... (%d/15)", i+1)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatalf("❌ Не удалось подключиться к PostgreSQL: %v", err)
	}

	log.Println("✅ Соединение с PostgreSQL успешно")
	createUsersTable()
	createTaskSolutionsTable()
	createTasksTable()
	fixTasksTable()
	createDefaultUsers()
	createSampleTasks()
}

func createUsersTable() {
	// Сначала создаем таблицу с правильными колонками
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
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("❌ Ошибка при создании таблицы users: %v", err)
	}

	// Добавляем колонку updated_at, если её нет
	addColumnQuery := `
	DO $$ 
	BEGIN
		IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
			WHERE table_name='users' AND column_name='updated_at') THEN
			ALTER TABLE users ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
		END IF;
	END $$;
	`
	_, err = DB.Exec(addColumnQuery)
	if err != nil {
		log.Printf("⚠️ Предупреждение при добавлении колонки updated_at: %v", err)
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

func createTasksTable() {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        language VARCHAR(50) NOT NULL,
        title VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        difficulty VARCHAR(20) DEFAULT 'beginner',
        template TEXT,                           -- Добавляем
        starter_code TEXT,                      -- Уже есть
        tests JSONB NOT NULL,
        created_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        is_published BOOLEAN DEFAULT TRUE
    );
    
    -- Создаем индексы отдельно
    CREATE INDEX IF NOT EXISTS idx_tasks_language ON tasks(language);
    CREATE INDEX IF NOT EXISTS idx_tasks_created_by ON tasks(created_by);
    CREATE INDEX IF NOT EXISTS idx_tasks_is_published ON tasks(is_published);
    CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);
    `

	_, err := DB.Exec(query)
	if err != nil {
		log.Printf("❌ Ошибка при создании таблицы tasks: %v", err)
	} else {
		log.Println("✅ Таблица tasks готова")
	}

	// Добавляем колонку template если её нет
	addTemplateQuery := `
    DO $$ 
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
            WHERE table_name='tasks' AND column_name='template') THEN
            ALTER TABLE tasks ADD COLUMN template TEXT;
        END IF;
    END $$;
    `

	_, err = DB.Exec(addTemplateQuery)
	if err != nil {
		log.Printf("⚠️ Предупреждение при добавлении колонки template: %v", err)
	}
}

func createDefaultUsers() {
	log.Println("🔄 Начинаем создание тестовых пользователей...")

	// Создаём тестовых пользователей с РАЗНЫМИ email
	defaultUsers := []struct {
		username string
		email    string
		password string
		role     string
	}{
		// Обновляем существующего пользователя
		{
			username: "teacher_avg",
			email:    "teacher@mail.com",
			password: "123456789",
			role:     "teacher",
		},
		// Новые пользователи с УНИКАЛЬНЫМИ email
		{
			username: "student_ivan",
			email:    "student@trenager.ru",
			password: "123456789",
			role:     "student",
		},
		{
			username: "admin_root",
			email:    "admin@trenager.ru",
			password: "123456789",
			role:     "admin",
		},
		{
			username: "teacher_alex",
			email:    "alex@teacher.ru",
			password: "123456789",
			role:     "teacher",
		},
		{
			username: "student_olga",
			email:    "olga@student.ru",
			password: "studen123456789t123",
			role:     "student",
		},
	}

	for _, user := range defaultUsers {
		// Хэшируем пароль
		hash, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("⚠️ Ошибка хэширования для %s: %v", user.email, err)
			continue
		}

		// Используем UPSERT без updated_at (она будет установлена по умолчанию)
		query := `
			INSERT INTO users (username, email, password_hash, role, created_at)
			VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
			ON CONFLICT (email) 
			DO UPDATE SET 
				username = EXCLUDED.username,
				password_hash = EXCLUDED.password_hash,
				role = EXCLUDED.role
			RETURNING id
		`

		var userID int
		err = DB.QueryRow(query, user.username, user.email, string(hash), user.role).Scan(&userID)
		if err != nil {
			log.Printf("⚠️ Ошибка с пользователем %s: %v", user.email, err)
		} else {
			log.Printf("✅ Пользователь обработан: %s (%s) [ID: %d]", user.email, user.role, userID)
		}
	}

	// Выводим итоговый список пользователей
	log.Println("📋 Итоговый список пользователей:")
	rows, err := DB.Query("SELECT id, username, email, role FROM users ORDER BY id")
	if err != nil {
		log.Printf("⚠️ Ошибка при чтении пользователей: %v", err)
		return
	}
	defer rows.Close()

	userCount := 0
	for rows.Next() {
		var id int
		var username, email, role string
		if err := rows.Scan(&id, &username, &email, &role); err != nil {
			log.Printf("⚠️ Ошибка при сканировании: %v", err)
			continue
		}
		log.Printf("   %d: %s <%s> - %s", id, username, email, role)
		userCount++
	}

	if userCount == 0 {
		log.Println("⚠️ В базе нет пользователей!")
	} else {
		log.Printf("✅ Всего пользователей: %d", userCount)
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}

func fixTasksTable() {
	// Исправленный код для добавления колонок
	queries := []string{
		`ALTER TABLE tasks ADD COLUMN IF NOT EXISTS difficulty VARCHAR(20) DEFAULT 'beginner'`,
		`ALTER TABLE tasks ADD COLUMN IF NOT EXISTS template TEXT`,
		`ALTER TABLE tasks ADD COLUMN IF NOT EXISTS created_by INTEGER REFERENCES users(id) ON DELETE SET NULL`,
		`ALTER TABLE tasks ADD COLUMN IF NOT EXISTS is_published BOOLEAN DEFAULT TRUE`,
	}

	for _, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			log.Printf("⚠️ Ошибка при добавлении колонки: %v", err)
		} else {
			log.Printf("✅ Колонка добавлена/проверена")
		}
	}
}

func createSampleTasks() {
	// Проверяем, есть ли уже задачи
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		log.Printf("⚠️ Ошибка при проверке количества задач: %v", err)
		return
	}

	if count > 0 {
		log.Printf("✅ В БД уже есть %d задач, пропускаем создание тестовых", count)
		return
	}

	sampleTasks := []struct {
		title       string
		description string
		language    string
		difficulty  string
		template    string
		starterCode string
		tests       string
		createdBy   int
	}{
		{
			title:       "Hello World на Python",
			description: "Напишите программу, которая выводит 'Hello, World!'",
			language:    "python",
			difficulty:  "beginner",
			template:    `print("Hello, World!")`,
			starterCode: `print("Hello, World!")`,
			tests:       `[{"input": "", "expected_output": "Hello, World!"}]`,
			createdBy:   1, // teacher_avg
		},
		{
			title:       "Сумма двух чисел",
			description: "Напишите программу, которая принимает два числа через input() и выводит их сумму",
			language:    "python",
			difficulty:  "beginner",
			template:    `num1 = int(input())\nnum2 = int(input())\nprint(num1 + num2)`,
			starterCode: `num1 = int(input())\nnum2 = int(input())\nprint(num1 + num2)`,
			tests:       `[{"input": "5\\n3", "expected_output": "8"}, {"input": "10\\n20", "expected_output": "30"}]`,
			createdBy:   1,
		},
		{
			title:       "Факториал числа",
			description: "Напишите функцию для вычисления факториала числа",
			language:    "python",
			difficulty:  "intermediate",
			template:    `def factorial(n):\n    if n == 0:\n        return 1\n    result = 1\n    for i in range(1, n + 1):\n        result *= i\n    return result\n\n# Тестирование\nprint(factorial(5))`,
			starterCode: `def factorial(n):\n    # Ваш код здесь\n    pass\n\nprint(factorial(5))`,
			tests:       `[{"input": "5", "expected_output": "120"}, {"input": "0", "expected_output": "1"}]`,
			createdBy:   1,
		},
		{
			title:       "Проверка числа на четность",
			description: "Напишите программу, которая проверяет, является ли число четным",
			language:    "python",
			difficulty:  "beginner",
			template:    `num = int(input())\nif num % 2 == 0:\n    print("Четное")\nelse:\n    print("Нечетное")`,
			starterCode: `num = int(input())\n# Ваш код здесь`,
			tests:       `[{"input": "4", "expected_output": "Четное"}, {"input": "7", "expected_output": "Нечетное"}]`,
			createdBy:   1,
		},
		{
			title:       "Hello World на JavaScript",
			description: "Напишите программу, которая выводит 'Hello, World!'",
			language:    "javascript",
			difficulty:  "beginner",
			template:    `console.log("Hello, World!")`,
			starterCode: `console.log("Hello, World!")`,
			tests:       `[{"input": "", "expected_output": "Hello, World!"}]`,
			createdBy:   1,
		},
		{
			title:       "Сумма массивов",
			description: "Напишите функцию, которая суммирует все элементы массива",
			language:    "javascript",
			difficulty:  "intermediate",
			template:    `function sumArray(arr) {\n    return arr.reduce((a, b) => a + b, 0);\n}\n\nconsole.log(sumArray([1, 2, 3, 4, 5]));`,
			starterCode: `function sumArray(arr) {\n    // Ваш код здесь\n}\n\nconsole.log(sumArray([1, 2, 3, 4, 5]));`,
			tests:       `[{"input": "", "expected_output": "15"}, {"input": "", "expected_output": "0"}]`,
			createdBy:   1,
		},
	}

	successCount := 0
	for _, task := range sampleTasks {
		query := `
        INSERT INTO tasks (title, description, language, difficulty, template, 
                          starter_code, tests, created_by, is_published, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, true, CURRENT_TIMESTAMP)
        `

		_, err := DB.Exec(query,
			task.title,
			task.description,
			task.language,
			task.difficulty,
			task.template,
			task.starterCode,
			task.tests,
			task.createdBy,
		)

		if err != nil {
			log.Printf("⚠️ Ошибка при добавлении задачи '%s': %v", task.title, err)
		} else {
			log.Printf("✅ Добавлена тестовая задача: %s (%s)", task.title, task.language)
			successCount++
		}
	}

	log.Printf("📊 Всего добавлено %d тестовых задач", successCount)
}
