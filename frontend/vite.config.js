import vue from '@vitejs/plugin-vue'
import { defineConfig } from 'vite'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': '/src'
    }
  },
  server: {
    host: true, // Для Docker
    port: 5173
  },
  build: {
    chunkSizeWarningLimit: 1000
  }
})