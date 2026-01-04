<template>
  <section class="output-section">
    <div class="output-header">
      <div class="tabs">
        <button 
          class="tab" 
          :class="{ 'tab--active': activeTab === 'console' }"
          @click="activeTab = 'console'"
        >
          Консоль
        </button>
        <button 
          class="tab" 
          :class="{ 'tab--active': activeTab === 'ai' }"
          @click="activeTab = 'ai'"
        >
          AI Анализ
        </button>
      </div>
      <button 
        v-if="activeTab === 'console'" 
        class="btn btn--small" 
        @click="$emit('clear')"
      >
        Очистить
      </button>
    </div>

    <div v-if="activeTab === 'console'" class="console-output">
      <pre>{{ output }}</pre>
    </div>

    <div v-else-if="activeTab === 'ai' && aiAnalysis" class="ai-analysis">
      <div v-if="aiLoading" class="ai-loading">
        <div class="spinner"></div>
        <p>AI анализирует код...</p>
      </div>
      
      <div v-else class="ai-results">
        <!-- Оценка -->
        <div class="ai-score">
          <div class="score-circle" :class="getScoreClass(aiAnalysis.score)">
            {{ aiAnalysis.score }}/10
          </div>
          <div class="score-label">
            <h4>Оценка</h4>
            <p>Сложность: {{ aiAnalysis.complexity }}</p>
          </div>
        </div>

        <!-- Комментарии -->
        <div class="ai-section">
          <h4>Комментарии</h4>
          <ul>
            <li v-for="(comment, index) in aiAnalysis.comments" :key="index">
              {{ comment }}
            </li>
          </ul>
        </div>

        <!-- Предложения -->
        <div class="ai-section">
          <h4>Предложения по улучшению</h4>
          <ul>
            <li v-for="(suggestion, index) in aiAnalysis.suggestions" :key="index">
              {{ suggestion }}
            </li>
          </ul>
        </div>

        <!-- Best Practices -->
        <div class="ai-section">
          <h4>Best Practices</h4>
          <ul>
            <li v-for="(practice, index) in aiAnalysis.best_practices" :key="index">
              {{ practice }}
            </li>
          </ul>
        </div>

        <!-- Альтернативные решения -->
        <div v-if="aiAnalysis.alternative_solutions && aiAnalysis.alternative_solutions.length" class="ai-section">
          <h4>Альтернативные решения</h4>
          <ul>
            <li v-for="(solution, index) in aiAnalysis.alternative_solutions" :key="index">
              {{ solution }}
            </li>
          </ul>
        </div>
      </div>
    </div>

    <div v-else-if="activeTab === 'ai'" class="ai-empty">
      <p>Нажмите "Проанализировать код" для получения AI-анализа</p>
    </div>
  </section>
</template>

<script>
export default {
  name: 'OutputSection',
  props: {
    output: String,
    aiAnalysis: Object,
    aiLoading: Boolean
  },
  emits: ['clear'],
  data() {
    return {
      activeTab: 'console'
    }
  },
  methods: {
    getScoreClass(score) {
      if (score >= 8) return 'score-high';
      if (score >= 6) return 'score-medium';
      return 'score-low';
    }
  }
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
  height: 100%;
  display: flex;
  flex-direction: column;
}

.output-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  flex-shrink: 0;
}

.tabs {
  display: flex;
  gap: 0.5rem;
  background: #404040;
  padding: 0.25rem;
  border-radius: 8px;
}

.tab {
  padding: 0.5rem 1rem;
  border: none;
  background: transparent;
  color: #94a3b8;
  border-radius: 6px;
  font-weight: 600;
  font-size: 0.85rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.tab--active {
  background: #6366f1;
  color: white;
}

.tab:hover:not(.tab--active) {
  background: #4b5563;
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
  background: #dc2626;
  color: white;
}

.btn--small:hover {
  background: #b91c1c;
}

.btn--small:active:not(:disabled) {
  transform: translateY(1px);
}

.console-output, .ai-analysis, .ai-empty {
  background: #1E1E1E;
  border: 1px solid #404040;
  border-radius: 12px;
  padding: 1rem;
  color: #E2E8F0;
  font-family: 'Monaco', 'Consolas', monospace;
  font-size: 0.875rem;
  overflow-y: auto;
  flex-grow: 1;
}

.console-output pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

/* AI Analysis Styles */
.ai-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #94a3b8;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #404040;
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.ai-results {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.ai-score {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #404040;
}

.score-circle {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 1.1rem;
  color: white;
}

.score-high {
  background: linear-gradient(135deg, #10b981, #059669);
}

.score-medium {
  background: linear-gradient(135deg, #f59e0b, #d97706);
}

.score-low {
  background: linear-gradient(135deg, #ef4444, #dc2626);
}

.score-label h4 {
  margin: 0;
  color: #f8fafc;
  font-size: 1rem;
}

.score-label p {
  margin: 0.25rem 0 0 0;
  color: #94a3b8;
  font-size: 0.85rem;
}

.ai-section h4 {
  color: #f8fafc;
  margin: 0 0 0.75rem 0;
  font-size: 1rem;
  font-weight: 600;
}

.ai-section ul {
  margin: 0;
  padding-left: 1.25rem;
}

.ai-section li {
  color: #e2e8f0;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
  line-height: 1.4;
}

.ai-section li:last-child {
  margin-bottom: 0;
}

.ai-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: #94a3b8;
  text-align: center;
}

@media (max-width: 768px) {
  .output-section {
    padding: 1rem;
    border-radius: 12px;
  }
  
  .output-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }
  
  .tabs {
    align-self: stretch;
  }
  
  .tab {
    flex: 1;
    text-align: center;
  }
  
  .output-header .btn--small {
    align-self: flex-end;
  }
  
  .console-output, .ai-analysis {
    min-height: 200px;
  }
  
  .ai-score {
    flex-direction: column;
    text-align: center;
    gap: 0.75rem;
  }
}

@media (max-width: 480px) {
  .output-section {
    padding: 0.75rem;
  }
  
  .console-output, .ai-analysis {
    padding: 0.75rem;
  }
}
</style>