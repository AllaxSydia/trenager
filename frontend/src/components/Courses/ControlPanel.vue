<template>
  <section class="control-panel">
    <div class="control-buttons">
      <button 
        class="btn btn--primary" 
        @click="$emit('run')"
        :disabled="isRunning"
      >
        Запустить
      </button>
      <button 
        class="btn btn--secondary" 
        @click="$emit('test')"
        :disabled="isTesting"
      >
        Проверить тесты
      </button>
      <button 
        class="btn btn--success" 
        @click="$emit('submit')"
        :disabled="isSubmitting"
      >
        Отправить решение
      </button>
      <button 
        class="btn btn--outline" 
        @click="$emit('save')"
      >
        Сохранить
      </button>
    </div>
  </section>
</template>

<script>
export default {
  name: 'ControlPanel',
  props: {
    isRunning: Boolean,
    isTesting: Boolean,
    isSubmitting: Boolean
  },
  emits: ['run', 'test', 'submit', 'save']
}
</script>

<style scoped>
.control-panel {
  background: #303030;
  border-radius: 16px;
  padding: 1.25rem;
  border: 1px solid #404040;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
  margin: 0;
}

.control-buttons {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
  width: 100%;
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
  flex: 1;
  min-width: 120px;
  position: relative;
  overflow: hidden;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}

.btn:active:not(:disabled) {
  transform: translateY(2px);
}

.btn--primary {
  background: #3B82F6;
  color: white;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}

.btn--primary:hover:not(:disabled) {
  background: #2563EB;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(37, 99, 235, 0.4);
}

.btn--primary:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.4);
}

.btn--secondary {
  background: #303030;
  color: #E2E8F0;
  border: 1px solid #404040;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.btn--secondary:hover:not(:disabled) {
  background: #404040;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 255, 255, 0.1);
}

.btn--secondary:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
}

.btn--success {
  background: #303030;
  color: #E2E8F0;
  border: 1px solid #404040;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.btn--success:hover:not(:disabled) {
  background: #404040;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 255, 255, 0.1);
}

.btn--success:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
}

.btn--outline {
  background: #303030;
  color: #E2E8F0;
  border: 1px solid #404040;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.btn--outline:hover:not(:disabled) {
  background: #404040;
  color: #E2E8F0;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 255, 255, 0.1);
}

.btn--outline:active:not(:disabled) {
  transform: translateY(0);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
}

@media (max-width: 768px) {
  .control-panel {
    padding: 1rem;
    border-radius: 12px;
  }
  
  .control-buttons {
    flex-direction: column;
  }
  
  .btn {
    justify-content: center;
    min-width: auto;
  }
}

@media (max-width: 480px) {
  .control-panel {
    padding: 0.75rem;
  }
  
  .btn {
    padding: 0.6rem 1rem;
    font-size: 0.8rem;
  }
}
</style>