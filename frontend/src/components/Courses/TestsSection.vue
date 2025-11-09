<template>
  <section class="tests-section">
    <div class="tests-header">
      <h3>Тесты</h3>
      <div class="tests-summary">
        <span class="tests-passed">Пройдено: {{ passedTests }}/{{ tests.length }}</span>
      </div>
    </div>
    <div class="tests-list">
      <div 
        v-for="(test, index) in tests" 
        :key="index"
        class="test-item"
        :class="{
          'test-item--passed': test.status === 'passed',
          'test-item--failed': test.status === 'failed',
          'test-item--running': test.status === 'running'
        }"
      >
        <div class="test-icon">
          <span v-if="test.status === 'passed'">✓</span>
          <span v-else-if="test.status === 'failed'">✗</span>
          <span v-else-if="test.status === 'running'">⏳</span>
          <span v-else>○</span>
        </div>
        <div class="test-content">
          <div class="test-input">
            <strong>Вход:</strong> {{ test.input }}
          </div>
          <div class="test-expected">
            <strong>Ожидается:</strong> {{ test.expected }}
          </div>
          <div v-if="test.actual" class="test-actual">
            <strong>Получено:</strong> {{ test.actual }}
          </div>
          <div v-if="test.error" class="test-error">
            <strong>Ошибка:</strong> {{ test.error }}
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script>
export default {
  name: 'TestsSection',
  props: {
    tests: Array,
    passedTests: Number
  }
}
</script>

<style scoped>
.tests-section {
  background: #303030;
  border-radius: 16px;
  padding: 1.25rem;
  border: 1px solid #404040;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
  margin: 0;
  position: relative;
}

.tests-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  position: relative;
  z-index: 2;
}

.tests-header h3 {
  margin: 0;
  color: #F8FAFC;
  font-size: 1.1rem;
  font-weight: 600;
}

.tests-summary {
  color: #94A3B8;
  font-size: 0.875rem;
  background: rgba(30, 41, 59, 0.8);
  padding: 0.4rem 0.8rem;
  border-radius: 8px;
  border: 1px solid #404040;
  position: relative;
  z-index: 3;
}

.tests-passed {
  color: #3B82F6;
  font-weight: 600;
}

.tests-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  overflow: visible;
  margin-top: 0.5rem;
}

.test-item {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 1rem;
  background: #1E1E1E;
  border-radius: 12px;
  border: 1px solid #404040;
  transition: all 0.2s ease;
  position: relative;
  z-index: 1;
}

.test-item:hover {
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.test-item:active {
  transform: translateY(0);
}

.test-item--passed {
  border-color: #3B82F6;
  background: rgba(59, 130, 246, 0.05);
}

.test-item--failed {
  border-color: #EF4444;
  background: rgba(239, 68, 68, 0.05);
}

.test-item--running {
  border-color: #F59E0B;
  background: rgba(245, 158, 11, 0.05);
}

.test-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  flex-shrink: 0;
  margin-top: 0.125rem;
  transition: all 0.2s ease;
}

.test-item:active .test-icon {
  transform: scale(0.95);
}

.test-item--passed .test-icon {
  background: #3B82F6;
  color: white;
}

.test-item--failed .test-icon {
  background: #EF4444;
  color: white;
}

.test-item--running .test-icon {
  background: #F59E0B;
  color: white;
}

.test-content {
  flex: 1;
  font-size: 0.875rem;
}

.test-input,
.test-expected,
.test-actual,
.test-error {
  margin-bottom: 0.25rem;
}

.test-input strong,
.test-expected strong,
.test-actual strong,
.test-error strong {
  color: #F8FAFC;
}

.test-actual {
  color: #3B82F6;
}

.test-error {
  color: #EF4444;
}

@media (max-width: 768px) {
  .tests-section {
    padding: 1rem;
    border-radius: 12px;
  }
  
  .tests-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .tests-summary {
    align-self: flex-end;
  }
  
  .test-item {
    padding: 0.75rem;
    flex-direction: column;
    gap: 0.5rem;
  }
  
  .test-icon {
    align-self: flex-start;
  }
}

@media (max-width: 480px) {
  .tests-section {
    padding: 0.75rem;
  }
  
  .test-item {
    padding: 0.75rem;
  }
}
</style>