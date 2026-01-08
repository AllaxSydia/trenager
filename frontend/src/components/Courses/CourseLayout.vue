<template>
  <div class="course-page">
    <!-- Мобильный хедер -->
    <MobileHeader
      :title="courseTitle"
      :progress="progress"
      :is-mobile="isMobile"
      @toggle-sidebar="showSidebar = true"
    />

    <div class="course-layout">
      <!-- Боковая панель с заданиями -->
      <LessonSidebar
        :lessons="internalLessons"
        :current-lesson="currentLesson"
        :progress="progress"
        :title="courseTitle"
        :is-mobile="isMobile"
        @select-lesson="selectLesson"
      />

      <!-- Основной контент -->
      <main class="main-content" v-if="currentLesson">

        <ProblemSection :lesson="currentLesson" />
        
        <CodeSection
          v-model:code="userCode"
          :language="language"
          @reset="resetCode"
          @execute="executeCode"
        />
        
        <ControlPanel
          :is-running="isRunning"
          :is-testing="isTesting"
          :is-submitting="isSubmitting"
          :ai-loading="aiLoading"
          @run="runCode"
          @test="runTests"
          @submit="submitSolution"
          @save="saveCode"
          @analyze="analyzeWithAI"
        />
        
        <InputSection 
          v-model:input="consoleInput"
          @execute="executeCode"
        />
        
        <OutputSection
          :output="consoleOutput"
          @clear="clearOutput"
        />
        
        <TestsSection
          :tests="currentLesson.tests || []"
          :passed-tests="passedTests"
        />
      </main>

      <!-- Сообщение если урок не выбран -->
      <div v-else class="no-lesson-selected">
        <div class="loading-message">
          <h3>Загрузка курса {{ language }}...</h3>
          <p v-if="isLoading">Загружаем задачи с сервера...</p>
          <p v-else-if="apiTasks.length > 0">Выберите задачу из списка слева</p>
          <p v-else>Нет доступных задач для этого языка</p>
          <div class="debug-info">
            <p><strong>Отладка:</strong></p>
            <p>Задач из БД: {{ apiTasks.length }}</p>
            <p>Уроков доступно: {{ internalLessons.length }}</p>
            <p>Текущий урок: {{ currentLesson ? currentLesson.title : 'не выбран' }}</p>
            <p>Язык: {{ language }}</p>
          </div>
        </div>
      </div>
    </div>

    

    <!-- Мобильный сайдбар -->
    <MobileSidebar
        v-if="isMobile && showSidebar"
        :lessons="internalLessons"
        :current-lesson="currentLesson"
        :title="courseTitle"
        :is-mobile="isMobile"
        :show-sidebar="showSidebar"
        @select-lesson="selectLessonMobile"
        @close="showSidebar = false"
        />
  </div>
</template>

<script>
import { api } from '@/utils/api.js'
import CodeSection from './CodeSection.vue'
import ControlPanel from './ControlPanel.vue'
import InputSection from './InputSection.vue'
import LessonSidebar from './LessonSidebar.vue'
import MobileHeader from './MobileHeader.vue'
import MobileSidebar from './MobileSidebar.vue'
import OutputSection from './OutputSection.vue'
import ProblemSection from './ProblemSection.vue'
import TestsSection from './TestsSection.vue'

