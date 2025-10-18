<template>
  <div>
    <header class="app-header">
      <h1>🚀 Тренажёр программирования</h1>
    </header>
    <main class="app-main">
      <div class="container">
        <div class="layout">
          <TaskList />
          <div class="editor-area">
            <TaskInfo />
            <CodeEditor />
            <OutputPanel />
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import TaskList from '@/components/TaskList.vue'
import TaskInfo from '@/components/TaskInfo.vue'
import CodeEditor from '@/components/CodeEditor.vue'
import OutputPanel from '@/components/OutputPanel.vue'
import { useTaskStore } from '@/stores/taskStore'
import { onMounted } from 'vue'

const taskStore = useTaskStore()

onMounted(() => {
  taskStore.loadTasks()
})
</script>

<style scoped>
/* Стили из App.vue переносим сюда */
.app-header {
  background: #252526;
  padding: 15px 20px;
  border-bottom: 1px solid #464647;
  flex-shrink: 0;
}

.app-header h1 {
  margin: 0;
  font-size: 1.5em;
  color: #ffffff;
}

.app-main {
  flex: 1;
  overflow: hidden;
}

.container {
  height: 100%;
  padding: 20px;
}

.layout {
  display: grid;
  grid-template-columns: 350px 1fr;
  gap: 20px;
  height: 100%;
}

.editor-area {
  display: flex;
  flex-direction: column;
  gap: 15px;
  height: 100%;
}
</style>