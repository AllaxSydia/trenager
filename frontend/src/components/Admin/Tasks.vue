<template>
  <div class="teacher-tasks-page">
    <!-- Хедер -->
    <div class="page-header">
      <div class="header-content">
        <h1>📝 Создание задач</h1>
        <p class="subtitle">Панель учителя для создания и управления заданиями</p>
        
        <div class="header-actions">
          <button @click="goBack" class="btn btn-secondary">
            ← Назад
          </button>
          <div class="user-info">
            <span class="username">{{ userEmail }}</span>
            <span class="role-badge">👨‍🏫 Учитель</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Основной контент -->
    <div class="tasks-container">
      <!-- Левая колонка - создание задачи -->
      <div class="create-task-section">
        <div class="section-card">
          <h2>Создать новую задачу</h2>
          
          <!-- Выбор языка -->
          <div class="form-group">
            <label class="form-label">Язык программирования</label>
            <div class="language-selector">
              <button 
                v-for="lang in languages" 
                :key="lang.id"
                @click="selectLanguage(lang.id)"
                :class="['lang-btn', { 'lang-btn--active': newTask.language === lang.id }]"
              >
                {{ lang.name }}
              </button>
            </div>
          </div>

          <!-- Название задачи -->
          <div class="form-group">
            <label class="form-label">Название задачи</label>
            <input 
              v-model="newTask.title" 
              type="text" 
              placeholder="Введите название задачи"
              class="form-input"
            >
          </div>

          <!-- Описание задачи -->
          <div class="form-group">
            <label class="form-label">Описание задачи</label>
            <textarea 
              v-model="newTask.description" 
              placeholder="Опишите задачу для студентов..."
              class="form-textarea"
              rows="5"
            ></textarea>
          </div>

          <!-- Стартовый код -->
          <div class="form-group">
            <label class="form-label">Стартовый код</label>
            <div class="code-editor-wrapper">
              <div class="editor-header">
                <span class="editor-title">{{ getLanguageName(newTask.language) }}</span>
              </div>
              <textarea 
                v-model="newTask.starterCode" 
                placeholder="Введите начальный код..."
                class="code-editor"
                spellcheck="false"
              ></textarea>
            </div>
          </div>

          <!-- Тесты -->
          <div class="form-group">
            <div class="tests-header">
              <label class="form-label">Тесты для проверки</label>
              <button @click="addTest" class="btn btn-sm btn-success">
                + Добавить тест
              </button>
            </div>
            
            <div v-if="newTask.tests.length === 0" class="no-tests">
              <p>Нет добавленных тестов. Нажмите "Добавить тест" чтобы создать первый.</p>
            </div>
            
            <div v-else class="tests-list">
              <div v-for="(test, index) in newTask.tests" :key="index" class="test-item">
                <div class="test-header">
                  <h4>Тест {{ index + 1 }}</h4>
                  <button 
                    @click="removeTest(index)"
                    class="btn btn-sm btn-danger"
                    title="Удалить тест"
                  >
                    ×
                  </button>
                </div>
                
                <div class="test-inputs">
                  <div class="input-group">
                    <label>Входные данные:</label>
                    <textarea 
                      v-model="test.input" 
                      placeholder="Введите входные данные (каждое значение с новой строки)"
                      class="test-input"
                      rows="3"
                    ></textarea>
                  </div>
                  
                  <div class="input-group">
                    <label>Ожидаемый вывод:</label>
                    <textarea 
                      v-model="test.expectedOutput" 
                      placeholder="Ожидаемый результат"
                      class="test-input"
                      rows="2"
                    ></textarea>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Кнопки действий -->
          <div class="form-actions">
            <button 
              @click="saveTask" 
              :disabled="!isFormValid"
              class="btn btn-primary btn-large"
            >
              💾 Сохранить задачу
            </button>
            
            <button 
              @click="resetForm" 
              class="btn btn-outline"
            >
              🗑️ Очистить форму
            </button>
            
            <button 
              @click="previewTask" 
              :disabled="!isFormValid"
              class="btn btn-info"
            >
              👁️ Предпросмотр
            </button>
          </div>
        </div>
      </div>

      <!-- Правая колонка - существующие задачи -->
      <div class="existing-tasks-section">
        <div class="section-card">
          <h2>Существующие задачи</h2>
          
          <div class="tasks-filter">
            <select v-model="filterLanguage" class="filter-select">
              <option value="">Все языки</option>
              <option v-for="lang in languages" :key="lang.id" :value="lang.id">
                {{ lang.name }}
              </option>
            </select>
            
            <div class="tasks-count">
              {{ filteredTasks.length }} задач
            </div>
          </div>

          <div v-if="filteredTasks.length === 0" class="no-tasks">
            <p>Нет созданных задач</p>
          </div>

          <div v-else class="tasks-list">
            <div 
              v-for="task in filteredTasks" 
              :key="task.id"
              class="task-card"
              :class="{ 'task-card--selected': selectedTaskId === task.id }"
              @click="selectTask(task)"
            >
              <div class="task-header">
                <span class="task-language">{{ getLanguageName(task.language) }}</span>
                <span class="task-date">{{ formatDate(task.createdAt) }}</span>
              </div>
              
              <h3 class="task-title">{{ task.title }}</h3>
              
              <p class="task-description">{{ truncateText(task.description, 100) }}</p>
              
              <div class="task-footer">
                <span class="task-tests">🧪 {{ task.tests?.length || 0 }} тестов</span>
                <div class="task-actions">
                  <button 
                    @click.stop="editTask(task)"
                    class="btn-icon"
                    title="Редактировать"
                  >
                    ✏️
                  </button>
                  <button 
                    @click.stop="deleteTask(task.id)"
                    class="btn-icon btn-icon--danger"
                    title="Удалить"
                  >
                    🗑️
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Модалка предпросмотра -->
    <div v-if="showPreview" class="modal-overlay" @click="showPreview = false">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h2>Предпросмотр задачи</h2>
          <button @click="showPreview = false" class="btn-close">×</button>
        </div>
        
        <div class="modal-body">
          <h3>{{ newTask.title || 'Без названия' }}</h3>
          <p class="preview-language">Язык: {{ getLanguageName(newTask.language) }}</p>
          
          <div class="preview-section">
            <h4>Описание:</h4>
            <pre class="preview-description">{{ newTask.description || 'Нет описания' }}</pre>
          </div>
          
          <div class="preview-section">
            <h4>Стартовый код:</h4>
            <pre class="preview-code">{{ newTask.starterCode || 'Нет стартового кода' }}</pre>
          </div>
          
          <div class="preview-section">
            <h4>Тесты ({{ newTask.tests.length }}):</h4>
            <div v-if="newTask.tests.length === 0" class="no-tests">
              <p>Нет тестов</p>
            </div>
            <div v-else class="preview-tests">
              <div v-for="(test, index) in newTask.tests" :key="index" class="preview-test">
                <strong>Тест {{ index + 1 }}:</strong>
                <div class="test-preview">
                  <div><strong>Вход:</strong> {{ test.input || '(пусто)' }}</div>
                  <div><strong>Ожидаемый вывод:</strong> {{ test.expectedOutput || '(пусто)' }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <div class="modal-footer">
          <button @click="showPreview = false" class="btn btn-secondary">Закрыть</button>
          <button @click="saveTask" class="btn btn-primary">Сохранить задачу</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'TeacherTasks',

  props: {
    editTaskId: {
      type: String,
      default: null
    }
  },
  
  data() {
    return {
      // Данные новой задачи
      newTask: {
        language: 'python',
        title: '',
        description: '',
        starterCode: '',
        tests: []
      },
      
      // Существующие задачи (будут загружаться из API)
      existingTasks: [],
      
      // UI состояние
      languages: [
        { id: 'python', name: 'Python' },
        { id: 'javascript', name: 'JavaScript' },
        { id: 'java', name: 'Java' },
        { id: 'cpp', name: 'C++' }
      ],
      filterLanguage: '',
      selectedTaskId: null,
      showPreview: false,
      userEmail: ''
    }
  },
  
  computed: {
    // Валидация формы
    isFormValid() {
      return (
        this.newTask.title.trim() !== '' &&
        this.newTask.description.trim() !== '' &&
        this.newTask.tests.length > 0
      )
    },
    
    // Фильтрация задач
    filteredTasks() {
      if (!this.filterLanguage) return this.existingTasks
      return this.existingTasks.filter(task => task.language === this.filterLanguage)
    }
  },
  
  mounted() {
    this.loadUserData()
    this.loadExistingTasksFromAPI() // Загрузим реальные задачи
    this.checkEditMode()
  },
  
  methods: {
    checkEditMode() {
      // Если есть ID задачи для редактирования из props
      if (this.editTaskId) {
        this.loadTaskForEdit(this.editTaskId)
      }
    },

    async loadExistingTasksFromAPI() {
      try {
        const token = localStorage.getItem('token')
        const response = await fetch('/api/teacher/tasks', {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        
        if (response.ok) {
          const tasks = await response.json()
          console.log('Сырые данные от API:', tasks)
          
          // Конвертируем названия полей для отображения
          this.existingTasks = tasks.map(task => {
            // Обрабатываем тесты
            const convertedTests = (task.tests || []).map(test => {
              // Если в БД поле называется expected_output, добавляем expectedOutput для фронтенда
              return {
                input: test.input || '',
                expectedOutput: test.expected_output || test.expectedOutput || '',
                expected_output: test.expected_output || test.expectedOutput || '' // Оставляем для совместимости
              }
            })
            
            return {
              ...task,
              tests: convertedTests
            }
          })
          
          console.log('Загруженные и конвертированные задачи:', this.existingTasks)
        } else {
          console.warn('Не удалось загрузить задачи, используем демо-данные')
        }
      } catch (error) {
        console.error('Ошибка загрузки задач:', error)
      }
    },

    async loadTaskForEdit(taskId) {
      try {
        const token = localStorage.getItem('token')
        const response = await fetch(`/api/teacher/tasks/${taskId}`, {
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        
        if (response.ok) {
          const task = await response.json()
          this.editTask(task)
        }
      } catch (error) {
        console.error('Ошибка загрузки задачи:', error)
      }
    },

    loadUserData() {
      const savedUser = JSON.parse(localStorage.getItem('user') || '{}')
      this.userEmail = savedUser.email || 'Неизвестный'
    },
    
    // Навигация
    goBack() {
      this.$router.back()
    },
    
    // Выбор языка
    selectLanguage(langId) {
      this.newTask.language = langId
    },
    
    getLanguageName(langId) {
      const lang = this.languages.find(l => l.id === langId)
      return lang ? lang.name : 'Неизвестный'
    },
    
    // Работа с тестами
    addTest() {
      this.newTask.tests.push({
        input: '',
        expectedOutput: ''
      })
    },
    
    removeTest(index) {
      this.newTask.tests.splice(index, 1)
    },
    
    // Форматирование
    formatDate(dateString) {
      const date = new Date(dateString)
      return date.toLocaleDateString('ru-RU')
    },
    
    truncateText(text, length) {
      if (!text) return ''
      return text.length > length ? text.substring(0, length) + '...' : text
    },
    
    // Действия с задачами
    async saveTask() {
      if (!this.isFormValid) {
        alert('Заполните все обязательные поля: название, описание и хотя бы один тест!')
        return
      }
      
      try {
        const token = localStorage.getItem('token')
        const url = this.selectedTaskId && !isNaN(this.selectedTaskId) 
          ? `/api/teacher/tasks/${this.selectedTaskId}`
          : '/api/teacher/tasks'
        
        const method = this.selectedTaskId && !isNaN(this.selectedTaskId) ? 'PUT' : 'POST'
        
        // ОТЛАДКА: посмотрим что у нас в тестах
        console.log('Исходные тесты (frontend формат):', this.newTask.tests)
        
        // ФОРМАТИРУЕМ ТЕСТЫ ПРАВИЛЬНО ДЛЯ БЭКЕНДА
        const formattedTests = this.newTask.tests.map(test => ({
          input: test.input || '',
          expected_output: test.expectedOutput || '' // ← конвертируем в snake_case
        }))
        
        console.log('Отформатированные тесты (backend формат):', formattedTests)
        
        const taskData = {
          title: this.newTask.title,
          description: this.newTask.description,
          language: this.newTask.language,
          starter_code: this.newTask.starterCode,
          tests: formattedTests, // ← используем отформатированные тесты
          difficulty: 'beginner',
          is_published: true
        }
        
        console.log('Отправляемые данные:', taskData)
        
        const response = await fetch(url, {
          method,
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${token}`
          },
          body: JSON.stringify(taskData)
        })
        
        if (response.ok) {
          const savedTask = await response.json()
          console.log('Ответ от сервера:', savedTask)
          alert(`Задача "${this.newTask.title}" сохранена!`)
          
          // Обновляем список задач
          await this.loadExistingTasksFromAPI()
          
          this.resetForm()
          this.showPreview = false
          
          // Если мы в режиме редактирования через роутер - возвращаемся
          if (this.editTaskId) {
            this.$router.push('/teacher/tasks')
          }
        } else {
          const errorText = await response.text()
          console.error('Ошибка сохранения:', errorText)
          throw new Error(`Ошибка сохранения: ${response.status}`)
        }
      } catch (error) {
        console.error('Ошибка сохранения задачи:', error)
        alert('Не удалось сохранить задачу: ' + error.message)
      }
    },
    
    resetForm() {
      this.newTask = {
        language: 'python',
        title: '',
        description: '',
        starterCode: '',
        tests: []
      }
      this.selectedTaskId = null
    },
    
    previewTask() {
      if (!this.isFormValid) {
        alert('Заполните все обязательные поля для предпросмотра!')
        return
      }
      this.showPreview = true
    },
    
    selectTask(task) {
      this.selectedTaskId = task.id
    },
    
    editTask(task) {
      console.log('Редактируем задачу:', task)
      console.log('Тесты задачи:', task.tests)
      
      // Конвертируем тесты из формата бэкенда в формат фронтенда
      const convertedTests = (task.tests || []).map(test => {
        console.log('Тест:', test)
        return {
          input: test.input || '',
          expectedOutput: test.expected_output || test.expectedOutput || '' // Обрабатываем оба варианта
        }
      })
      
      console.log('Конвертированные тесты:', convertedTests)
      
      this.newTask = {
        language: task.language,
        title: task.title,
        description: task.description,
        starterCode: task.starter_code || task.starterCode || '',
        tests: convertedTests
      }
      this.selectedTaskId = task.id
      
      console.log('Новые данные формы:', this.newTask)
      
      // Скролл к форме
      document.querySelector('.create-task-section')?.scrollIntoView({ behavior: 'smooth' })
    },
    
    async deleteTask(taskId) {
      if (!confirm('Вы уверены, что хотите удалить эту задачу?')) return
      
      try {
        const token = localStorage.getItem('token')
        const response = await fetch(`/api/teacher/tasks/${taskId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${token}`
          }
        })
        
        if (response.ok) {
          // Удаляем из списка
          const index = this.existingTasks.findIndex(t => t.id == taskId)
          if (index !== -1) {
            this.existingTasks.splice(index, 1)
          }
          
          if (this.selectedTaskId == taskId) {
            this.resetForm()
          }
          
          alert('Задача удалена')
        } else {
          throw new Error('Ошибка удаления')
        }
      } catch (error) {
        console.error('Ошибка удаления задачи:', error)
        alert('Не удалось удалить задачу')
      }
    }
  }
}
</script>

<style scoped>
.teacher-tasks-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  color: #e2e8f0;
  padding: 20px;
}

/* Хедер */
.page-header {
  background: rgba(30, 41, 59, 0.8);
  border-radius: 12px;
  padding: 1.5rem 2rem;
  margin-bottom: 2rem;
  border: 1px solid #334155;
}

.header-content h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
  background: linear-gradient(135deg, #60a5fa, #3b82f6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  color: #94a3b8;
  margin: 0.5rem 0 1.5rem;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.username {
  color: #cbd5e1;
  font-weight: 500;
}

.role-badge {
  background: linear-gradient(135deg, #8b5cf6, #7c3aed);
  color: white;
  padding: 0.25rem 0.75rem;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
}

/* Основной контейнер */
.tasks-container {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  max-width: 1600px;
  margin: 0 auto;
}

/* Карточки */
.section-card {
  background: rgba(30, 41, 59, 0.8);
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid #334155;
  height: fit-content;
}

.section-card h2 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #f8fafc;
  font-size: 1.5rem;
}

/* Форма */
.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  color: #cbd5e1;
  font-weight: 500;
}

/* Выбор языка */
.language-selector {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.lang-btn {
  background: #1e293b;
  border: 1px solid #475569;
  color: #94a3b8;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.lang-btn:hover {
  background: #334155;
  color: #e2e8f0;
}

.lang-btn--active {
  background: #3b82f6;
  border-color: #3b82f6;
  color: white;
}

/* Поля ввода */
.form-input,
.form-textarea,
.code-editor {
  width: 100%;
  background: #0f172a;
  border: 1px solid #475569;
  color: #e2e8f0;
  border-radius: 6px;
  padding: 0.75rem;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 0.9rem;
  transition: border-color 0.2s;
}

.form-input:focus,
.form-textarea:focus,
.code-editor:focus {
  outline: none;
  border-color: #3b82f6;
}

.form-textarea,
.code-editor {
  resize: vertical;
  min-height: 100px;
}

.code-editor-wrapper {
  background: #0f172a;
  border-radius: 6px;
  overflow: hidden;
}

.editor-header {
  background: #1e293b;
  padding: 0.5rem 0.75rem;
  border-bottom: 1px solid #334155;
}

.editor-title {
  color: #94a3b8;
  font-size: 0.85rem;
  font-weight: 500;
}

/* Тесты */
.tests-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.no-tests {
  background: rgba(255, 255, 255, 0.05);
  border: 1px dashed #475569;
  border-radius: 6px;
  padding: 1.5rem;
  text-align: center;
  color: #94a3b8;
}

.tests-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.test-item {
  background: #0f172a;
  border: 1px solid #475569;
  border-radius: 8px;
  padding: 1rem;
}

.test-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.test-header h4 {
  margin: 0;
  color: #cbd5e1;
}

.test-inputs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.input-group label {
  font-size: 0.9rem;
  color: #94a3b8;
}

.test-input {
  width: 100%;
  background: #1e293b;
  border: 1px solid #334155;
  color: #e2e8f0;
  border-radius: 4px;
  padding: 0.5rem;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 0.85rem;
  resize: vertical;
}

/* Кнопки действий формы */
.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #334155;
}

/* Существующие задачи */
.tasks-filter {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.filter-select {
  background: #0f172a;
  border: 1px solid #475569;
  color: #e2e8f0;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  min-width: 150px;
}

.tasks-count {
  color: #94a3b8;
  font-size: 0.9rem;
}

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-height: 600px;
  overflow-y: auto;
  padding-right: 0.5rem;
}

.task-card {
  background: #0f172a;
  border: 1px solid #475569;
  border-radius: 8px;
  padding: 1rem;
  cursor: pointer;
  transition: all 0.2s;
}

.task-card:hover {
  border-color: #60a5fa;
  transform: translateY(-1px);
}

.task-card--selected {
  border-color: #3b82f6;
  background: rgba(59, 130, 246, 0.05);
}

.task-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
}

.task-language {
  background: #3b82f6;
  color: white;
  padding: 0.2rem 0.6rem;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 600;
}

.task-date {
  color: #94a3b8;
  font-size: 0.85rem;
}

.task-title {
  margin: 0 0 0.5rem;
  color: #f8fafc;
  font-size: 1.1rem;
}

.task-description {
  color: #94a3b8;
  font-size: 0.9rem;
  line-height: 1.4;
  margin-bottom: 0.75rem;
}

.task-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.task-tests {
  color: #60a5fa;
  font-size: 0.85rem;
  font-weight: 500;
}

.task-actions {
  display: flex;
  gap: 0.5rem;
}

.btn-icon {
  background: none;
  border: none;
  color: #94a3b8;
  cursor: pointer;
  font-size: 1.1rem;
  padding: 0.25rem;
  border-radius: 4px;
  transition: all 0.2s;
}

.btn-icon:hover {
  color: #e2e8f0;
  background: rgba(255, 255, 255, 0.05);
}

.btn-icon--danger:hover {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
}

/* Модалка */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: #1e293b;
  border-radius: 12px;
  width: 90%;
  max-width: 800px;
  max-height: 90vh;
  overflow-y: auto;
  border: 1px solid #475569;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid #334155;
}

.modal-header h2 {
  margin: 0;
  color: #f8fafc;
}

.btn-close {
  background: none;
  border: none;
  color: #94a3b8;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
}

.btn-close:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #e2e8f0;
}

.modal-body {
  padding: 1.5rem;
}

.modal-body h3 {
  margin: 0 0 0.5rem;
  color: #f8fafc;
}

.preview-language {
  color: #60a5fa;
  margin-bottom: 1.5rem;
}

.preview-section {
  margin-bottom: 1.5rem;
}

.preview-section h4 {
  color: #cbd5e1;
  margin: 0 0 0.5rem;
}

.preview-description,
.preview-code {
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 6px;
  padding: 1rem;
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
  font-family: inherit;
}

.preview-code {
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 0.9rem;
}

.preview-tests {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.preview-test {
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 6px;
  padding: 0.75rem;
}

.test-preview {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  margin-top: 0.25rem;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 0.9rem;
}

.modal-footer {
  padding: 1.5rem;
  border-top: 1px solid #334155;
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

/* Кнопки */
.btn {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.9rem;
}

.btn-sm {
  padding: 0.25rem 0.5rem;
  font-size: 0.8rem;
}

.btn-large {
  padding: 0.75rem 1.5rem;
  font-size: 1rem;
}

.btn-primary {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #2563eb, #1d4ed8);
  transform: translateY(-1px);
}

.btn-primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: #475569;
  color: white;
}

.btn-secondary:hover {
  background: #334155;
}

.btn-success {
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
}

.btn-success:hover {
  background: linear-gradient(135deg, #059669, #047857);
}

.btn-danger {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
}

.btn-danger:hover {
  background: linear-gradient(135deg, #dc2626, #b91c1c);
}

.btn-info {
  background: linear-gradient(135deg, #06b6d4, #0891b2);
  color: white;
}

.btn-info:hover {
  background: linear-gradient(135deg, #0891b2, #0e7490);
}

.btn-outline {
  background: transparent;
  border: 1px solid #475569;
  color: #94a3b8;
}

.btn-outline:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #e2e8f0;
}

/* Адаптация */
@media (max-width: 1200px) {
  .tasks-container {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .teacher-tasks-page {
    padding: 10px;
  }
  
  .page-header {
    padding: 1rem;
  }
  
  .section-card {
    padding: 1rem;
  }
  
  .test-inputs {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .modal-content {
    width: 95%;
  }
}

@media (max-width: 480px) {
  .header-actions {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
  
  .language-selector {
    justify-content: center;
  }
  
  .lang-btn {
    flex: 1;
    min-width: 80px;
  }
}
</style>