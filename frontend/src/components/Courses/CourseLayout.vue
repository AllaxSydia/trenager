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
          @run="runCode"
          @test="runTests"
          @submit="submitSolution"
          @save="saveCode"
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
          <p>Выберите урок из списка слева</p>
          <div class="debug-info">
            <p><strong>Отладка:</strong></p>
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
      required: true
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
      internalLessons: []
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
        this.internalLessons = JSON.parse(JSON.stringify(newLessons))
        this.ensureLessonSelected()
      }
    }
  },
  mounted() {
    console.log('CourseLayout mounted for language:', this.language)
    console.log('Initial lessons:', this.lessons)
    
    this.checkMobile()
    window.addEventListener('resize', this.checkMobile)
    this.checkBackendConnection()
    
    // Инициализируем internalLessons из пропсов
    this.internalLessons = JSON.parse(JSON.stringify(this.lessons))
    this.ensureLessonSelected()
  },
  beforeUnmount() {
    window.removeEventListener('resize', this.checkMobile)
  },
  methods: {
    ensureLessonSelected() {
      if (this.internalLessons.length > 0 && !this.currentLesson) {
        console.log('Selecting first lesson:', this.internalLessons[0])
        this.selectLesson(this.internalLessons[0])
      } else if (this.internalLessons.length === 0) {
        console.warn('No lessons available for selection')
      }
    },

    async checkBackendConnection() {
      try {
        const health = await api.healthCheck()
        if (health.status !== 'healthy') {
          this.consoleOutput = '⚠️ Бэкенд недоступен. Убедитесь, что сервер запущен на localhost:8080\n'
        }
      } catch (error) {
        this.consoleOutput = '⚠️ Не удалось подключиться к бэкенду\n'
      }
    },

    async loadTasksFromAPI() {
      this.isLoading = true
      try {
        console.log(`Loading tasks for language: ${this.language}`)
        const tasks = await api.getTasks(this.language)
        console.log('Tasks from API:', tasks)
        
        // ВСЕГДА используем задачи из пропсов, игнорируем задачи из API
        this.internalLessons = JSON.parse(JSON.stringify(this.lessons))
        console.log('Using lessons from props:', this.internalLessons)
        
        this.ensureLessonSelected()
        this.updateProgress()
        
      } catch (error) {
        console.error('Failed to load tasks from API:', error)
        console.log('Не удалось загрузить задачи из API, используем пропсы:', error.message)
        this.internalLessons = JSON.parse(JSON.stringify(this.lessons))
        this.ensureLessonSelected()
        this.updateProgress()
      } finally {
        this.isLoading = false
      }
    },

    updateProgress() {
      const completedCount = this.internalLessons.filter(lesson => lesson.completed).length
      this.progress = Math.round((completedCount / this.internalLessons.length) * 100)
    },

    checkMobile() {
      this.isMobile = window.innerWidth <= 1024
    },

    selectLesson(lesson) {
      if (!lesson) {
        console.error('Attempted to select null lesson')
        return
      }
      
      console.log('Selecting lesson:', lesson.title)
      
      if (this.currentLesson?.id === lesson.id) {
        console.log('Lesson already selected')
        return
      }
      
      this.currentLesson = { 
        ...lesson,
        tests: lesson.tests ? lesson.tests.map(test => ({ 
          ...test,
          status: null,
          actual: null,
          error: null
        })) : []
      }
      
      this.userCode = lesson.starterCode || ''
      this.consoleInput = ''
      this.consoleOutput = ''
      this.loadSavedCode()
    },

    selectLessonMobile(lesson) {
      this.selectLesson(lesson)
      this.showSidebar = false
    },

    resetCode() {
      this.userCode = this.currentLesson?.starterCode || ''
      this.consoleOutput = '🔄 Код сброшен к начальному состоянию\n'
    },

    executeCode() {
      this.runCode()
    },

    async runCode() {
      if (!this.userCode?.trim()) {
        this.consoleOutput = '❌ Введите код для выполнения\n'
        return
      }

      this.isRunning = true
      this.consoleOutput = '🚀 Выполнение кода...\n\n'

      try {
        const inputs = this.consoleInput.trim() ? [this.consoleInput] : []
        
        const result = await api.executeCode({
          code: this.userCode,
          language: this.language,
          inputs: inputs
        })
        
        if (result.success) {
          this.consoleOutput += `✅ Успешно!\n${result.output || 'Программа выполнена без вывода'}\n`
        } else {
          this.consoleOutput += `❌ Ошибка выполнения:\n${result.output || result.message}\n`
        }
      } catch (error) {
        this.consoleOutput += `❌ Ошибка соединения: ${error.message}\n`
      } finally {
        this.isRunning = false
      }
    },

    async runTests() {
      if (!this.currentLesson?.tests || this.currentLesson.tests.length === 0) {
        this.consoleOutput = 'ℹ️ Для этой задачи нет тестов\n'
        return
      }
      
      this.isTesting = true
      this.consoleOutput = '🧪 Запуск тестов...\n\n'

      let passedCount = 0

      for (let i = 0; i < this.currentLesson.tests.length; i++) {
        const test = this.currentLesson.tests[i]
        test.status = 'running'
        
        try {
          const inputs = this.consoleInput.trim() ? [this.consoleInput] : []
          
          const result = await api.checkSolution({
            task_id: this.currentLesson.id.toString(),
            code: this.userCode,
            language: this.language,
            inputs: inputs
          })
          
          if (result.success && result.passed) {
            test.status = 'passed'
            test.actual = result.actual || 'Тест пройден'
            passedCount++
          } else {
            test.status = 'failed'
            test.actual = result.actual || 'Тест не пройден'
            test.error = result.message || 'Ошибка выполнения'
          }
          
          this.consoleOutput += `Тест ${i + 1}: ${test.status === 'passed' ? '✅' : '❌'} ${test.input || 'без ввода'}\n`
          
        } catch (error) {
          test.status = 'failed'
          test.error = `Ошибка: ${error.message}`
          this.consoleOutput += `Тест ${i + 1}: ❌ Ошибка выполнения\n`
        }
        
        await new Promise(resolve => setTimeout(resolve, 300))
      }
      
      this.consoleOutput += `\n📊 Результат: ${passedCount}/${this.currentLesson.tests.length} тестов пройдено\n`
      this.isTesting = false
    },

    async submitSolution() {
      if (!this.currentLesson) return
      
      this.isSubmitting = true
      this.consoleOutput = '📤 Отправка решения...\n\n'
      
      await this.runTests()
      
      const allPassed = this.passedTests === this.currentLesson.tests.length
      
      if (allPassed) {
        const lessonIndex = this.internalLessons.findIndex(l => l.id === this.currentLesson.id)
        if (lessonIndex !== -1) {
          this.internalLessons[lessonIndex].completed = true
          this.updateProgress()
        }
        this.consoleOutput += '\n🎉 Поздравляем! Все тесты пройдены! Решение принято.\n'
      } else {
        this.consoleOutput += '\n❌ Не все тесты пройдены. Продолжайте работать!\n'
      }
      
      this.isSubmitting = false
    },

    saveCode() {
      if (!this.currentLesson) return
      localStorage.setItem(`${this.language}_lesson_${this.currentLesson.id}`, this.userCode)
      this.consoleOutput = '💾 Код сохранен локально.\n'
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
</style>