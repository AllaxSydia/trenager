<template>
  <div class="output-panel">
    <div class="output-header">
      <h3>Результат:</h3>
      <div class="controls">
        <button 
          @click="runCode" 
          :disabled="isRunning"
          class="run-btn"
        >
          {{ isRunning ? 'Выполняется...' : 'Запустить' }}
        </button>
        <button 
          @click="checkSolution" 
          :disabled="isRunning"
          class="check-btn"
        >
          Проверить
        </button>
      </div>
    </div>

    <div class="output-content">
      <!-- Результат выполнения -->
      <div v-if="output" class="output-section">
        <pre :class="{ 'error': !lastSuccess }">{{ output }}</pre>
      </div>

      <!-- Результат проверки -->
      <div v-if="checkResult" class="check-result">
        <div :class="['result-badge', checkResult.passed ? 'success' : 'error']">
          {{ checkResult.message }}
        </div>
        <div v-if="!checkResult.passed" class="test-details">
          <div class="expected">Ожидалось: <code>{{ checkResult.expected }}</code></div>
          <div class="actual">Получено: <code>{{ checkResult.actual }}</code></div>
        </div>
      </div>

      <!-- Состояние загрузки -->
      <div v-if="isRunning" class="loading">
        Выполнение кода...
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineEmits, defineProps, ref } from 'vue'

const props = defineProps({
  code: String,
  language: String,
  taskId: String
})

const emit = defineEmits(['codeExecuted'])

const isRunning = ref(false)
const output = ref('')
const lastSuccess = ref(true)
const checkResult = ref(null)

const runCode = async () => {
  isRunning.value = true
  output.value = ''
  checkResult.value = null

  try {
    const response = await fetch('/api/execute', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        code: props.code,
        language: props.language,
        task_id: props.taskId
      })
    })
    
    const result = await response.json()
    output.value = result.output
    lastSuccess.value = result.success
    emit('codeExecuted', result)
  } catch (error) {
    output.value = 'Ошибка соединения: ' + error.message
    lastSuccess.value = false
  } finally {
    isRunning.value = false
  }
}

const checkSolution = async () => {
  isRunning.value = true
  output.value = ''
  checkResult.value = null

  try {
    const response = await fetch('/api/check', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        code: props.code,
        language: props.language,
        task_id: props.taskId || 'hello_world'
      })
    })
    
    const result = await response.json()
    output.value = result.output
    lastSuccess.value = result.success
    checkResult.value = result
    emit('codeExecuted', result)
  } catch (error) {
    output.value = 'Ошибка проверки: ' + error.message
    lastSuccess.value = false
  } finally {
    isRunning.value = false
  }
}

// Экспортируем методы для использования в родительском компоненте
defineExpose({
  runCode,
  checkSolution
})
</script>

<style scoped>
.output-panel {
  border: 1px solid #e1e5e9;
  border-radius: 8px;
  background: white;
}

.output-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #e1e5e9;
  background: #f8f9fa;
}

.output-header h3 {
  margin: 0;
  font-size: 16px;
  color: #2c3e50;
}

.controls {
  display: flex;
  gap: 8px;
}

button {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  background: white;
  cursor: pointer;
  font-size: 14px;
  transition: all 0.2s;
}

button:hover:not(:disabled) {
  background: #f8f9fa;
}

button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.run-btn {
  border-color: #3498db;
  color: #3498db;
}

.check-btn {
  border-color: #27ae60;
  color: #27ae60;
}

.output-content {
  padding: 16px;
  min-height: 100px;
}

.output-section pre {
  margin: 0;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 4px;
  white-space: pre-wrap;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 13px;
}

.output-section pre.error {
  background: #fee;
  color: #c53030;
}

.check-result {
  margin-top: 12px;
  padding: 12px;
  border-radius: 6px;
}

.result-badge {
  padding: 8px 12px;
  border-radius: 4px;
  font-weight: 600;
  text-align: center;
}

.result-badge.success {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.result-badge.error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.test-details {
  margin-top: 8px;
  padding: 8px;
  background: #f8f9fa;
  border-radius: 4px;
  font-size: 14px;
}

.test-details code {
  background: #e9ecef;
  padding: 2px 6px;
  border-radius: 3px;
  font-family: monospace;
}

.loading {
  padding: 12px;
  text-align: center;
  color: #6c757d;
  font-style: italic;
}
</style>