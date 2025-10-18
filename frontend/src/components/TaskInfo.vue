<template>
  <div class="task-info" v-if="currentTask">
    <h3>{{ currentTask.title }}</h3>
    <p>{{ currentTask.description }}</p>
  </div>
  <div class="task-info" v-else>
    <p>Выберите задачу из списка слева</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useTaskStore } from '@/stores/taskStore'

const taskStore = useTaskStore()
const currentTask = computed(() => taskStore.currentTask)
</script>

<style scoped>
.task-info {
  background: #252526;
  border-radius: 8px;
  padding: 20px;
  border: 1px solid #3c3c3c;
}

h3 {
  margin: 0 0 10px 0;
  color: #ffffff;
}

p {
  margin: 0;
  color: #cccccc;
  line-height: 1.5;
}
</style>