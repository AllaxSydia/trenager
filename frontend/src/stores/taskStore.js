import { taskService } from '@/services/api'
import { defineStore } from 'pinia'

export const useTaskStore = defineStore('tasks', {
  state: () => ({
    tasks: [],
    currentTask: null,
    userCode: '',
    executionResult: null,
    isLoading: false,
    language: 'python'
  }),
  
  actions: {
    async loadTasks() {
      try {
        this.tasks = await taskService.getTasks()
        
        if (this.tasks.length > 0 && !this.currentTask) {
          this.currentTask = this.tasks[0]
          this.userCode = this.tasks[0].template
        }
      } catch (error) {
        console.error('Ошибка загрузки задач:', error)
        this.tasks = [
          {
            id: '1',
            title: 'Hello World',
            description: 'Напишите программу которая выводит "Hello, World!"',
            template: 'print("Hello, World!")'
          }
        ]
      }
    },
    
    setCurrentTask(task) {
      this.currentTask = task
      this.userCode = task.template || ''
    },
    
    async executeCode() {
      if (!this.currentTask) return
      
      this.isLoading = true
      this.executionResult = null
      
      try {
        const result = await taskService.executeCode(
          this.currentTask.id,
          this.userCode,
          this.language
        )
        this.executionResult = result
      } catch (error) {
        this.executionResult = {
          success: false,
          message: 'Ошибка соединения с сервером: ' + error.message,
          output: 'Проверьте что бэкенд сервер запущен на localhost:8080'
        }
      } finally {
        this.isLoading = false
      }
    },
    
    updateCode(code) {
      this.userCode = code
    },
    
    // ДОБАВЛЯЕМ ЭТОТ МЕТОД
    setLanguage(lang) {
      this.language = lang
    }
  }
})