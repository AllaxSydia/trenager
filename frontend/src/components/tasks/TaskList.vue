<template>
  <div class="task-list">
    <div class="tasks">
      <div 
        v-for="task in tasks"
        :key="task.id"
        class="task-card"
        @click="selectTask(task.id)"
      >
        <h3>{{ task.title }}</h3>
        <p>{{ task.description }}</p>
        <span class="difficulty" :class="task.difficulty">
          {{ getDifficultyText(task.difficulty) }}
        </span>
      </div>
    </div>

    <div v-if="tasks.length === 0" class="no-tasks">
      <p>Задачи для этого языка пока не добавлены</p>
    </div>
  </div>
</template>

<script setup>
import { TASK_DATA } from '@/utils/constants'
import { computed } from 'vue'
import { useRouter } from 'vue-router'

const props = defineProps({
  language: String
})

const router = useRouter()

const tasks = computed(() => {
  return TASK_DATA[props.language] || []
})

const selectTask = (taskId) => {
  router.push(`/task/${props.language}/${taskId}`)
}

const getDifficultyText = (difficulty) => {
  const map = { easy: 'Легкая', medium: 'Средняя', hard: 'Сложная' }
  return map[difficulty] || difficulty
}
</script>

<style scoped>
.task-list {
  background: white;
  border-radius: 12px;
  padding: 30px;
}

.tasks {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.task-card {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s ease;
  border-left: 4px solid #007bff;
}

.task-card:hover {
  background: #e9ecef;
  transform: translateX(5px);
}

.task-card h3 {
  color: #333;
  margin-bottom: 8px;
}

.task-card p {
  color: #666;
  margin-bottom: 10px;
}

.difficulty {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: bold;
}

.difficulty.easy {
  background: #d4edda;
  color: #155724;
}

.difficulty.medium {
  background: #fff3cd;
  color: #856404;
}

.difficulty.hard {
  background: #f8d7da;
  color: #721c24;
}

.no-tasks {
  text-align: center;
  padding: 40px;
  color: #666;
}
</style>