// frontend/src/utils/api.js
const API_BASE = '/api'

export const api = {
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

  /**
   * Получение списка задач
   */
  async getTasks(language) {
    try {
      const response = await fetch(`${API_BASE}/tasks?language=${language}`)
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }
      
      return await response.json()
    } catch (error) {
      console.error('API Tasks error:', error)
      return []
    }
  },

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
  }
}