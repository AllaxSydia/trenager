import axios from 'axios'

// Определяем базовый URL в зависимости от среды
const isDevelopment = import.meta.env.DEV
const API_BASE_URL = isDevelopment 
  ? 'http://localhost:10000/api'  // разработка - порт 10000
  : '/api'                        // продакшен - относительный путь

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
})

// Добавляем обработчик ошибок
api.interceptors.response.use(
  response => response,
  error => {
    if (error.code === 'ECONNREFUSED' || !error.response) {
      throw new Error('Не удалось подключиться к серверу. Проверьте, запущен ли бэкенд.')
    }
    throw error
  }
)

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