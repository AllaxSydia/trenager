<template>
  <div class="task-list">
    <h3>Задачи</h3>
    <div 
      v-for="task in tasks" 
      :key="task.id"
      :class="['task-item', { active: isActive(task) }]"
      @click="selectTask(task)"
    >
      <strong>{{ task.title }}</strong>
      <p>{{ task.description }}</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useTaskStore } from '@/stores/taskStore'

const taskStore = useTaskStore()

const tasks = computed(() => taskStore.tasks)

function selectTask(task) {
  taskStore.setCurrentTask(task)
}

function isActive(task) {
  return taskStore.currentTask?.id === task.id
}
</script>

<style scoped>
.task-list {
  background: #252526;
  border-radius: 8px;
  padding: 15px;
  overflow-y: auto;
  height: 100%;
}

.task-item {
  padding: 15px;
  border: 1px solid #3c3c3c;
  margin: 10px 0;
  cursor: pointer;
  border-radius: 5px;
  transition: all 0.3s ease;
}

.task-item:hover {
  background: #2a2d2e;
  border-color: #007acc;
}

.task-item.active {
  background: #04395e;
  border-color: #007acc;
}

h3 {
  margin-top: 0;
  color: #cccccc;
}

p {
  margin: 5px 0 0 0;
  font-size: 0.9em;
  color: #999;
}
</style>