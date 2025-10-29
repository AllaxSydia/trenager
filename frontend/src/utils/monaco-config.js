// Monaco configuration for Vite
import * as monaco from 'monaco-editor'

// Configure Monaco environment for workers
self.MonacoEnvironment = {
  getWorkerUrl: function (moduleId, label) {
    if (label === 'json') {
      return new URL('monaco-editor/esm/vs/language/json/json.worker.js', import.meta.url).toString()
    }
    if (label === 'css' || label === 'scss' || label === 'less') {
      return new URL('monaco-editor/esm/vs/language/css/css.worker.js', import.meta.url).toString()
    }
    if (label === 'html' || label === 'handlebars' || label === 'razor') {
      return new URL('monaco-editor/esm/vs/language/html/html.worker.js', import.meta.url).toString()
    }
    if (label === 'typescript' || label === 'javascript') {
      return new URL('monaco-editor/esm/vs/language/typescript/ts.worker.js', import.meta.url).toString()
    }
    return new URL('monaco-editor/esm/vs/editor/editor.worker.js', import.meta.url).toString()
  }
}

export { monaco }
