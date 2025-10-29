<template>
  <div ref="editorContainer" class="monaco-editor"></div>
</template>

<script setup>
import { monaco } from '@/utils/monaco-config'
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  language: {
    type: String,
    default: 'javascript'
  },
  theme: {
    type: String,
    default: 'vs-dark'
  },
  readOnly: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const editorContainer = ref(null)
let editor = null

// Маппинг языков
const languageMap = {
  python: 'python',
  javascript: 'javascript',
  cpp: 'cpp',
  java: 'java'
}

onMounted(async () => {
  await nextTick()
  initEditor()
})

onUnmounted(() => {
  if (editor) {
    editor.dispose()
  }
})

watch(() => props.modelValue, (newValue) => {
  if (editor && editor.getValue() !== newValue) {
    editor.setValue(newValue)
  }
})

watch(() => props.language, (newLang) => {
  if (editor) {
    const model = editor.getModel()
    monaco.editor.setModelLanguage(model, languageMap[newLang] || 'javascript')
  }
})

watch(() => props.theme, (newTheme) => {
  if (editor) {
    monaco.editor.setTheme(newTheme)
  }
})

const initEditor = () => {
  if (!editorContainer.value) {
    console.error('Editor container not found')
    return
  }

  try {
    console.log('🔄 Initializing Monaco Editor...')

    // Очищаем контейнер
    editorContainer.value.innerHTML = ''

    // Создаем редактор
    editor = monaco.editor.create(editorContainer.value, {
      value: props.modelValue || getDefaultCode(props.language),
      language: languageMap[props.language] || 'javascript',
      theme: props.theme,
      fontSize: 14,
      fontFamily: "'Cascadia Code', 'Consolas', 'Courier New', monospace",
      lineNumbers: 'on',
      roundedSelection: true,
      scrollBeyondLastLine: false,
      readOnly: props.readOnly,
      automaticLayout: true,
      tabSize: 2,
      insertSpaces: true,
      detectIndentation: true,
      minimap: { enabled: false },
      bracketPairColorization: { enabled: true },
      guides: { indentation: true },
      wordWrap: 'on'
    })

    // Обработчик изменений
    editor.onDidChangeModelContent(() => {
      const value = editor.getValue()
      emit('update:modelValue', value)
      emit('change', value)
    })

    console.log('✅ Monaco Editor initialized successfully')

  } catch (error) {
    console.error('❌ Failed to initialize Monaco Editor:', error)
    createFallbackEditor()
  }
}

const createFallbackEditor = () => {
  if (!editorContainer.value) return
  
  console.log('🔄 Creating fallback textarea editor...')
  
  const textarea = document.createElement('textarea')
  textarea.value = props.modelValue || getDefaultCode(props.language)
  textarea.style.width = '100%'
  textarea.style.height = '400px'
  textarea.style.fontFamily = 'Consolas, Monaco, monospace'
  textarea.style.fontSize = '14px'
  textarea.style.padding = '10px'
  textarea.style.border = '1px solid #ccc'
  textarea.style.borderRadius = '4px'
  textarea.style.backgroundColor = props.theme === 'vs-dark' ? '#1e1e1e' : '#ffffff'
  textarea.style.color = props.theme === 'vs-dark' ? '#d4d4d4' : '#000000'
  
  textarea.addEventListener('input', (e) => {
    emit('update:modelValue', e.target.value)
    emit('change', e.target.value)
  })
  
  editorContainer.value.appendChild(textarea)
}

const getDefaultCode = (language) => {
  const defaults = {
    python: '# Write your Python code here\nprint("Hello World")',
    javascript: '// Write your JavaScript code here\nconsole.log("Hello World")',
    cpp: '// Write your C++ code here\n#include <iostream>\n\nint main() {\n    std::cout << "Hello World" << std::endl;\n    return 0;\n}',
    java: '// Write your Java code here\npublic class Main {\n    public static void main(String[] args) {\n        System.out.println("Hello World");\n    }\n}'
  }
  return defaults[language] || defaults.javascript
}

// Публичные методы
defineExpose({
  focus: () => {
    if (editor) {
      editor.focus()
    }
  },
  getValue: () => {
    return editor ? editor.getValue() : ''
  },
  setValue: (value) => {
    if (editor) {
      editor.setValue(value || '')
    }
  }
})
</script>

<style scoped>
.monaco-editor {
  width: 100%;
  height: 100%;
  min-height: 400px;
  border: 1px solid #444;
  border-radius: 6px;
  overflow: hidden;
}
</style>