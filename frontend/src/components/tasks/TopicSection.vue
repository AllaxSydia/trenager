<template>
  <div class="topic-section">
    <div class="topic-header" :style="{ borderLeftColor: difficultyColor }">
      <h2 class="topic-title">
        <span class="topic-icon">{{ topicIcon }}</span>
        {{ topicName }}
      </h2>
      <span class="difficulty-badge" :style="{ background: difficultyColor }">
        {{ difficultyName }}
      </span>
    </div>
    
    <div class="tasks-grid">
      <div
        v-for="task in tasks"
        :key="task.id"
        class="task-card"
        @click="selectTask(task.id)"
      >
        <div class="task-header">
          <h3 class="task-title">{{ task.title }}</h3>
          <span class="task-difficulty" :style="{ background: getTaskDifficultyColor(task.difficulty) }">
            {{ getDifficultyText(task.difficulty) }}
          </span>
        </div>
        <p class="task-description">{{ task.description }}</p>
        <div class="task-footer">
          <span class="task-id">#{{ task.id }}</span>
          <button class="start-btn">Решить →</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { DIFFICULTY_LEVELS, TOPICS } from '@/utils/constants'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  language: {
    type: String,
    required: true
  },
  topic: {
    type: String,
    required: true
  },
  tasks: {
    type: Array,
    default: () => []
  }
})

const router = useRouter()

const topicName = computed(() => TOPICS[props.topic]?.name || props.topic)
const topicIcon = computed(() => {
  const icons = {
    output: '📤',
    input: '📥', 
    conditions: '🔀',
    loops: '🔄',
    arrays: '📊',
    functions: '⚙️',
    oop: '🏗️',
    algorithms: '🧩'
  }
  return icons[props.topic] || '📝'
})

const difficultyColor = computed(() => {
  const difficulties = props.tasks.map(t => t.difficulty)
  if (difficulties.includes('advanced')) return DIFFICULTY_LEVELS.advanced.color
  if (difficulties.includes('intermediate')) return DIFFICULTY_LEVELS.intermediate.color
  return DIFFICULTY_LEVELS.beginner.color
})

const difficultyName = computed(() => {
  const difficulties = props.tasks.map(t => t.difficulty)
  if (difficulties.includes('advanced')) return DIFFICULTY_LEVELS.advanced.name
  if (difficulties.includes('intermediate')) return DIFFICULTY_LEVELS.intermediate.name
  return DIFFICULTY_LEVELS.beginner.name
})

const selectTask = (taskId) => {
  router.push(`/task/${props.language}/${props.topic}/${taskId}`)
}

const getTaskDifficultyColor = (difficulty) => {
  return DIFFICULTY_LEVELS[difficulty]?.color || '#6c757d'
}

const getDifficultyText = (difficulty) => {
  return DIFFICULTY_LEVELS[difficulty]?.name || difficulty
}
</script>

<style scoped>
.topic-section {
  margin-bottom: 40px;
}

.topic-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-left: 15px;
  border-left: 4px solid #007bff;
}

.topic-title {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #333;
  font-size: 1.5rem;
  margin: 0;
}

.topic-icon {
  font-size: 1.8rem;
}

.difficulty-badge {
  color: white;
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: bold;
}

.tasks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 20px;
}

.task-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  border: 1px solid #e9ecef;
}

.task-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 15px rgba(0,0,0,0.15);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.task-title {
  color: #333;
  font-size: 1.1rem;
  margin: 0;
  flex: 1;
}

.task-difficulty {
  color: white;
  padding: 3px 8px;
  border-radius: 12px;
  font-size: 0.7rem;
  font-weight: bold;
  margin-left: 10px;
}

.task-description {
  color: #666;
  margin-bottom: 15px;
  line-height: 1.4;
  font-size: 0.9rem;
}

.task-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.task-id {
  color: #999;
  font-size: 0.8rem;
}

.start-btn {
  background: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  padding: 6px 12px;
  cursor: pointer;
  font-size: 0.8rem;
  transition: background 0.3s ease;
}

.start-btn:hover {
  background: #0056b3;
}
</style>