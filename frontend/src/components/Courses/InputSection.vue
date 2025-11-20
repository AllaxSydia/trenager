<template>
  <section class="input-section">
    <div class="input-header">
      <h3>Консольный ввод</h3>
    </div>
    <div class="console-input">
      <textarea 
        :value="input"
        @input="$emit('update:input', $event.target.value)"
        class="input-field"
        placeholder="Введите входные данные здесь..."
        rows="3"
      ></textarea>
    </div>
  </section>
</template>

<script>
export default {
  name: 'InputSection',
  props: {
    input: String
  },
  emits: ['update:input']
}
</script>

<style scoped>
.input-section {
  background: #303030;
  border-radius: 16px;
  padding: 1.25rem;
  border: 1px solid #404040;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
  margin: 0;
}

.input-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.input-header h3 {
  margin: 0;
  color: #F8FAFC;
  font-size: 1.1rem;
  font-weight: 600;
}

.console-input {
  background: #1E1E1E;
  border: 1px solid #404040;
  border-radius: 12px;
  padding: 1rem;
  overflow: hidden;
}

.input-field {
  width: 100%;
  background: transparent;
  border: none;
  color: #E2E8F0;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 0.875rem;
  line-height: 1.5;
  resize: vertical;
  outline: none;
  min-height: 60px;
  transition: all 0.2s ease;
}

.input-field:focus {
  box-shadow: 0 0 0 2px #3B82F6;
}

.input-field::placeholder {
  color: #94A3B8;
}

@media (max-width: 768px) {
  .input-section {
    padding: 1rem;
    border-radius: 12px;
  }
}

@media (max-width: 480px) {
  .input-section {
    padding: 0.75rem;
  }
  
  .console-input {
    padding: 0.75rem;
  }
  
  .input-field {
    min-height: 50px;
  }
}
</style>