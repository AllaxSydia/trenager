import { createApp } from 'vue'
import App from './App.vue'
import router from './router/index'

const app = createApp(App)

// Простой store
const store = {
  state: {
    isLoggedIn: false,
    username: '',
    userAvatar: ''
  },
  setUser(userData) {
    this.state.isLoggedIn = userData.isLoggedIn
    this.state.username = userData.username
    this.state.userAvatar = userData.userAvatar
  },
  logout() {
    this.state.isLoggedIn = false
    this.state.username = ''
    this.state.userAvatar = ''
  },
  initialize() {
    const savedUser = localStorage.getItem('user')
    if (savedUser) {
      try {
        const userData = JSON.parse(savedUser)
        this.setUser(userData)
      } catch (e) {
        localStorage.removeItem('user')
      }
    }
  }
}

// Инициализируем store
store.initialize()

// Делаем store глобальным
app.config.globalProperties.$store = store

app.use(router)
app.mount('#app')