<template>
  <div class="code-editor">
    <!-- НОВАЯ КРАСИВАЯ ШАПКА -->
    <div class="editor-header">
      <div class="header-left">
        <div class="language-selector">
          <label for="language">Язык:</label>
          <select id="language" v-model="language" @change="onLanguageChange">
            <option value="python">🐍 Python</option>
            <option value="javascript">🟨 JavaScript</option>
            <option value="java">☕ Java</option>
            <option value="cpp">⚙️ C++</option>
            <option value="go">🐹 Go</option>
          </select>
        </div>
        
        <div class="status-indicator" :class="{ running: isLoading }">
          <div class="status-dot"></div>
          <span>{{ isLoading ? 'Выполняется...' : 'Готов' }}</span>
        </div>
      </div>
      
      <div class="header-right">
        <button class="run-button" @click="execute" :disabled="isLoading">
          <span class="button-icon">🚀</span>
          <span class="button-text">{{ isLoading ? 'Выполнение...' : 'Выполнить код' }}</span>
          <span class="shortcut">Ctrl+Enter</span>
        </button>
      </div>
    </div>
    
    <!-- КРАСИВЫЙ РЕДАКТОР С НУМЕРАЦИЕЙ -->
    <div class="editor-container">
      <div class="editor-wrapper">
        <div class="editor-line-numbers">
          <div v-for="line in lineCount" :key="line" class="line-number">
            {{ line }}
          </div>
        </div>
        <textarea 
          v-model="code"
          @input="onCodeChange"
          @keydown="handleKeydown"
          placeholder="// Начните писать код здесь..."
          class="text-editor"
          ref="editorTextarea"
        ></textarea>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useTaskStore } from '@/stores/taskStore'

const taskStore = useTaskStore()
const editorTextarea = ref(null)

const language = ref('python')
const code = ref('')

const isLoading = computed(() => taskStore.isLoading)

// Считаем количество строк для нумерации
const lineCount = computed(() => {
  return code.value.split('\n').length
})

// Следим за сменой задачи
watch(() => taskStore.currentTask, (newTask) => {
  if (newTask) {
    updateTemplateForCurrentTask()
  }
}, { immediate: true })

// Следим за изменениями кода в store
watch(() => taskStore.userCode, (newCode) => {
  if (newCode !== code.value) {
    code.value = newCode
  }
})

function onCodeChange() {
  taskStore.updateCode(code.value)
}

function onLanguageChange() {
  taskStore.setLanguage(language.value)
  if (taskStore.currentTask) {
    updateTemplateForCurrentTask()
  }
}

// ОБНОВЛЕННАЯ ФУНКЦИЯ ДЛЯ ОБНОВЛЕНИЯ ШАБЛОНА
function updateTemplateForCurrentTask() {
  const templates = {
    python: taskStore.currentTask?.template || '# Напишите ваш код здесь\nprint("Hello, World!")',
    javascript: getJSTemplate(),
    java: getJavaTemplate(),
    cpp: getCppTemplate(),
    go: getGoTemplate()
  }
  
  code.value = templates[language.value] || templates.python
  taskStore.updateCode(code.value)
}

// ШАБЛОНЫ ДЛЯ РАЗНЫХ ЯЗЫКОВ
function getJSTemplate() {
  if (!taskStore.currentTask) return '// Напишите ваш код здесь\nconsole.log("Hello, World!");'
  
  switch (taskStore.currentTask.id) {
    case '1': // Hello World
      return 'console.log("Hello, World!");'
    case '2': // Сумма двух чисел
      return `function sum(a, b) {
    // Ваш код здесь
}

// Тестирование
console.log(sum(2, 3));`
    case '3': // Факториал
      return `function factorial(n) {
    // Ваш код здесь
}

// Тестирование
console.log(factorial(5));`
    default:
      return '// Напишите ваш код здесь\nconsole.log("Hello, World!");'
  }
}

function getJavaTemplate() {
  if (!taskStore.currentTask) return getDefaultJavaTemplate()
  
  switch (taskStore.currentTask.id) {
    case '1': // Hello World
      return `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}`
    case '2': // Сумма двух чисел
      return `public class Main {
    public static int sum(int a, int b) {
        return a + b;  // УБРАЛИ КОММЕНТАРИЙ И ДОБАВИЛИ RETURN
    }
    
    public static void main(String[] args) {
        System.out.println(sum(2, 3));
    }
}`
    case '3': // Факториал
      return `public class Main {
    public static int factorial(int n) {
        if (n == 0) return 1;
        int result = 1;
        for (int i = 1; i <= n; i++) {
            result *= i;
        }
        return result;  // ДОБАВИЛИ RETURN
    }
    
    public static void main(String[] args) {
        System.out.println(factorial(5));
    }
}`
    default:
      return getDefaultJavaTemplate()
  }
}

function getCppTemplate() {
  if (!taskStore.currentTask) return getDefaultCppTemplate()
  
  switch (taskStore.currentTask.id) {
    case '1': // Hello World
      return `#include <iostream>
using namespace std;

int main() {
    cout << "Hello, World!" << endl;
    return 0;
}`
    case '2': // Сумма двух чисел
      return `#include <iostream>
using namespace std;

int sum(int a, int b) {
    // Ваш код здесь
    return 0;
}

int main() {
    cout << sum(2, 3) << endl;
    return 0;
}`
    case '3': // Факториал
      return `#include <iostream>
using namespace std;

int factorial(int n) {
    // Ваш код здесь
    return 0;
}

int main() {
    cout << factorial(5) << endl;
    return 0;
}`
    default:
      return getDefaultCppTemplate()
  }
}

