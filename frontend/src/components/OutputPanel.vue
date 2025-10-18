<template>
  <div class="output-panel">
    <h3>Результат выполнения:</h3>
    <div class="output-content" :class="resultClass">
      <pre>{{ outputText }}</pre>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useTaskStore } from '@/stores/taskStore'

const taskStore = useTaskStore()

const outputText = computed(() => {
  if (!taskStore.executionResult) {
    return 'Здесь будет результат выполнения вашего кода...'
  }
  
  if (taskStore.executionResult.success) {
    return taskStore.executionResult.output || taskStore.executionResult.message
  } else {
    return taskStore.executionResult.message
  }
})

const resultClass = computed(() => {
  if (!taskStore.executionResult) return ''
  return taskStore.executionResult.success ? 'success' : 'error'
})
</script>

<style scoped>
.output-panel {
  background: #252526;
  border-radius: 8px;
  padding: 16px;
  border: 1px solid #3c3c3c;
  /* ФИКСИРУЕМ ВЫСОТУ */
  height: 200px;
  display: flex;
  flex-direction: column;
}

h3 {
  margin: 0 0 12px 0;
  color: #cccccc;
  font-size: 1em;
  flex-shrink: 0;
}

.output-content {
  background: #1e1e1e;
  border-radius: 4px;
  padding: 12px;
  /* ЗАНИМАЕМ ВСЮ ОСТАВШУЮСЯ ВЫСОТУ */
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.output-content.success {
  border: 1px solid #4ec9b0;
}

.output-content.error {
  border: 1px solid #f44747;
}

pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 13px;
  white-space: pre-wrap;
  word-wrap: break-word;
  color: #cccccc;
}

.output-content.success pre {
  color: #4ec9b0;
}

.output-content.error pre {
  color: #f44747;
}
</style>