export default {
  name: 'CourseLayout',
  components: {
    MobileHeader,
    LessonSidebar,
    MobileSidebar,
    ProblemSection,
    CodeSection,
    ControlPanel,
    InputSection,
    OutputSection,
    TestsSection
  },
  props: {
    courseTitle: {
      type: String,
      required: true
    },
    lessons: {
      type: Array,
      default: () => []
    },
    language: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      userCode: '',
      consoleInput: '',
      consoleOutput: '',
      isRunning: false,
      isTesting: false,
      isSubmitting: false,
      progress: 0,
      currentLesson: null,
      showSidebar: false,
      isMobile: false,
      isLoading: false,
      internalLessons: [],
      aiLoading: false,
      apiTasks: [],  // Задачи из API/БД
      aiResult: null
    }
  },
  computed: {
    passedTests() {
      if (!this.currentLesson?.tests) return 0
      return this.currentLesson.tests.filter(test => test.status === 'passed').length
    }
  },
  watch: {
    language: {
      immediate: true,
      handler(newLanguage) {
        console.log(`Language changed to: ${newLanguage}`)
        if (newLanguage) {
          this.loadTasksFromAPI()
        }
      }
    },
    lessons: {
      immediate: true,
      handler(newLessons) {
        console.log('Lessons prop updated:', newLessons)
        // Если есть переданные уроки, используем их
        if (newLessons && newLessons.length > 0) {
          this.internalLessons = this.formatLessons(newLessons)
        } else {
          // Иначе используем задачи из API
          this.internalLessons = this.formatLessons(this.apiTasks)
        }
        this.ensureLessonSelected()
      }
    },
    apiTasks: {
      handler(newTasks) {
        console.log('API tasks updated:', newTasks)
        // Обновляем уроки когда загружаются задачи из API
        this.internalLessons = this.formatLessons(newTasks)
        this.ensureLessonSelected()
      }
    }
  },
  mounted() {
    console.log('CourseLayout mounted for language:', this.language)
    
    this.checkMobile()
    window.addEventListener('resize', this.checkMobile)
    this.checkBackendConnection()
    
    // Инициализируем уроки
    if (this.lessons && this.lessons.length > 0) {
      this.internalLessons = this.formatLessons(this.lessons)
    } else {
      this.loadTasksFromAPI()
    }
    
    this.ensureLessonSelected()
  },
  methods: {
    // Форматируем задачи в формат уроков
    formatLessons(tasks) {
      if (!tasks || !Array.isArray(tasks)) return []
      
      console.log('=== ФОРМАТИРОВАНИЕ УРОКОВ ===')
      console.log('Входные задачи:', tasks.length)
      
      const result = tasks.map((task, index) => {
        const formattedLesson = {
          id: task.id || `task_${index + 1}`,
          title: task.title || `Задача ${index + 1}`,
          description: task.description || '',
          starterCode: task.template || task.starter_code || '',
          code: task.template || task.starter_code || '',
          language: task.language || this.language,
          difficulty: task.difficulty || 'beginner',
          completed: false,
          tests: this.prepareTests(task.tests || []),
          apiData: task
        }
        
        console.log(`Урок ${index}:`, formattedLesson.title)
        console.log('Кол-во тестов:', formattedLesson.tests.length)
        if (formattedLesson.tests.length > 0) {
          console.log('Тесты:', formattedLesson.tests)
        }
        
        return formattedLesson
      })
      
      return result
    },
    
    async loadTasksFromAPI() {
      this.isLoading = true
      try {
        console.log(`Загрузка задач для языка: ${this.language}`)
        const tasks = await api.getTasks(this.language)
        console.log('Задачи получены:', tasks)
        
        if (tasks && tasks.length > 0) {
          // Форматируем задачи в уроки
          this.internalLessons = this.formatLessons(tasks)
          console.log(`Успешно загружено ${tasks.length} задач`)
          this.ensureLessonSelected()
        } else {
          // Используем статические уроки
          this.useStaticLessons()
        }
        
      } catch (error) {
        console.error('Ошибка загрузки задач:', error)
        this.useStaticLessons()
      } finally {
        this.isLoading = false
      }
    },
    
    formatLessons(tasks) {
      if (!tasks || !Array.isArray(tasks)) return []
      
      console.log('Форматируем задачи в уроки...')
      
      return tasks.map((task, index) => ({
        id: task.id || index + 1,
        title: task.title || `Задача ${index + 1}`,
        description: task.description || '',
        starterCode: task.template || task.starter_code || task.code_template || '',
        code: task.template || task.starter_code || task.code_template || '',
        language: task.language || this.language,
        difficulty: task.difficulty || 'beginner',
        completed: false,
        tests: this.prepareTests(task.tests || [])
      }))
    },
    
    prepareTests(tests) {
      console.log('Подготавливаем тесты:', tests)
      
      if (!Array.isArray(tests)) return []
      
      return tests.map(test => ({
        input: test.input || '',
        expected_output: test.expected_output || '',  // Оставляем как есть, даже если пустое
        status: null,
        actual: null,
        error: null
      }))
    },
    
    useStaticLessons() {
      console.log('Используем статические уроки')
      
      const staticLessons = {
        python: [
          {
            id: 1,
            title: "Проверка числа на четность",
            description: "Напишите программу, которая проверяет, является ли число четным",
            language: "python",
            difficulty: "beginner",
            starterCode: `num = int(input())\nif num % 2 == 0:\n    print("Четное")\nelse:\n    print("Нечетное")`,
            code: `num = int(input())\nif num % 2 == 0:\n    print("Четное")\nelse:\n    print("Нечетное")`,
            tests: [
              { input: "4", expected_output: "Четное" },
              { input: "7", expected_output: "Нечетное" }
            ]
          },
          {
            id: 2,
            title: "Сумма двух чисел",
            description: "Напишите программу, которая принимает два числа через input() и выводит их сумму",
            language: "python",
            difficulty: "beginner",
            starterCode: `num1 = int(input())\nnum2 = int(input())\nprint(num1 + num2)`,
            code: `num1 = int(input())\nnum2 = int(input())\nprint(num1 + num2)`,
            tests: [
              { input: "5\n3", expected_output: "8" },
              { input: "10\n20", expected_output: "30" }
            ]
          }
        ],
        javascript: [
          {
            id: 1,
            title: "Hello World на JavaScript",
            description: "Напишите программу, которая выводит 'Hello, World!'",
            language: "javascript",
            difficulty: "beginner",
            starterCode: `console.log("Hello, World!")`,
            code: `console.log("Hello, World!")`,
            tests: [
              { input: "", expected_output: "Hello, World!" }
            ]
          }
        ]
      }
      
      this.internalLessons = staticLessons[this.language] || []
      this.ensureLessonSelected()
    },
    
    ensureLessonSelected() {
      if (this.internalLessons.length > 0 && !this.currentLesson) {
        console.log('Выбираем первый урок:', this.internalLessons[0])
        this.selectLesson(this.internalLessons[0])
      } else if (this.internalLessons.length === 0) {
        console.warn('Нет уроков для выбора')
        this.currentLesson = null
      }
    },

    async checkBackendConnection() {
      try {
        const health = await api.healthCheck()
        if (health.status === 'healthy' || health.status === 'api_healthy') {
          this.consoleOutput += 'Все системы работают нормально\n'
        }
      } catch (error) {
        console.log('Бэкенд недоступен:', error.message)
      }
    },

    selectLesson(lesson) {
      if (!lesson) {
        console.error('Attempted to select null lesson')
        return
      }
      
      console.log('Selecting lesson:', lesson.title)
      console.log('Lesson data:', lesson)
      
      // Используем тесты из задачи
      const tests = lesson.tests || []
      
      this.currentLesson = { 
        ...lesson,
        tests: tests.map(test => ({ 
          ...test,
          status: null,
          actual: null,
          error: null
        }))
      }
      
      console.log('Current lesson with tests:', this.currentLesson)
      
      this.userCode = lesson.starterCode || lesson.code || ''
      this.consoleInput = ''
      this.consoleOutput = ''
      this.loadSavedCode()
    },

    selectLessonMobile(lesson) {
      this.selectLesson(lesson)
      this.showSidebar = false
    },

    resetCode() {
      this.userCode = this.currentLesson?.starterCode || this.currentLesson?.code || ''
      this.consoleOutput = 'Код сброшен к начальному состоянию\n'
    },

    executeCode() {
      this.runCode()
    },

    async runCode() {
      if (!this.userCode?.trim()) {
        this.consoleOutput = 'Введите код для выполнения\n'
        return
      }

      this.isRunning = true
      this.consoleOutput = 'Выполнение кода...\n\n'

      try {
        const inputs = this.consoleInput.trim() ? [this.consoleInput] : []
        
        const result = await api.executeCode({
          code: this.userCode,
          language: this.language,
          inputs: inputs
        })
        
        if (result.success) {
          this.consoleOutput += `Успешно!\n${result.output || 'Программа выполнена без вывода'}\n`
        } else {
          this.consoleOutput += `Ошибка выполнения:\n${result.output || result.message}\n`
        }
      } catch (error) {
        this.consoleOutput += `Ошибка соединения: ${error.message}\n`
      } finally {
        this.isRunning = false
      }
    },

    async analyzeWithAI() {
      if (!this.userCode?.trim()) {
        console.log('Нет кода для анализа')
        this.consoleOutput += '\nВведите код для AI анализа\n'
        return
      }
      
      if (this.aiLoading) {
        console.log('AI анализ уже выполняется')
        return
      }
      
      console.log('Запуск AI анализа...')
      
      this.aiLoading = true
      this.aiResult = null  // Сбрасываем предыдущий результат
      this.consoleOutput += '\nЗапуск AI анализа кода...\n'

      try {
        const aiResult = await api.analyzeCode({
          code: this.userCode,
          language: this.language,
          task_context: this.currentLesson?.description || 'Анализ кода студента'
        })
        
        console.log('AI анализ завершен, результат:', aiResult)
        
        if (aiResult && aiResult.score !== undefined) {
          this.aiResult = aiResult  // Сохраняем результат для UI
          this.formatAIResponse(aiResult)  // Также выводим в консоль
        } else {
          console.error('Неверный формат ответа от AI:', aiResult)
          this.consoleOutput += '\nОшибка: неверный формат ответа от AI\n'
        }
        
      } catch (error) {
        console.error('Ошибка AI-анализа:', error)
        this.consoleOutput += `\nОшибка AI анализа: ${error.message}\n`
      } finally {
        this.aiLoading = false
      }
    },

    formatAIResponse(aiData) {
      this.consoleOutput += '='.repeat(50) + '\n'
      this.consoleOutput += 'AI АНАЛИЗ КОДА:\n'
      this.consoleOutput += '='.repeat(50) + '\n\n'
      
      // Оценка
      const score = aiData.score || 0
      this.consoleOutput += `ОЦЕНКА: ${score}/10\n`
      this.consoleOutput += `СЛОЖНОСТЬ: ${aiData.complexity || 'неизвестно'}\n\n`
      
      // Комментарии
      if (aiData.comments && aiData.comments.length > 0) {
        this.consoleOutput += 'КОММЕНТАРИИ:\n'
        aiData.comments.forEach((comment, index) => {
          this.consoleOutput += `  ${index + 1}. ${comment}\n`
        })
        this.consoleOutput += '\n'
      }
      
      // Предложения
      if (aiData.suggestions && aiData.suggestions.length > 0) {
        this.consoleOutput += '💡 ПРЕДЛОЖЕНИЯ ПО УЛУЧШЕНИЮ:\n'
        aiData.suggestions.forEach((suggestion, index) => {
          this.consoleOutput += `  ${index + 1}. ${suggestion}\n`
        })
        this.consoleOutput += '\n'
      }
      
      // Best Practices
      if (aiData.best_practices && aiData.best_practices.length > 0) {
        this.consoleOutput += 'РЕКОМЕНДУЕМЫЕ ПРАКТИКИ:\n'
        aiData.best_practices.forEach((practice, index) => {
          this.consoleOutput += `  ${index + 1}. ${practice}\n`
        })
        this.consoleOutput += '\n'
      }
      
      // Альтернативные решения
      if (aiData.alternative_solutions && aiData.alternative_solutions.length > 0) {
        this.consoleOutput += 'АЛЬТЕРНАТИВНЫЕ РЕШЕНИЯ:\n'
        aiData.alternative_solutions.forEach((solution, index) => {
          this.consoleOutput += `  ${index + 1}. ${solution}\n`
        })
        this.consoleOutput += '\n'
      }
      
      this.consoleOutput += '='.repeat(50) + '\n'
      this.consoleOutput += 'AI анализ завершен!\n'
    },

    async runTests() {
      if (!this.currentLesson?.tests || this.currentLesson.tests.length === 0) {
        this.consoleOutput = 'Для этой задачи нет тестов\n'
        return
      }
      
      this.isTesting = true
      this.consoleOutput = 'Запуск тестов...\n\n'

      let passedCount = 0

      for (let i = 0; i < this.currentLesson.tests.length; i++) {
        const test = this.currentLesson.tests[i]
        test.status = 'running'
        
        try {
          const result = await api.executeCode({
            code: this.userCode,
            language: this.language,
            inputs: test.input ? test.input.split('\n') : []
          })
          
          const output = result.output || ''
          const expected = test.expected_output || ''
          
          // Сохраняем фактический вывод для отображения в тестах
          test.actual = output.trim()
          
          const testPassed = test.actual === expected.trim()
          
          if (testPassed) {
            test.status = 'passed'
            passedCount++
          this.consoleOutput += `Тест ${i + 1}: Пройден\n`
          } else {
            test.status = 'failed'
          this.consoleOutput += `Тест ${i + 1}: Не пройден\n`
          }
          
        } catch (error) {
          test.status = 'failed'
          test.error = error.message
          test.actual = '' // Устанавливаем пустую строку при ошибке
        this.consoleOutput += `Тест ${i + 1}: Ошибка выполнения\n`
        }
        
        await new Promise(resolve => setTimeout(resolve, 500))
      }
      
      this.consoleOutput += `Итог: ${passedCount}/${this.currentLesson.tests.length} тестов пройдено\n`
      
      if (passedCount === this.currentLesson.tests.length) {
        this.consoleOutput += 'Отлично! Все тесты пройдены!\n'
      }
      
      this.isTesting = false
    },

    async submitSolution() {
      if (!this.currentLesson) return
      
      this.isSubmitting = true
      this.consoleOutput = 'Проверка решения...\n\n'
      
      await this.runTests()
      
      const allPassed = this.currentLesson.tests.every(test => test.status === 'passed')
      
      if (allPassed) {
        // Отмечаем урок как пройденный
        const lessonIndex = this.internalLessons.findIndex(l => l.id === this.currentLesson.id)
        if (lessonIndex !== -1) {
          this.internalLessons[lessonIndex].completed = true
          this.updateProgress()
        }
        this.consoleOutput += '\nПоздравляем! Все тесты пройдены! Задача решена правильно.\n'
        
        // Автоматически запускаем AI анализ при успешной сдаче
        setTimeout(() => {
          this.analyzeWithAI()
        }, 1000)
      } else {
        this.consoleOutput += '\nНе все тесты пройдены. Продолжайте работать над решением!\n'
      }
      
      this.isSubmitting = false
    },

    updateProgress() {
      const completedCount = this.internalLessons.filter(lesson => lesson.completed).length
      this.progress = Math.round((completedCount / this.internalLessons.length) * 100)
    },

    saveCode() {
      if (!this.currentLesson) return
      localStorage.setItem(`${this.language}_lesson_${this.currentLesson.id}`, this.userCode)
      this.consoleOutput = 'Код сохранен локально.\n'
    },

    loadSavedCode() {
      if (!this.currentLesson) return
      const savedCode = localStorage.getItem(`${this.language}_lesson_${this.currentLesson.id}`)
      if (savedCode) {
        this.userCode = savedCode
      }
    },

    clearOutput() {
      this.consoleOutput = ''
    },
    
    checkMobile() {
      this.isMobile = window.innerWidth <= 1024
    }
  }
}
</script>

