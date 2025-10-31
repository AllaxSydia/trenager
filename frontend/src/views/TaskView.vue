<template>
  <div class="task-view">
    <div class="container">
      <header class="header">
        <button class="back-btn" @click="goBack">← Назад</button>
        <h1>{{ task?.title || 'Задача' }}</h1>
        <div class="language-selector">
          <select v-model="currentLanguage" @change="updateEditorLanguage">
            <option value="python">Python</option>
            <option value="javascript">JavaScript</option>
            <option value="cpp">C++</option>
            <option value="java">Java</option>
          </select>
        </div>
      </header>

      <div class="task-layout">
        <div class="task-description">
          <h2>📋 Описание задачи</h2>
          <p>{{ task?.description || 'Описание задачи' }}</p>
          
          <div class="examples" v-if="task?.solution">
            <h3>📝 Пример решения:</h3>
            <pre class="solution-code">{{ task.solution }}</pre>
          </div>

          <div class="task-tips">
            <h3>💡 Подсказки:</h3>
            <ul>
              <li>Пишите код который реально работает</li>
              <li>Для получения зеленого цвета - решите задачу правильно</li>
              <li>Ошибки подсвечиваются красным</li>
            </ul>
          </div>
        </div>

        <div class="editor-section">
          <h2>💻 Редактор кода</h2>
          <div class="editor-wrapper">
            <MonacoEditor
              v-model="code"
              :language="currentLanguage"
              :theme="editorTheme"
              @change="onCodeChange"
              ref="editorRef"
            />
          </div>
          <div class="editor-actions">
            <button class="btn btn-secondary" @click="toggleTheme">
              {{ editorTheme === 'vs-dark' ? '☀️ Светлая' : '🌙 Тёмная' }}
            </button>
            <button class="btn btn-primary" @click="runCode" :disabled="isRunning">
              {{ isRunning ? '⏳ Выполняется...' : '▶️ Выполнить' }}
            </button>
            <button class="btn btn-success" @click="checkSolution" :disabled="isRunning">
              ✅ Проверить
            </button>
          </div>
        </div>

        <div class="output-section">
          <h2>📊 Вывод</h2>
          <div class="output-container" :class="outputColor">
            <pre v-if="output" class="output">{{ output }}</pre>
            <div v-else class="no-output">
              <p>Здесь будет отображаться результат выполнения кода</p>
            </div>
          </div>
          
          <div class="output-info" :class="outputColor" v-if="output">
            <span class="status-icon">
              {{ outputColor === 'green' ? '✅' : outputColor === 'red' ? '❌' : '🟠' }}
            </span>
            <span class="status-text">
              {{ getStatusText(outputColor) }}
            </span>
          </div>

          <div class="execution-info" v-if="executionTime">
            <p>⏱ Время выполнения: {{ executionTime }}s</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import MonacoEditor from '@/components/editor/MonacoEditor.vue'
import { CodeRunner } from '@/services/codeRunner'
import { TASK_DATA } from '@/utils/constants'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()
const editorRef = ref(null)

const code = ref('')
const output = ref('')
const outputColor = ref('') // 'green', 'orange', 'red'
const currentLanguage = ref('python')
const editorTheme = ref('vs-dark')
const isRunning = ref(false)
const executionTime = ref(0)

// Получаем задачу из TASK_DATA
const task = computed(() => {
  const lang = route.params.lang
  const topic = route.params.topic
  const taskId = parseInt(route.params.taskId)
  
  console.log('Loading task:', { lang, topic, taskId })
  
  const tasks = TASK_DATA[lang]?.[topic] || []
  const foundTask = tasks.find(t => t.id === taskId)
  
  console.log('Found task:', foundTask)
  
  return foundTask
})

onMounted(() => {
  currentLanguage.value = route.params.lang
  setDefaultCode()
})

