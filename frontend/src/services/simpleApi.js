import axios from 'axios'

const API_BASE = 'http://localhost:8080/api'

// Простая версия API для тестирования
export const simpleAPI = {
  async executeCode(language, code) {
    try {
      console.log(`🚀 Sending request to ${API_BASE}/execute`)
      
      const response = await axios.post(`${API_BASE}/execute`, {
        language,
        code
      }, {
        timeout: 10000 // Уменьшаем таймаут
      })
      
      console.log('✅ API Response:', response.data)
      return response.data
      
    } catch (error) {
      console.error('❌ API Request failed:', error.message)
      
      // ФИКС: Не симулируем успех, а возвращаем ошибку
      throw new Error(`Бэкенд недоступен: ${error.message}`)
    }
  },

  async checkSolution(language, code, taskId) {
    try {
      console.log(`🚀 Sending request to ${API_BASE}/check`)
      
      const response = await axios.post(`${API_BASE}/check`, {
        language,
        code,
        task_id: taskId
      }, {
        timeout: 10000
      })
      
      console.log('✅ API Check Response:', response.data)
      return response.data
      
    } catch (error) {
      console.error('❌ API Check request failed:', error.message)
      
      // ФИКС: Пробрасываем ошибку вместо симуляции
      throw new Error(`Проверка недоступна: ${error.message}`)
    }
  }
}