<style scoped>
.course-page {
  background-color: #0E1117;
  color: #E2E8F0;
  padding: 20px;
}

.course-layout {
  display: grid;
  grid-template-columns: 340px 1fr;
  gap: 20px;
  align-items: start;
}

/* Кастомные скроллбары */
::-webkit-scrollbar {
  width: 12px;
  height: 12px;
}

::-webkit-scrollbar-track {
  background: #1E1E1E;
  border-radius: 6px;
}

::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #3B82F6, #8B5CF6);
  border-radius: 6px;
  border: 2px solid #1E1E1E;
  transition: all 0.3s ease;
}

::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, #2563EB, #7C3AED);
  transform: scale(1.05);
}

::-webkit-scrollbar-thumb:active {
  background: linear-gradient(135deg, #1D4ED8, #6D28D9);
}

* {
  scrollbar-width: thin;
  scrollbar-color: #3B82F6 #1E1E1E;
}

.main-content {
  padding: 0;
  overflow-y: auto;
  max-height: calc(100vh - 140px);
  margin-top: 0;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.main-content::-webkit-scrollbar {
  width: 10px;
}

.main-content::-webkit-scrollbar-track {
  background: #0E1117;
  border-radius: 6px;
}

.main-content::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #3B82F6, #8B5CF6);
  border-radius: 6px;
  border: 2px solid #0E1117;
}

@media (max-width: 1024px) {
  .course-layout {
    grid-template-columns: 1fr;
  }
  
  .main-content {
    overflow-y: visible;
    max-height: none;
  }
}

@media (max-width: 768px) {
  .course-page {
    padding: 10px;
  }
}

@media (max-width: 480px) {
  .course-page {
    padding: 5px;
  }
}

.no-lesson-selected {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 400px;
  grid-column: 2;
}

.loading-message {
  text-align: center;
  color: #94A3B8;
}

.loading-message h3 {
  margin-bottom: 10px;
  color: #E2E8F0;
}

.debug-info {
  margin-top: 20px;
  padding: 15px;
  background: #1E293B;
  border-radius: 8px;
  border: 1px solid #334155;
  font-family: monospace;
  font-size: 14px;
  text-align: left;
}

.debug-info p {
  margin: 5px 0;
}
</style>