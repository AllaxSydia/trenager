<template>
  <section class="problem-section">
    <div class="problem-header">
      <h1>{{ lesson.title }}</h1>
      <div class="difficulty" :class="'difficulty--' + lesson.difficulty">
        {{ difficultyText[lesson.difficulty] }}
      </div>
    </div>
    <div class="problem-description">
      <p>{{ lesson.description }}</p>
    </div>
  </section>
</template>

<script>
export default {
  name: 'ProblemSection',
  props: {
    lesson: Object
  },
  data() {
    return {
      difficultyText: {
        easy: 'Легко',
        medium: 'Средне',
        hard: 'Сложно'
      }
    }
  }
}
</script>

<style scoped>
.problem-section {
  background: #303030;
  border-radius: 16px;
  padding: 1.25rem;
  border: 1px solid #404040;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
  margin: 0;
}

.problem-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.problem-header h1 {
  margin: 0;
  color: #F8FAFC;
  font-size: 1.4rem;
  font-weight: 600;
}

.difficulty {
  padding: 0.35rem 0.875rem;
  border-radius: 16px;
  font-size: 0.75rem;
  font-weight: 600;
}

.difficulty--easy {
  background: #3B82F6;
  color: white;
}

.difficulty--medium {
  background: #8B5CF6;
  color: white;
}

.difficulty--hard {
  background: #EF4444;
  color: white;
}

.problem-description {
  line-height: 1.6;
  color: #E2E8F0;
  font-size: 0.95rem;
}

@media (max-width: 768px) {
  .problem-section {
    padding: 1rem;
    border-radius: 12px;
  }
  
  .problem-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
  
  .problem-header h1 {
    font-size: 1.2rem;
  }
  
  .difficulty {
    align-self: flex-start;
  }
}

@media (max-width: 480px) {
  .problem-section {
    padding: 0.75rem;
  }
}
</style>