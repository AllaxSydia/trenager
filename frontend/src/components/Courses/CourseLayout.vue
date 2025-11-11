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
    console.log('Initial lessons (raw):', this.lessons)
    console.log('First lesson (raw):', this.lessons[0])
    console.log('First lesson tests (raw):', this.lessons[0]?.tests)
    
    // Проверим структуру тестов
    if (this.lessons[0]?.tests) {
      console.log('First test structure:', this.lessons[0].tests[0])
      console.log('First test keys:', Object.keys(this.lessons[0].tests[0]))
    }
    
    this.checkMobile()
    window.addEventListener('resize', this.checkMobile)
    this.checkBackendConnection()
    
    this.internalLessons = JSON.parse(JSON.stringify(this.lessons))
    this.ensureLessonSelected()
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
      if (health.status === 'healthy' || health.status === 'api_healthy') {
        this.consoleOutput += '✅ Все системы работают нормально\n'
      }
    } catch (error) {
      // Не показываем ошибку пользователю при загрузке
      console.log('Бэкенд недоступен:', error.message)
    }
  },

    async loadTasksFromAPI() {
      this.isLoading = true
      try {
        console.log(`Loading tasks for language: ${this.language}`)
        const tasks = await api.getTasks(this.language)
        console.log('Tasks from API:', tasks)
        
        // ДЕБАГ: посмотрим структуру полученных задач
        if (tasks && tasks.length > 0) {
          console.log('First task from API:', tasks[0])
          console.log('Tests in first task:', tasks[0].tests)
        }
        
        this.internalLessons = JSON.parse(JSON.stringify(this.lessons))
        console.log('Using lessons from props:', this.internalLessons)
        
        // ДЕБАГ: посмотрим тесты в пропсах
        if (this.internalLessons && this.internalLessons.length > 0) {
          console.log('First lesson tests from props:', this.internalLessons[0].tests)
        }
        
        this.ensureLessonSelected()
        this.updateProgress()
        
      } catch (error) {
        console.error('Failed to load tasks from API:', error)
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
      
      // ИСПРАВЛЕННЫЕ тесты с правильными входными данными
      const hardcodedTests = {
        'python_1': [{ 
          input: '', 
          expected_output: 'Hello, World!' 
        }],
        'python_2': [
          { 
            input: '5\n3', 
            expected_output: '8' 
          },
          { 
            input: '10\n20', 
            expected_output: '30' 
          },
          { 
            input: '-5\n8', 
            expected_output: '3' 
          }
        ],
        'python_3': [
          { 
            input: '5', 
            expected_output: '120' 
          },
          { 
            input: '3', 
            expected_output: '6' 
          },
          { 
            input: '1', 
            expected_output: '1' 
          }
        ],
        'python_4': [
          { 
            input: '4', 
            expected_output: 'чётное' 
          },
          { 
            input: '7', 
            expected_output: 'нечётное' 
          }
        ],
        'python_5': [
          { 
            input: '1\n2\n3', 
            expected_output: '3' 
          },
          { 
            input: '10\n5\n8', 
            expected_output: '10' 
          }
        ],
        'javascript_1': [{ 
          input: '', 
          expected_output: 'Hello, World!' 
        }],
        'javascript_2': [{ 
          input: '', 
          expected_output: '8' 
        }],
      }
      
      const testKey = `${this.language}_${lesson.id}`
      const tests = hardcodedTests[testKey] || []
      
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
          // РАЗБИВАЕМ input на массив строк если есть \n
          let inputs = []
          if (test.input && test.input.trim() !== '') {
            inputs = test.input.split('\n').filter(line => line.trim() !== '')
          }
          
          const result = await api.executeCode({
            code: this.userCode,
            language: this.language,
            inputs: inputs // ← передаем массив строк
          })
          
          const output = result.output || ''
          const expected = test.expected_output || ''
          
          const testPassed = output.trim() === expected.trim()
          
          if (testPassed) {
            test.status = 'passed'
            passedCount++
            this.consoleOutput += `✅ Тест ${i + 1}: Пройден\n`
            if (output) {
              this.consoleOutput += `   Вывод: "${output}"\n`
            }
            this.consoleOutput += '\n'
          } else {
            test.status = 'failed'
            this.consoleOutput += `❌ Тест ${i + 1}: Не пройден\n`
            this.consoleOutput += `   Ожидалось: "${expected}"\n`
            this.consoleOutput += `   Получено:  "${output}"\n\n`
          }
          
        } catch (error) {
          test.status = 'failed'
          this.consoleOutput += `❌ Тест ${i + 1}: Ошибка выполнения\n`
          this.consoleOutput += `   Ошибка: ${error.message}\n\n`
        }
        
        await new Promise(resolve => setTimeout(resolve, 500))
      }
      
      this.consoleOutput += `📊 Итог: ${passedCount}/${this.currentLesson.tests.length} тестов пройдено\n`
      
      if (passedCount === this.currentLesson.tests.length) {
        this.consoleOutput += '🎉 Отлично! Все тесты пройдены!\n'
      }
      
      this.isTesting = false
    },

    async submitSolution() {
      if (!this.currentLesson) return
      
      this.isSubmitting = true
      this.consoleOutput = '📤 Проверка решения...\n\n'
      
      await this.runTests()
      
      const allPassed = this.currentLesson.tests.every(test => test.status === 'passed')
      
      if (allPassed) {
        // Отмечаем урок как пройденный
        const lessonIndex = this.internalLessons.findIndex(l => l.id === this.currentLesson.id)
        if (lessonIndex !== -1) {
          this.internalLessons[lessonIndex].completed = true
          this.updateProgress()
        }
        this.consoleOutput += '\n🎉 Поздравляем! Все тесты пройдены! Задача решена правильно.\n'
      } else {
        this.consoleOutput += '\n❌ Не все тесты пройдены. Продолжайте работать над решением!\n'
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