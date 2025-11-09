<template>
  <section class="code-section">
    <div class="code-header">
      <h3>Редактор кода</h3>
      <div class="code-actions">
        <button class="btn btn--small" @click="$emit('reset')">Сбросить</button>
        <button class="btn btn--small" @click="$emit('execute')">▶ Запустить</button>
      </div>
    </div>
    <div class="code-editor-container">
      <CodeMirrorEditor
        :modelValue="code"
        :language="language"
        @update:modelValue="$emit('update:code', $event)"
      />
    </div>
  </section>
</template>

<script>
import CodeMirrorEditor from './CodeMirrorEditor.vue';

export default {
  name: 'CodeSection',
  components: {
    CodeMirrorEditor
  },
  props: {
    code: String,
    language: {
      type: String,
      default: 'python'
    }
  },
  emits: ['update:code', 'reset', 'execute']
}
</script>

<style scoped>
/* Стили остаются такими же как были */
.code-section {
  background: #303030;
  border-radius: 16px;
  padding: 1.25rem;
  border: 1px solid #404040;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
  margin: 0;
}

.code-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.code-header h3 {
  margin: 0;
  color: #F8FAFC;
  font-size: 1.1rem;
  font-weight: 600;
}

.code-actions {
  display: flex;
  gap: 0.5rem;
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

.btn--small {
  padding: 0.45rem 0.9rem;
  font-size: 0.7rem;
  border-radius: 8px;
  flex: none;
  background: #3B82F6;
  color: white;
}

.btn--small:hover {
  background: #2563EB;
}

.code-editor-container {
  position: relative;
  background: #1E1E1E;
  border-radius: 12px;
  border: 1px solid #404040;
  height: 400px;
  overflow: hidden;
}

@media (max-width: 768px) {
  .code-section {
    padding: 1rem;
  }
  
  .code-header {
    flex-direction: column;
    gap: 0.5rem;
    align-items: flex-start;
  }
  
  .code-actions {
    align-self: flex-end;
  }
  
  .code-editor-container {
    height: 300px;
  }
}
</style>