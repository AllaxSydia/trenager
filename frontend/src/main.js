import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

console.log('🚀 Vue app starting...')

const app = createApp(App)
app.use(router)

// Добавьте обработчик ошибок
app.config.errorHandler = (err, instance, info) => {
  console.error('Vue error:', err)
  console.log('Error info:', info)
}

app.mount('#app')

console.log('✅ Vue app mounted')