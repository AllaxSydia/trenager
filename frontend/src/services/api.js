import axios from 'axios'

// Для разных сред
const getApiBaseUrl = () => {
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL  // из переменной окружения
  }
  return import.meta.env.DEV 
    ? 'http://localhost:8080/api'
    : '/api'
}

const API_BASE_URL = getApiBaseUrl()

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
})

export const taskService = {
  async getTasks() {
    const response = await api.get('/tasks')
    return response.data
  },
  
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