const setDefaultCode = () => {
  if (task.value?.defaultCode) {
    code.value = task.value.defaultCode
  } else {
    // Заглушка если код не найден
    const defaults = {
      python: '# Напишите ваш код здесь\nprint("Hello World")',
      javascript: '// Напишите ваш код здесь\nconsole.log("Hello World")',
      cpp: '// Напишите ваш код здесь\n#include <iostream>\n\nint main() {\n    std::cout << "Hello World" << std::endl;\n    return 0;\n}',
      java: '// Напишите ваш код здесь\npublic class Main {\n    public static void main(String[] args) {\n        System.out.println("Hello World");\n    }\n}'
    }
    code.value = defaults[currentLanguage.value] || defaults.python
  }
}

const goBack = () => {
  router.push(`/language/${route.params.lang}`)
}

const updateEditorLanguage = () => {
  setDefaultCode()
  output.value = ''
  outputColor.value = ''
  executionTime.value = 0
}

const toggleTheme = () => {
  editorTheme.value = editorTheme.value === 'vs-dark' ? 'vs' : 'vs-dark'
}

const onCodeChange = (newCode) => {
  code.value = newCode
}

const runCode = async () => {
  isRunning.value = true
  output.value = 'Отправка кода на выполнение...'
  outputColor.value = 'orange'
  executionTime.value = 0

  const startTime = Date.now()

  try {
    const result = await CodeRunner.executeCode(currentLanguage.value, code.value)
    output.value = result.output
    outputColor.value = result.success ? 'green' : 'red'
  } catch (error) {
    output.value = `❌ Ошибка: ${error.message}`
    outputColor.value = 'red'
  } finally {
    isRunning.value = false
    executionTime.value = ((Date.now() - startTime) / 1000).toFixed(2)
  }
}

const checkSolution = async () => {
  isRunning.value = true
  output.value = 'Проверка решения...'
  outputColor.value = 'orange'
  executionTime.value = 0

  const startTime = Date.now()

  try {
    const result = await CodeRunner.checkSolution(currentLanguage.value, code.value, task.value?.id)
    output.value = result.output
    outputColor.value = result.passed ? 'green' : 'red'
    
    if (result.passed) {
      output.value = '✅ Поздравляем! Задача решена правильно!\n\n' + result.output
    } else {
      output.value = '❌ Решение не прошло проверку\n\n' + result.output
    }
  } catch (error) {
    output.value = `❌ Ошибка проверки: ${error.message}`
    outputColor.value = 'red'
  } finally {
    isRunning.value = false
    executionTime.value = ((Date.now() - startTime) / 1000).toFixed(2)
  }
}

const getStatusText = (color) => {
  const statusMap = {
    green: '✅ Задача решена правильно!',
    orange: '🟠 Программа выполнена, но решение не оптимально',
    red: '❌ Ошибка выполнения'
  }
  return statusMap[color] || ''
}
</script>

<style scoped>
.task-view {
  min-height: 100vh;
  background: #f8f9fa;
  padding: 20px 0; /* ФИКС: убираем горизонтальный padding */
}

/* ФИКС: Контейнер с жесткими ограничениями */
.container {
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
  padding: 0 20px;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  padding: 20px;
  background: white;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
}

.back-btn {
  background: #6c757d;
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px 15px;
  cursor: pointer;
  font-size: 14px;
  transition: background 0.3s ease;
}

.back-btn:hover {
  background: #5a6268;
}

h1 {
  color: #333;
  font-size: 1.8rem;
  margin: 0;
}

.language-selector select {
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: white;
  font-size: 14px;
}

/* ФИКС: Grid layout с ограничениями */
.task-layout {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr;
  gap: 20px;
  height: calc(100vh - 150px);
  max-width: 100%;
  overflow: hidden;
}

.task-description, .editor-section, .output-section {
  background: white;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.1);
  display: flex;
  flex-direction: column;
}

.task-description h2, .editor-section h2, .output-section h2 {
  color: #333;
  margin-bottom: 15px;
  font-size: 1.2rem;
  border-bottom: 2px solid #f0f0f0;
  padding-bottom: 10px;
}

.examples {
  margin-top: 20px;
}

