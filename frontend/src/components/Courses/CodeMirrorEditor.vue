<template>
  <div class="code-mirror-editor">
    <div class="editor-header">
      <span class="language-badge">{{ language }}</span>
      <span class="file-name">{{ getFileName() }}</span>
    </div>
    <div ref="editorElement" class="editor-container"></div>
  </div>
</template>

<script>
import { autocompletion, completionKeymap } from '@codemirror/autocomplete'
import { defaultKeymap, indentWithTab } from '@codemirror/commands'
import { cpp } from '@codemirror/lang-cpp'
import { java } from '@codemirror/lang-java'
import { javascript } from '@codemirror/lang-javascript'
import { python } from '@codemirror/lang-python'
import { lintKeymap } from '@codemirror/lint'
import { EditorState } from '@codemirror/state'
import { oneDark } from '@codemirror/theme-one-dark'
import { EditorView, keymap } from '@codemirror/view'

export default {
  name: 'CodeMirrorEditor',
  props: {
    language: {
      type: String,
      required: true,
      validator: (value) => ['python', 'javascript', 'cpp', 'java'].includes(value)
    },
    modelValue: {
      type: String,
      default: ''
    }
  },
  emits: ['update:modelValue'],
  data() {
    return {
      editor: null,
      code: this.modelValue
    }
  },
  watch: {
    modelValue(newVal) {
      if (this.editor && newVal !== this.editor.state.doc.toString()) {
        this.editor.dispatch({
          changes: {
            from: 0,
            to: this.editor.state.doc.length,
            insert: newVal
          }
        })
      }
    },
    language(newLang) {
      if (this.editor) {
        this.initEditor()
      }
    }
  },
  mounted() {
    this.initEditor()
  },
  beforeUnmount() {
    if (this.editor) {
      this.editor.destroy()
    }
  },
  methods: {
    initEditor() {
      if (this.editor) {
        this.editor.destroy()
      }

      const languageExtensions = {
        python: python(),
        javascript: javascript(),
        cpp: cpp(),
        java: java()
      }

      const startState = EditorState.create({
        doc: this.modelValue, // Используем modelValue вместо this.code
        extensions: [
          keymap.of([...defaultKeymap, indentWithTab, ...completionKeymap, ...lintKeymap]),
          languageExtensions[this.language] || python(),
          oneDark,
          autocompletion(),
          EditorView.updateListener.of((update) => {
            if (update.docChanged) {
              const code = update.state.doc.toString()
              this.code = code
              this.$emit('update:modelValue', code)
            }
          }),
          EditorView.theme({
            "&": {
              height: "100%",
              fontSize: "14px",
              fontFamily: "'Monaco', 'Menlo', 'Ubuntu Mono', monospace"
            },
            ".cm-content": {
              fontFamily: "'Monaco', 'Menlo', 'Ubuntu Mono', monospace",
              padding: "8px 0"
            },
            ".cm-line": {
              padding: "0 8px"
            },
            ".cm-gutters": {
              backgroundColor: "#1a1a1a",
              color: "#666",
              borderRight: "1px solid #333"
            },
            ".cm-activeLine": {
              backgroundColor: "#2a2a2a"
            },
            ".cm-selectionBackground": {
              background: "#3a3a3a"
            }
          })
        ]
      })

      this.editor = new EditorView({
        state: startState,
        parent: this.$refs.editorElement
      })

      // Авто-фокус
      this.editor.focus()
    },

    getFileName() {
      const extensions = {
        python: 'main.py',
        javascript: 'script.js',
        cpp: 'main.cpp',
        java: 'Main.java'
      }
      return extensions[this.language] || 'code.txt'
    },

    getCode() {
      return this.editor ? this.editor.state.doc.toString() : ''
    },

    setCode(newCode) {
      if (this.editor) {
        this.editor.dispatch({
          changes: {
            from: 0,
            to: this.editor.state.doc.length,
            insert: newCode
          }
        })
      }
    },

    focus() {
      if (this.editor) {
        this.editor.focus()
      }
    }
  }
}
</script>

<style scoped>
.code-mirror-editor {
  border: 1px solid #374151;
  border-radius: 8px;
  overflow: hidden;
  background: #1f2937;
  height: 100%;
}

.editor-header {
  background: #111827;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid #374151;
}

.language-badge {
  background: #3b82f6;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
}

.file-name {
  color: #9ca3af;
  font-size: 13px;
  font-family: 'Monaco', 'Menlo', monospace;
}

.editor-container {
  background: #1e1e1e;
  height: calc(100% - 49px); /* Высота минус заголовок */
}

/* Стили для CodeMirror */
:deep(.cm-editor) {
  height: 100%;
  outline: none;
}

:deep(.cm-content) {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}

:deep(.cm-tooltip) {
  background: #2d2d2d;
  border: 1px solid #444;
  border-radius: 4px;
}

:deep(.cm-completionLabel) {
  font-family: 'Monaco', 'Menlo', monospace;
}
</style>