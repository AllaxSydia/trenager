<template>
  <div class="statistics-page">
    <div class="container">
      <h1 class="page-title">Статистика студентов</h1>
      
      <div v-if="loading" class="loading">
        <p>Загрузка данных...</p>
      </div>
      
      <div v-else-if="error" class="error">
        <p>{{ error }}</p>
        <button @click="loadStatistics" class="retry-btn">Повторить</button>
      </div>
      
      <div v-else class="statistics-content">
        <div class="summary-cards">
          <div class="summary-card">
            <h3>Всего студентов</h3>
            <p class="summary-value">{{ statistics.total_students }}</p>
          </div>
          <div class="summary-card">
            <h3>Всего решений</h3>
            <p class="summary-value">{{ statistics.total_solutions }}</p>
          </div>
        </div>
        
        <div class="students-table">
          <table>
            <thead>
              <tr>
                <th>Студент</th>
                <th>Email</th>
                <th>Всего задач</th>
                <th>Решено</th>
                <th>Успешность</th>
                <th>По языкам</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="student in statistics.students" :key="student.user_id">
                <td>{{ student.username }}</td>
                <td>{{ student.email }}</td>
                <td>{{ student.total_tasks }}</td>
                <td>{{ student.solved_tasks }}</td>
                <td>
                  <div class="progress-bar">
                    <div 
                      class="progress-fill" 
                      :style="{ width: student.success_rate + '%' }"
                    ></div>
                    <span class="progress-text">{{ student.success_rate.toFixed(1) }}%</span>
                  </div>
                </td>
                <td>
                  <div class="languages">
                    <span 
                      v-for="(stats, lang) in student.languages" 
                      :key="lang"
                      class="language-badge"
                    >
                      {{ lang }}: {{ stats.solved_tasks }}/{{ stats.total_tasks }}
                    </span>
                    <span v-if="Object.keys(student.languages).length === 0" class="no-data">
                      Нет данных
                    </span>
                  </div>
                </td>
              </tr>
              <tr v-if="statistics.students.length === 0">
                <td colspan="6" class="no-students">Нет студентов</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Statistics',
  data() {
    return {
      loading: true,
      error: null,
      statistics: {
        total_students: 0,
        total_solutions: 0,
        students: []
      }
    }
  },
  mounted() {
    this.checkAuth()
    this.loadStatistics()
  },
  methods: {
    checkAuth() {
      const user = JSON.parse(localStorage.getItem('user') || '{}')
      if (!user.isLoggedIn || user.role !== 'teacher') {
        this.$router.push('/')
      }
    },
    async loadStatistics() {
      this.loading = true
      this.error = null
      
      try {
        const token = localStorage.getItem('token')
        const response = await fetch('/api/admin/statistics', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        
        if (!response.ok) {
          if (response.status === 403) {
            this.error = 'Доступ запрещен. Только для преподавателей.'
            this.$router.push('/')
            return
          }
          throw new Error('Ошибка при загрузке статистики')
        }
        
        const data = await response.json()
        this.statistics = data
      } catch (err) {
        this.error = err.message || 'Не удалось загрузить статистику'
        console.error('Ошибка загрузки статистики:', err)
      } finally {
        this.loading = false
      }
    }
  }
}
</script>

<style scoped>
.statistics-page {
  min-height: 100vh;
  background: #0E1117;
  padding: 120px 20px 40px;
  color: #E5E7EB;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
}

.page-title {
  font-size: 2.5rem;
  font-weight: bold;
  margin-bottom: 2rem;
  color: #E5E7EB;
}

.loading, .error {
  text-align: center;
  padding: 3rem;
  color: #9CA3AF;
}

.error {
  color: #EF4444;
}

.retry-btn {
  margin-top: 1rem;
  padding: 10px 20px;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 16px;
}

.retry-btn:hover {
  background: #2563EB;
}

.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.summary-card {
  background: #1F2937;
  border: 1px solid #374151;
  border-radius: 12px;
  padding: 1.5rem;
}

.summary-card h3 {
  font-size: 1rem;
  color: #9CA3AF;
  margin-bottom: 0.5rem;
}

.summary-value {
  font-size: 2.5rem;
  font-weight: bold;
  color: #3B82F6;
}

.students-table {
  background: #1F2937;
  border: 1px solid #374151;
  border-radius: 12px;
  overflow: hidden;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  background: #374151;
}

th {
  padding: 1rem;
  text-align: left;
  font-weight: 600;
  color: #E5E7EB;
  border-bottom: 2px solid #4B5563;
}

td {
  padding: 1rem;
  border-bottom: 1px solid #374151;
  color: #9CA3AF;
}

tbody tr:hover {
  background: #374151;
}

.progress-bar {
  position: relative;
  width: 100px;
  height: 24px;
  background: #374151;
  border-radius: 12px;
  overflow: hidden;
}

.progress-fill {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background: #10B981;
  transition: width 0.3s ease;
}

.progress-text {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 12px;
  font-weight: 600;
  color: #E5E7EB;
  z-index: 1;
}

.languages {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.language-badge {
  background: #3B82F6;
  color: #E5E7EB;
  padding: 4px 8px;
  border-radius: 6px;
  font-size: 12px;
}

.no-data {
  color: #6B7280;
  font-style: italic;
}

.no-students {
  text-align: center;
  padding: 2rem;
  color: #6B7280;
}

@media (max-width: 768px) {
  .page-title {
    font-size: 2rem;
  }
  
  .summary-cards {
    grid-template-columns: 1fr;
  }
  
  table {
    font-size: 14px;
  }
  
  th, td {
    padding: 0.75rem 0.5rem;
  }
  
  .languages {
    flex-direction: column;
  }
}
</style>

