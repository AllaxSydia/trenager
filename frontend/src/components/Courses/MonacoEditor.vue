<template>
  <div class="monaco-editor-container">
    <div ref="editorContainer" class="editor"></div>
  </div>
</template>

<script>
import * as monaco from 'monaco-editor';

export default {
  name: 'MonacoEditor',
  props: {
    code: {
      type: String,
      default: ''
    },
    language: {
      type: String,
      default: 'python'
    },
    readOnly: {
      type: Boolean,
      default: false
    }
  },
  emits: ['update:code'],
  data() {
    return {
      editor: null
    }
  },
  watch: {
    code(newCode) {
      if (this.editor && newCode !== this.editor.getValue()) {
        this.editor.setValue(newCode)
      }
    },
    language(newLanguage) {
      if (this.editor) {
        monaco.editor.setModelLanguage(this.editor.getModel(), newLanguage)
      }
    }
  },
  mounted() {
    this.initMonaco()
  },
  beforeUnmount() {
    if (this.editor) {
      this.editor.dispose()
    }
  },
  methods: {
    initMonaco() {
      if (!this.$refs.editorContainer) return

      this.editor = monaco.editor.create(this.$refs.editorContainer, {
        value: this.code,
        language: this.language,
        theme: 'vs-dark',
        fontSize: 14,
        lineNumbers: 'on',
        roundedSelection: true,
        scrollBeyondLastLine: false,
        readOnly: this.readOnly,
        automaticLayout: true,
        minimap: { enabled: false },
        folding: true,
        wordWrap: 'on',
        tabSize: 2,
        insertSpaces: true,
        detectIndentation: true
      })

      // Слушаем изменения кода
      this.editor.onDidChangeModelContent(() => {
        this.$emit('update:code', this.editor.getValue())
      })

      // Горячие клавиши
      this.editor.addCommand(
        monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter,
        () => {
          this.$emit('execute')
        }
      )
    },

    getCode() {
      return this.editor ? this.editor.getValue() : ''
    },

    setCode(code) {
      if (this.editor) {
        this.editor.setValue(code)
      }
    }
  }
}
</script>

<style scoped>
.monaco-editor-container {
  height: 400px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #404040;
}

.editor {
  height: 100%;
  width: 100%;
}

/* Кастомные скроллбары для Monaco */
.monaco-editor-container :deep(.monaco-scrollable-element) {
  scrollbar-width: thin;
}

.monaco-editor-container :deep(.monaco-scrollable-element)::-webkit-scrollbar {
  width: 10px;
  height: 10px;
}

.monaco-editor-container :deep(.monaco-scrollable-element)::-webkit-scrollbar-track {
  background: #1e1e1e;
}

.monaco-editor-container :deep(.monaco-scrollable-element)::-webkit-scrollbar-thumb {
  background: #424242;
  border-radius: 5px;
}

.monaco-editor-container :deep(.monaco-scrollable-element)::-webkit-scrollbar-thumb:hover {
  background: #4f4f4f;
}
</style>