.solution-code {
  background: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 6px;
  padding: 15px;
  font-family: 'Courier New', monospace;
  font-size: 0.9rem;
  white-space: pre-wrap;
  overflow-x: auto;
  line-height: 1.4;
}

.task-tips {
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #eee;
}

.task-tips h3 {
  color: #666;
  font-size: 1rem;
  margin-bottom: 10px;
}

.task-tips ul {
  color: #666;
  padding-left: 20px;
}

.task-tips li {
  margin-bottom: 5px;
  font-size: 0.9rem;
  line-height: 1.4;
}

/* ФИКС: Обертка для редактора с жесткими ограничениями */
.editor-wrapper {
  width: 100%;
  max-width: 100%;
  overflow: hidden;
  position: relative;
  flex: 1;
  min-height: 400px;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
}

.editor-section {
  position: relative;
  min-height: 500px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* ФИКС: Стили для Monaco Editor */
:deep(.monaco-editor) {
  width: 100% !important;
  max-width: 100% !important;
}

:deep(.monaco-editor .overflow-guard) {
  width: 100% !important;
  max-width: 100% !important;
}

:deep(.monaco-editor .monaco-scrollable-element) {
  width: 100% !important;
  max-width: 100% !important;
  left: 0 !important;
}

.editor-actions {
  display: flex;
  gap: 10px;
  margin-top: 15px;
  flex-wrap: wrap;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.3s ease;
  font-weight: 500;
}

.btn-primary {
  background: #007bff;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #0056b3;
  transform: translateY(-1px);
}

.btn-primary:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

.btn-secondary {
  background: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background: #5a6268;
  transform: translateY(-1px);
}

.btn-success {
  background: #28a745;
  color: white;
}

.btn-success:hover:not(:disabled) {
  background: #218838;
  transform: translateY(-1px);
}

.btn-success:disabled {
  background: #ccc;
  cursor: not-allowed;
  transform: none;
}

.output-section {
  min-height: 400px;
}

.output-container {
  flex: 1;
  background: #1e1e1e;
  border-radius: 6px;
  padding: 15px;
  overflow-y: auto;
  border: 2px solid transparent;
  transition: border-color 0.3s ease;
  min-height: 200px;
}

.output-container.green {
  border-color: #28a745;
  background: #1a331a;
}

.output-container.orange {
  border-color: #fd7e14;
  background: #332211;
}

.output-container.red {
  border-color: #dc3545;
  background: #331a1a;
}

.output {
  color: #d4d4d4;
  font-family: 'Courier New', monospace;
  white-space: pre-wrap;
  margin: 0;
  line-height: 1.4;
  font-size: 13px;
}

.no-output {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #888;
  font-style: italic;
  text-align: center;
}

.output-info {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 15px;
  padding: 12px;
  border-radius: 6px;
  font-weight: 500;
}

.output-info.green {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.output-info.orange {
  background: #ffeaa7;
  color: #856404;
  border: 1px solid #ffecb3;
}

.output-info.red {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.status-icon {
  font-size: 1.2rem;
}

.status-text {
  font-size: 0.9rem;
}

.execution-info {
  margin-top: 10px;
  padding: 8px 12px;
  background: #e9ecef;
  border-radius: 4px;
  font-size: 0.85rem;
  color: #495057;
}

.execution-info p {
  margin: 0;
}

@media (max-width: 1200px) {
  .task-layout {
    grid-template-columns: 1fr 1fr;
    grid-template-rows: auto 1fr;
  }
  
  .task-description {
    grid-column: 1 / -1;
    height: auto;
  }
}

@media (max-width: 768px) {
  .task-layout {
    grid-template-columns: 1fr;
    grid-template-rows: auto 1fr auto;
  }
  
  .task-description {
    height: auto;
  }
  
  .header {
    flex-direction: column;
    gap: 15px;
    text-align: center;
  }
  
  .editor-actions {
    flex-direction: column;
  }
  
  .btn {
    width: 100%;
  }
}
</style>