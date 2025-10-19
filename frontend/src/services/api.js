import axios from 'axios'

// УПРОЩЕННАЯ ФУНКЦИЯ - всегда используем относительные пути
const getApiBaseUrl = () => {
    // В продакшене и в Docker - всегда относительный путь
    // Бэкенд и фронтенд в одном контейнере, поэтому /api проксируется на тот же сервер
    return '/api'
}

const API_BASE_URL = getApiBaseUrl()

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
})

// Логирование
api.interceptors.request.use(config => {
  console.log(`🔄 API Request: ${config.method?.toUpperCase()} ${config.baseURL}${config.url}`)
  console.log(`🌐 Frontend URL: ${window.location.host}`)
  console.log(`🎯 Full URL: ${config.baseURL}${config.url}`)
  return config
})

// Обработка ошибок
api.interceptors.response.use(
  response => response,
  error => {
    console.error('🚨 API Error:', error.response?.data || error.message)
    console.error('🔧 Error config:', error.config)
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

// Test connection
export const testConnection = async () => {
  try {
    const response = await api.get('/test')
    console.log('✅ Connection test:', response.data)
    return response.data
  } catch (error) {
    console.error('❌ Connection test failed:', error)
    throw error
  }
}

export default api