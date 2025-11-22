package handlers

import (
	"backend/internal/database"
	"encoding/json"
	"log"
	"net/http"
)

// StudentStatistics представляет статистику одного студента
type StudentStatistics struct {
	UserID      int64  `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	TotalTasks  int    `json:"total_tasks"`
	SolvedTasks int    `json:"solved_tasks"`
	SuccessRate float64 `json:"success_rate"`
	Languages   map[string]LanguageStats `json:"languages"`
}

// LanguageStats статистика по языку программирования
type LanguageStats struct {
	TotalTasks  int     `json:"total_tasks"`
	SolvedTasks int     `json:"solved_tasks"`
	SuccessRate float64 `json:"success_rate"`
}

// StatisticsResponse ответ со статистикой
type StatisticsResponse struct {
	TotalStudents   int                 `json:"total_students"`
	TotalSolutions  int                 `json:"total_solutions"`
	Students        []StudentStatistics `json:"students"`
}

// StatisticsHandler возвращает статистику по всем студентам
func StatisticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, `{"error":"method_not_allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Получаем всех студентов (не преподавателей)
	query := `
		SELECT u.id, u.username, u.email,
			COUNT(ts.id) as total_solutions,
			COUNT(CASE WHEN ts.success = true THEN 1 END) as solved_tasks
		FROM users u
		LEFT JOIN task_solutions ts ON u.id = ts.user_id
		WHERE u.role = 'student' OR u.role IS NULL
		GROUP BY u.id, u.username, u.email
		ORDER BY u.username
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("❌ Ошибка при получении статистики: %v", err)
		http.Error(w, `{"error":"database_error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []StudentStatistics
	totalSolutions := 0

	for rows.Next() {
		var student StudentStatistics
		var totalSolutionsCount, solvedTasksCount int

		err := rows.Scan(&student.UserID, &student.Username, &student.Email, &totalSolutionsCount, &solvedTasksCount)
		if err != nil {
			log.Printf("⚠️ Ошибка при сканировании строки: %v", err)
			continue
		}

		// Получаем статистику по языкам для этого студента
		langQuery := `
			SELECT language,
				COUNT(*) as total_tasks,
				COUNT(CASE WHEN success = true THEN 1 END) as solved_tasks
			FROM task_solutions
			WHERE user_id = $1
			GROUP BY language
		`
		langRows, err := database.DB.Query(langQuery, student.UserID)
		if err == nil {
			student.Languages = make(map[string]LanguageStats)
			for langRows.Next() {
				var lang string
				var langTotal, langSolved int
				if err := langRows.Scan(&lang, &langTotal, &langSolved); err == nil {
					successRate := 0.0
					if langTotal > 0 {
						successRate = float64(langSolved) / float64(langTotal) * 100
					}
					student.Languages[lang] = LanguageStats{
						TotalTasks:  langTotal,
						SolvedTasks: langSolved,
						SuccessRate: successRate,
					}
				}
			}
			langRows.Close()
		}

		student.TotalTasks = totalSolutionsCount
		student.SolvedTasks = solvedTasksCount
		if totalSolutionsCount > 0 {
			student.SuccessRate = float64(solvedTasksCount) / float64(totalSolutionsCount) * 100
		}

		totalSolutions += totalSolutionsCount
		students = append(students, student)
	}

	response := StatisticsResponse{
		TotalStudents:  len(students),
		TotalSolutions: totalSolutions,
		Students:       students,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("❌ Ошибка при кодировании ответа: %v", err)
		http.Error(w, `{"error":"encoding_error"}`, http.StatusInternalServerError)
		return
	}
}

