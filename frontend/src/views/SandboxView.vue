<template>
  <div class="sandbox-view">
    <div class="container">
      <header class="header">
        <button class="back-btn" @click="$router.push('/')">← На главную</button>
        <h1>🏖️ Режим песочницы</h1>
      </header>
      
      <div class="sandbox-content">
        <p>Свободное программирование без ограничений</p>
        <textarea v-model="code" class="code-editor" placeholder="Введите ваш код здесь..."></textarea>
        <button class="btn btn-primary" @click="runCode">Выполнить код</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const code = ref('')

const runCode = () => {
  alert('Выполняется код: ' + code.value)
}
</script>

<style scoped>
.sandbox-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #a29bfe 0%, #6c5ce7 100%);
  padding: 30px 20px;
}

.header {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 30px;
}

.back-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px 15px;
  cursor: pointer;
  backdrop-filter: blur(10px);
}

h1 {
  color: white;
  font-size: 2rem;
}

.sandbox-content {
  background: white;
  border-radius: 12px;
  padding: 30px;
  text-align: center;
}

.sandbox-content p {
  color: #666;
  margin-bottom: 20px;
  font-size: 1.1rem;
}

.code-editor {
  width: 100%;
  height: 400px;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 20px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  margin-bottom: 20px;
  resize: vertical;
}
</style>