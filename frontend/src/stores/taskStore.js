import { TASK_DATA } from '@/utils/constants'
import { defineStore } from 'pinia'

export const useTaskStore = defineStore('task', {
  state: () => ({
    currentLanguage: null,
    currentTask: null,
    userCode: ''
  }),

  actions: {
    setLanguage(lang) {
      this.currentLanguage = lang
    },
    
    setTask(taskId) {
      this.currentTask = taskId
      const tasks = TASK_DATA[this.currentLanguage] || []
      const task = tasks.find(t => t.id === taskId)
      if (task) {
        this.userCode = task.defaultCode
      }
    }
  }
})