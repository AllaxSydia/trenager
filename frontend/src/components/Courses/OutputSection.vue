<template>
  <section class="output-section">
    <div class="output-header">
      <h3>Консольный вывод</h3>
      <button class="btn btn--small" @click="$emit('clear')">Очистить</button>
    </div>
    <div class="console-output">
      <pre>{{ output }}</pre>
    </div>
  </section>
</template>

<script>
export default {
  name: 'OutputSection',
  props: {
    output: String
  },
  emits: ['clear']
}
</script>

<style scoped>
.output-section {
  background: #303030;
  border-radius: 16px;
  padding: 1.25rem;
  border: 1px solid #404040;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
  margin: 0;
}

.output-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.output-header h3 {
  margin: 0;
  color: #F8FAFC;
  font-size: 1.1rem;
  font-weight: 600;
}

.btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0.7rem 1.1rem;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn--small {
  padding: 0.45rem 0.9rem;
  font-size: 0.7rem;
  border-radius: 8px;
  flex: none;
}

.btn--small:active:not(:disabled) {
  transform: translateY(1px);
}

.console-output {
  background: #1E1E1E;
  border: 1px solid #404040;
  border-radius: 12px;
  padding: 1rem;
  color: #E2E8F0;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 0.875rem;
  min-height: 100px;
  overflow: visible;
}

.console-output pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

@media (max-width: 768px) {
  .output-section {
    padding: 1rem;
    border-radius: 12px;
  }
  
  .output-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .output-header .btn--small {
    align-self: flex-end;
  }
  
  .console-output {
    min-height: 80px;
  }
}

@media (max-width: 480px) {
  .output-section {
    padding: 0.75rem;
  }
  
  .console-output {
    padding: 0.75rem;
    min-height: 80px;
  }
}
</style>