function getGoTemplate() {
  if (!taskStore.currentTask) return getDefaultGoTemplate()
  
  switch (taskStore.currentTask.id) {
    case '1': // Hello World
      return `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`
    case '2': // Сумма двух чисел
      return `package main

import "fmt"

func sum(a, b int) int {
    // Ваш код здесь
    return 0
}

func main() {
    fmt.Println(sum(2, 3))
}`
    case '3': // Факториал
      return `package main

import "fmt"

func factorial(n int) int {
    // Ваш код здесь
    return 0
}

func main() {
    fmt.Println(factorial(5))
}`
    default:
      return getDefaultGoTemplate()
  }
}

function getDefaultJavaTemplate() {
  return `public class Main {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
    }
}`
}

function getDefaultCppTemplate() {
  return `#include <iostream>
using namespace std;

int main() {
    std::cout << "Hello, World!" << std::endl;
    return 0;
}`
}

function getDefaultGoTemplate() {
  return `package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}`
}

function execute() {
  taskStore.executeCode()
}

// Горячие клавиши
function handleKeydown(event) {
  if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
    event.preventDefault()
    execute()
  }
  
  // Автоматическое добавление отступов
  if (event.key === 'Tab') {
    event.preventDefault()
    const start = event.target.selectionStart
    const end = event.target.selectionEnd
    
    // Вставляем 4 пробела вместо таба
    code.value = code.value.substring(0, start) + '    ' + code.value.substring(end)
    
    // Устанавливаем курсор после отступа
    nextTick(() => {
      event.target.selectionStart = event.target.selectionEnd = start + 4
    })
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
.code-editor {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #1e1e1e;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #333;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

/* НОВАЯ ШАПКА */
.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  background: linear-gradient(135deg, #2d2d30 0%, #252526 100%);
  border-bottom: 1px solid #3e3e42;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 24px;
}

.language-selector {
  display: flex;
  align-items: center;
  gap: 8px;
}

.language-selector label {
  color: #cccccc;
  font-size: 14px;
  font-weight: 500;
}

.language-selector select {
  padding: 8px 12px;
  background: #3c3c3c;
  color: white;
  border: 1px solid #565656;
  border-radius: 6px;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s ease;
  min-width: 140px;
}

.language-selector select:hover {
  border-color: #007acc;
  background: #464647;
}

.language-selector select:focus {
  outline: none;
  border-color: #007acc;
  box-shadow: 0 0 0 2px rgba(0, 122, 204, 0.2);
}

.status-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: #2a2a2a;
  border-radius: 6px;
  border: 1px solid #404040;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #4ec9b0;
  animation: pulse 2s infinite;
}

.status-indicator.running .status-dot {
  background: #ffa500;
  animation: pulse 1s infinite;
}

.status-indicator span {
  color: #cccccc;
  font-size: 12px;
  font-weight: 500;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* КНОПКА ВЫПОЛНЕНИЯ */
.run-button {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: linear-gradient(135deg, #28a745 0%, #20c997 100%);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 2px 10px rgba(40, 167, 69, 0.3);
}

.run-button:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 4px 15px rgba(40, 167, 69, 0.4);
  background: linear-gradient(135deg, #20c997 0%, #28a745 100%);
}

.run-button:active:not(:disabled) {
  transform: translateY(0);
}

.run-button:disabled {
  background: #5a5a5a;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.button-icon {
  font-size: 16px;
}

.shortcut {
  background: rgba(255, 255, 255, 0.2);
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 500;
}

/* РЕДАКТОР С НУМЕРАЦИЕЙ СТРОК */
.editor-container {
  flex: 1;
  min-height: 0;
  background: #1e1e1e;
}

.editor-wrapper {
  display: flex;
  height: 100%;
  background: #1e1e1e;
}

.editor-line-numbers {
  background: #252526;
  padding: 16px 12px;
  border-right: 1px solid #333;
  overflow-y: auto;
  min-width: 50px;
  text-align: right;
}

.line-number {
  color: #6e7681;
  font-size: 13px;
  font-family: 'Courier New', monospace;
  line-height: 1.5;
  user-select: none;
}

.text-editor {
  flex: 1;
  background: #1e1e1e;
  color: #d4d4d4;
  border: none;
  padding: 16px;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.5;
  resize: none;
  outline: none;
  box-sizing: border-box;
  white-space: pre;
  overflow-wrap: normal;
  overflow-x: auto;
}

.text-editor::placeholder {
  color: #6e7681;
}

.text-editor:focus {
  background: #1e1e1e;
}

/* СКРОЛЛБАР ДЛЯ РЕДАКТОРА */
.text-editor::-webkit-scrollbar {
  width: 8px;
}

.text-editor::-webkit-scrollbar-track {
  background: #252526;
}

.text-editor::-webkit-scrollbar-thumb {
  background: #424242;
  border-radius: 4px;
}

.text-editor::-webkit-scrollbar-thumb:hover {
  background: #565656;
}

/* СКРОЛЛБАР ДЛЯ НУМЕРАЦИИ СТРОК */
.editor-line-numbers::-webkit-scrollbar {
  width: 6px;
}

.editor-line-numbers::-webkit-scrollbar-track {
  background: #252526;
}

.editor-line-numbers::-webkit-scrollbar-thumb {
  background: #424242;
  border-radius: 3px;
}
</style>