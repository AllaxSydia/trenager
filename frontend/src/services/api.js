import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000, // Увеличили таймаут для выполнения кода
})

export const taskService = {
  // Получить все задачи
  async getTasks() {
    const response = await api.get('/tasks')
    return response.data
  },
  
  // Выполнить код
  async executeCode(taskId, code, language) {
    const response = await api.post('/execute', {
      task_id: taskId,
      code,
      language
    })
    return response.data
  }
}

export default api