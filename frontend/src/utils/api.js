// frontend/src/utils/api.js
const API_BASE = '/api'

export const api = {
  // ============= АУТЕНТИФИКАЦИЯ =============
  
  /**
   * Вход в систему
   */
  async login(credentials) {
    try {
      const response = await fetch(`${API_BASE}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials)
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const data = await response.json()
      
      if (data.success && data.token) {
        // Сохраняем токен и данные пользователя
        localStorage.setItem('token', data.token)
        localStorage.setItem('user', JSON.stringify({
          id: data.userId,
          username: data.username,
          email: data.email,
          role: data.role || 'student',
          isLoggedIn: true
        }))
      }
      
      return data
    } catch (error) {
      console.error('API Login error:', error)
      return {
        success: false,
        error: `Connection error: ${error.message}`
      }
    }
  },

  /**
   * Быстрый вход для тестовых пользователей
   */
  async quickLogin(userType) {
    try {
      const response = await fetch(`${API_BASE}/auth/quick-login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ user_type: userType })
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const data = await response.json()
      
      if (data.success && data.token) {
        localStorage.setItem('token', data.token)
        localStorage.setItem('user', JSON.stringify({
          username: data.username,
          email: data.email,
          role: data.role,
          isLoggedIn: true
        }))
      }
      
      return data
    } catch (error) {
      console.error('API Quick Login error:', error)
      return {
        success: false,
        error: `Connection error: ${error.message}`
      }
    }
  },

  /**
   * Регистрация нового пользователя
   */
  async register(userData) {
    try {
      const response = await fetch(`${API_BASE}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData)
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const data = await response.json()
      
      if (data.success && data.token) {
        localStorage.setItem('token', data.token)
        localStorage.setItem('user', JSON.stringify({
          username: data.username,
          email: data.email,
          role: data.role || 'student',
          isLoggedIn: true
        }))
      }
      
      return data
    } catch (error) {
      console.error('API Register error:', error)
      return {
        success: false,
        error: `Connection error: ${error.message}`
      }
    }
  },

  /**
   * Гостевая авторизация
   */
  async guestLogin() {
    try {
      const response = await fetch(`${API_BASE}/auth/guest`)
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      const data = await response.json()
      
      if (data.success && data.token) {
        localStorage.setItem('token', data.token)
        localStorage.setItem('user', JSON.stringify({
          username: data.username,
          email: data.email,
          role: 'guest',
          isLoggedIn: true,
          guest: true
        }))
      }
      
      return data
    } catch (error) {
      console.error('API Guest Login error:', error)
      return {
        success: false,
        error: `Connection error: ${error.message}`
      }
    }
  },

  /**
   * Проверка валидности токена
   */
  async validateToken() {
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        return { success: false, error: 'No token' }
      }
      
      const response = await fetch(`${API_BASE}/auth/validate`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      
      if (!response.ok) {
        // Токен невалиден, удаляем
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        return { success: false, error: 'Invalid token' }
      }
      
      return await response.json()
    } catch (error) {
      console.error('API Validate Token error:', error)
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      return {
        success: false,
        error: `Connection error: ${error.message}`
      }
    }
  },

  /**
   * Получение информации о пользователе
   */
  async getUserInfo() {
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        return { success: false, error: 'No token' }
      }
      
      const response = await fetch(`${API_BASE}/auth/user-info`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      return await response.json()
    } catch (error) {
      console.error('API Get User Info error:', error)
      return {
        success: false,
        error: `Connection error: ${error.message}`
      }
    }
  },

  /**
   * Выход из системы
   */
  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    return { success: true, message: 'Logged out successfully' }
  },

  // ============= ЗАДАЧИ И КУРСЫ =============
  
  /**
   * Получение списка задач по языку
   */
  // В вашем api.js исправьте метод getTasks:
    async getTasks(language) {
      try {
        const response = await fetch(`/api/tasks?language=${language}`)
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`)
        }
        const data = await response.json()
        
        // ДОБАВЬТЕ ОТЛАДКУ
        console.log('API response for tasks:', data)
        if (Array.isArray(data) && data.length > 0) {
          console.log('Первая задача:', data[0])
          console.log('Тесты в первой задаче:', data[0].tests)
          
          // Проверяем expected_output
          if (data[0].tests && Array.isArray(data[0].tests)) {
            data[0].tests.forEach((test, i) => {
              console.log(`Тест ${i}:`, test)
              console.log(`Есть expected_output?`, test.expected_output)
            })
          }
        }
        
        return data
      } catch (error) {
        console.error('Error fetching tasks:', error)
        return []
      }
    },

  /**
   * Получение информации о задаче
   */
  async getTask(lang, topic, id) {
    try {
      const response = await fetch(`${API_BASE}/task/${lang}/${topic}/${id}`)
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      return await response.json()
    } catch (error) {
      console.error('API Task error:', error)
      return {
        error: `Connection error: ${error.message}`
      }
    }
  },

  // ============= ВЫПОЛНЕНИЕ И ПРОВЕРКА КОДА =============
  
  /**
   * Выполнение кода на бэкенде
   */
  async executeCode(requestData) {
    try {
      const response = await fetch(`${API_BASE}/execute`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestData)
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      return await response.json()
    } catch (error) {
      console.error('API Execute error:', error)
      return {
        success: false,
        message: `Connection error: ${error.message}`,
        output: ''
      }
    }
  },

  /**
   * Проверка решения задачи
   */
  async checkSolution(requestData) {
    try {
      const headers = {
        'Content-Type': 'application/json',
      }
      
      // Добавляем токен авторизации, если он есть
      const token = localStorage.getItem('token')
      if (token) {
        headers['Authorization'] = `Bearer ${token}`
      }
      
      const response = await fetch(`${API_BASE}/check`, {
        method: 'POST',
        headers: headers,
        body: JSON.stringify(requestData)
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      return await response.json()
    } catch (error) {
      console.error('API Check error:', error)
      return {
        success: false,
        passed: false,
        message: `Connection error: ${error.message}`,
        output: ''
      }
    }
  },

  // ============= AI АНАЛИЗ =============
  
  async analyzeCode(data) {
    try {
      console.log('[API] Запуск AI-анализа...')
      console.log('Длина кода:', data.code?.length)
      console.log('Язык:', data.language)
      console.log('Контекст задачи:', data.task_context)
      
      const headers = {
        'Content-Type': 'application/json',
      }
      
      // Добавляем токен для AI-анализа
      const token = localStorage.getItem('token')
      if (token) {
        headers['Authorization'] = `Bearer ${token}`
      }
      
      console.log('Отправка запроса к /api/ai/review')
      const response = await fetch(`${API_BASE}/ai/review`, {
        method: 'POST',
        headers: headers,
        body: JSON.stringify({
          code: data.code,
          language: data.language,
          task_context: data.task_context
        })
      })
      
      console.log('Статус ответа:', response.status, response.statusText)
      
      if (!response.ok) {
        const errorText = await response.text()
        console.error('Ошибка AI-анализа:', errorText)
        throw new Error(`AI анализ не удался: ${response.status} - ${errorText}`)
      }
      
      const result = await response.json()
      console.log('AI-анализ успешно выполнен:', result)
      return result
      
    } catch (error) {
      console.error('Ошибка API AI-анализа:', error)
      throw error
    }
  },

  // ============= ПАНЕЛЬ УЧИТЕЛЯ =============
  
  /**
   * Получение задач учителя
   */
  async getTeacherTasks() {
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        throw new Error('Authentication required')
      }
      
      const response = await fetch(`${API_BASE}/teacher/tasks`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      return await response.json()
    } catch (error) {
      console.error('API Get Teacher Tasks error:', error)
      return []
    }
  },

  /**
   * Создание новой задачи
   */
  async createTask(taskData) {
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        throw new Error('Authentication required')
      }
      
      const response = await fetch(`${API_BASE}/teacher/tasks`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(taskData)
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      return await response.json()
    } catch (error) {
      console.error('API Create Task error:', error)
      throw error
    }
  },

  /**
   * Обновление задачи
   */
  async updateTask(taskId, taskData) {
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        throw new Error('Authentication required')
      }
      
      const response = await fetch(`${API_BASE}/teacher/tasks/${taskId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(taskData)
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      return await response.json()
    } catch (error) {
      console.error('API Update Task error:', error)
      throw error
    }
  },

  /**
   * Удаление задачи
   */
  async deleteTask(taskId) {
    try {
      const token = localStorage.getItem('token')
      if (!token) {
        throw new Error('Authentication required')
      }
      
      const response = await fetch(`${API_BASE}/teacher/tasks/${taskId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      return await response.json()
    } catch (error) {
      console.error('API Delete Task error:', error)
      throw error
    }
  },

  // ============= СИСТЕМНЫЕ =============
  
  /**
   * Проверка здоровья бэкенда
   */
  async healthCheck() {
    try {
      const response = await fetch(`${API_BASE}/health`)
      return await response.json()
    } catch (error) {
      return { status: 'unhealthy', error: error.message }
    }
  },

  // ============= УТИЛИТЫ =============
  
  /**
   * Получение заголовков с авторизацией
   */
  getAuthHeaders() {
    const headers = {
      'Content-Type': 'application/json'
    }
    
    const token = localStorage.getItem('token')
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
    
    return headers
  },

  /**
   * Проверка авторизации
   */
  isAuthenticated() {
    const user = localStorage.getItem('user')
    if (!user) return false
    
    try {
      const userData = JSON.parse(user)
      return userData.isLoggedIn === true
    } catch (error) {
      return false
    }
  },

  /**
   * Получение данных пользователя
   */
  getUser() {
    const user = localStorage.getItem('user')
    if (!user) return null
    
    try {
      return JSON.parse(user)
    } catch (error) {
      return null
    }
  },

  /**
   * Получение роли пользователя
   */
  getUserRole() {
    const user = this.getUser()
    return user ? user.role : null
  },

  /**
   * Проверка, является ли пользователь учителем
   */
  isTeacher() {
    const role = this.getUserRole()
    return role === 'teacher' || role === 'admin'
  }
}