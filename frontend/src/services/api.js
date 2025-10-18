import axios from 'axios'

// Для разных сред
const getApiBaseUrl = () => {
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL  // из переменной окружения
  }
  return import.meta.env.DEV 
    ? 'http://localhost:8080/api'
    : 'https://trenager.onrender.com/api'  // ← АБСОЛЮТНЫЙ URL для продакшена
}

const API_BASE_URL = getApiBaseUrl()

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
})

// Добавляем логирование для отладки
api.interceptors.request.use(config => {
  console.log(`🔄 API Request: ${config.method?.toUpperCase()} ${config.baseURL}${config.url}`)
  return config
})

api.interceptors.response.use(
  response => {
    console.log(`✅ API Response: ${response.status} ${response.config.url}`)
    return response
  },
  error => {
    console.error(`❌ API Error: ${error.message}`)
    console.error(`📡 Failed URL: ${error.config?.baseURL}${error.config?.url}`)
    return Promise.reject(error)
  }
)

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