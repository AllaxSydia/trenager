<template>
  <div ref="editorContainer" class="codemirror-editor"></div>
</template>

<script setup>
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'

// Только базовые импорты
import { defaultKeymap } from '@codemirror/commands'
import { lineNumbers } from '@codemirror/gutter'
import { EditorState } from '@codemirror/state'
import { EditorView, highlightActiveLine, keymap } from '@codemirror/view'

// Языки
import { cpp } from '@codemirror/lang-cpp'
import { java } from '@codemirror/lang-java'
import { javascript } from '@codemirror/lang-javascript'
import { python } from '@codemirror/lang-python'

// Тема
import { oneDark } from '@codemirror/theme-one-dark'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  language: {
    type: String,
    default: 'python'
  },
  theme: {
    type: String,
    default: 'dark'
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

const editorContainer = ref(null)
let editorView = null

// Маппинг языков
const languageExtensions = {
  python: python(),
  javascript: javascript(),
  cpp: cpp(),
  java: java()
}

onMounted(() => {
  nextTick(() => {
    initEditor()
  })
})

onUnmounted(() => {
  if (editorView) {
    editorView.destroy()
  }
})

watch(() => props.modelValue, (newValue) => {
  if (editorView && editorView.state.doc.toString() !== newValue) {
    const transaction = editorView.state.update({
      changes: {
        from: 0,
        to: editorView.state.doc.length,
        insert: newValue
      }
    })
    editorView.dispatch(transaction)
  }
})

watch(() => props.language, () => {
  nextTick(() => {
    if (editorView) {
      editorView.destroy()
    }
    initEditor()
  })
})

const initEditor = () => {
  if (!editorContainer.value) {
    console.error('Editor container not found')
    return
  }

  try {
    // Очищаем контейнер
    if (editorContainer.value.firstChild) {
      editorContainer.value.removeChild(editorContainer.value.firstChild)
    }

    // Базовые расширения
    const extensions = [
      lineNumbers(),
      highlightActiveLine(),
      keymap.of([...defaultKeymap]),
      EditorView.updateListener.of(update => {
        if (update.docChanged) {
          const value = update.state.doc.toString()
          emit('update:modelValue', value)
          emit('change', value)
        }
      }),
      EditorView.lineWrapping
    ]

    // Добавляем язык
    const langExtension = languageExtensions[props.language] || python()
    extensions.push(langExtension)

    // Добавляем тему
    if (props.theme === 'dark') {
      extensions.push(oneDark)
    }

    // Код по умолчанию
    const defaultCode = getDefaultCode(props.language)
    const initialCode = props.modelValue || defaultCode

    // Создаем состояние
    const state = EditorState.create({
      doc: initialCode,
      extensions
    })

    // Создаем редактор
    editorView = new EditorView({
      state,
      parent: editorContainer.value
    })

    console.log('✅ Editor initialized successfully with language:', props.language)

  } catch (error) {
    console.error('❌ Failed to initialize editor:', error)
  }
}

// Код по умолчанию для каждого языка
const getDefaultCode = (language) => {
  const defaults = {
    python: '# Write your Python code here\nprint("Hello World")',
    javascript: '// Write your JavaScript code here\nconsole.log("Hello World")',
    cpp: '// Write your C++ code here\n#include <iostream>\n\nint main() {\n    std::cout << "Hello World" << std::endl;\n    return 0;\n}',
    java: '// Write your Java code here\npublic class Main {\n    public static void main(String[] args) {\n        System.out.println("Hello World");\n    }\n}'
  }
  return defaults[language] || defaults.python
}

// Публичные методы
defineExpose({
  focus: () => {
    if (editorView) {
      editorView.focus()
    }
  },
  getValue: () => {
    return editorView ? editorView.state.doc.toString() : ''
  },
  setValue: (value) => {
    if (editorView) {
      const transaction = editorView.state.update({
        changes: {
          from: 0,
          to: editorView.state.doc.length,
          insert: value || ''
        }
      })
      editorView.dispatch(transaction)
    }
  }
})
</script>

<style scoped>
.codemirror-editor {
  width: 100%;
  height: 100%;
  min-height: 400px;
  border: 1px solid #ccc;
  border-radius: 6px;
  overflow: hidden;
}

/* Стили для CodeMirror */
:deep(.cm-editor) {
  height: 100%;
  font-size: 14px;
  font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
}

:deep(.cm-content) {
  padding: 10px 0;
  caret-color: #fff;
}

:deep(.cm-line) {
  padding: 0 10px;
}

:deep(.cm-gutters) {
  background-color: #1e1e1e;
  border-right: 1px solid #333;
  color: #858585;
}

:deep(.cm-activeLineGutter) {
  background-color: #2a2a2a;
  color: #fff;
}

:deep(.cm-activeLine) {
  background-color: #2a2a2a;
}
</style>