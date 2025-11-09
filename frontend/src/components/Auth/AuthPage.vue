<template>
  <div class="auth-page">
    <main class="main">
      <div class="container">
        <div class="auth-container">
          <transition name="card-scale" mode="out-in">
            <div :key="activeTab" class="auth-card">
              <h1 class="auth-title">CodeForge</h1>
              
              <div class="auth-tabs">
                <button 
                  class="auth-tab" 
                  :class="{ 'auth-tab--active': activeTab === 'login' }"
                  @click="switchTab('login')"
                >
                  Вход
                </button>
                <button 
                  class="auth-tab" 
                  :class="{ 'auth-tab--active': activeTab === 'register' }"
                  @click="switchTab('register')"
                >
                  Регистрация
                </button>
              </div>
              
              <div class="forms-wrapper">
                <!-- Форма входа -->
                <transition name="scale-fade" mode="out-in">
                  <form 
                    v-if="activeTab === 'login'" 
                    key="login"
                    class="auth-form" 
                    @submit.prevent="handleLogin"
                  >
                    <div class="form-group">
                      <label for="login-email" class="form-label">Email:</label>
                      <input 
                        type="email" 
                        id="login-email" 
                        class="form-input" 
                        placeholder="Введите почту"
                        v-model="loginData.email"
                        required
                      >
                    </div>
                    
                    <div class="form-group">
                      <label for="login-password" class="form-label">Пароль:</label>
                      <input 
                        type="password" 
                        id="login-password" 
                        class="form-input" 
                        placeholder="Введите пароль"
                        v-model="loginData.password"
                        required
                      >
                    </div>
                    
                    <div class="form-options">
                      <a href="#" class="forgot-password">Забыли пароль?</a>
                    </div>
                    
                    <button class="auth-submit-btn" type="submit" :disabled="isLoading">
                      {{ isLoading ? 'Вход...' : 'Вход' }}
                    </button>
                  </form>
                  
                  <!-- Форма регистрации -->
                  <form 
                    v-else 
                    key="register"
                    class="auth-form" 
                    @submit.prevent="handleRegister"
                  >
                    <div class="form-group">
                      <label for="register-username" class="form-label">Имя пользователя:</label>
                      <input 
                        type="text" 
                        id="register-username" 
                        class="form-input" 
                        placeholder="Ваше имя"
                        v-model="registerData.username"
                        required
                      >
                    </div>
                    
                    <div class="form-group">
                      <label for="register-email" class="form-label">Email:</label>
                      <input 
                        type="email" 
                        id="register-email" 
                        class="form-input" 
                        placeholder="Введите почту"
                        v-model="registerData.email"
                        required
                      >
                    </div>
                    
                    <div class="form-group">
                      <label for="register-password" class="form-label">Пароль:</label>
                      <input 
                        type="password" 
                        id="register-password" 
                        class="form-input" 
                        placeholder="Придумайте пароль"
                        v-model="registerData.password"
                        required
                      >
                    </div>
                    
                    <div class="form-group">
                      <label for="register-confirm-password" class="form-label">Подтвердите пароль:</label>
                      <input 
                        type="password" 
                        id="register-confirm-password" 
                        class="form-input" 
                        placeholder="Повторите пароль"
                        v-model="registerData.confirmPassword"
                        required
                      >
                    </div>
                    
                    <div class="form-agreement">
                      <label class="checkbox-label">
                        <input 
                          type="checkbox" 
                          v-model="registerData.agreeToTerms"
                        >
                        <span class="checkmark"></span>
                        Я согласен с <a href="#" class="terms-link">условиями использования</a>
                      </label>
                    </div>
                    
                    <button class="auth-submit-btn" type="submit" :disabled="isLoading">
                      {{ isLoading ? 'Регистрация...' : 'Зарегистрироваться' }}
                    </button>
                  </form>
                </transition>
              </div>
            </div>
          </transition>
        </div>
      </div>
    </main>

    <!-- Toast уведомления -->
    <transition name="toast-slide">
      <div v-if="toast.show" :class="['toast', `toast--${toast.type}`]" @click="hideToast">
        <div class="toast-icon">
          <span v-if="toast.type === 'error'">⚠️</span>
          <span v-else-if="toast.type === 'success'">✅</span>
          <span v-else>ℹ️</span>
        </div>
        <div class="toast-content">
          <div class="toast-title">{{ toast.title }}</div>
          <div class="toast-message">{{ toast.message }}</div>
        </div>
        <button class="toast-close" @click.stop="hideToast">×</button>
      </div>
    </transition>
  </div>
</template>

<script>
export default {
  name: 'AuthPage',
  data() {
    return {
      activeTab: 'login',
      isLoading: false,
      loginData: {
        email: '',
        password: ''
      },
      registerData: {
        username: '',
        email: '',
        password: '',
        confirmPassword: '',
        agreeToTerms: false
      },
      toast: {
        show: false,
        type: 'error',
        title: '',
        message: '',
        timeout: null
      }
    }
  },
  methods: {
    showToast(type, title, message, duration = 4000) {
      console.log('Showing toast:', type, title, message) // Для отладки
      
      // Очищаем предыдущий таймер
      if (this.toast.timeout) {
        clearTimeout(this.toast.timeout)
      }
      
      this.toast.show = false // Сначала скрываем, чтобы анимация работала
      
      this.$nextTick(() => {
        this.toast = {
          show: true,
          type,
          title,
          message,
          timeout: setTimeout(() => {
            this.hideToast()
          }, duration)
        }
      })
    },
    
    hideToast() {
      console.log('Hiding toast') // Для отладки
      if (this.toast.timeout) {
        clearTimeout(this.toast.timeout)
        this.toast.timeout = null
      }
      this.toast.show = false
    },
    
    switchTab(tab) {
      if (this.activeTab === tab) return
      this.activeTab = tab
      this.resetForms()
      this.hideToast()
    },
    
    async handleLogin() {
      if (this.isLoading) return
      
      this.isLoading = true
      
      try {
        // Простая валидация
        if (!this.loginData.email || !this.loginData.password) {
          this.showToast('error', 'Ошибка', 'Заполните все поля')
          this.isLoading = false
          return
        }
        
        const userData = {
          isLoggedIn: true,
          username: this.loginData.email.split('@')[0] || 'User',
          userAvatar: (this.loginData.email.split('@')[0] || 'U').charAt(0).toUpperCase()
        }
        
        localStorage.setItem('user', JSON.stringify(userData))
        
        this.showToast('success', 'Успешный вход!', 'Добро пожаловать в CodeForge!', 2000)
        
        setTimeout(() => {
          this.$router.push('/')
        }, 1500)
        
      } catch (error) {
        console.error('Ошибка входа:', error)
        this.showToast('error', 'Ошибка входа', 'Проверьте правильность данных')
      } finally {
        this.isLoading = false
      }
    },
    
    async handleRegister() {
      if (this.isLoading) return
      
      // Валидация
      if (!this.registerData.username || !this.registerData.email || !this.registerData.password) {
        this.showToast('error', 'Ошибка', 'Заполните все обязательные поля')
        return
      }
      
      if (this.registerData.password !== this.registerData.confirmPassword) {
        this.showToast('error', 'Ошибка', 'Пароли не совпадают')
        return
      }
      
      if (this.registerData.password.length < 6) {
        this.showToast('error', 'Ошибка', 'Пароль должен содержать минимум 6 символов')
        return
      }
      
      if (!this.registerData.agreeToTerms) {
        this.showToast('error', 'Внимание', 'Необходимо согласиться с условиями использования')
        return
      }
      
      this.isLoading = true
      
      try {
        await new Promise(resolve => setTimeout(resolve, 1000))
        
        const userData = {
          isLoggedIn: true,
          username: this.registerData.username,
          userAvatar: this.registerData.username.charAt(0).toUpperCase()
        }
        
        localStorage.setItem('user', JSON.stringify(userData))
        
        this.showToast('success', 'Регистрация успешна!', 'Добро пожаловать в CodeForge!', 2000)
        
        setTimeout(() => {
          this.$router.push('/')
        }, 1500)
        
      } catch (error) {
        console.error('Ошибка регистрации:', error)
        this.showToast('error', 'Ошибка регистрации', 'Попробуйте еще раз')
      } finally {
        this.isLoading = false
      }
    },
    
    resetForms() {
      this.loginData = {
        email: '',
        password: ''
      }
      this.registerData = {
        username: '',
        email: '',
        password: '',
        confirmPassword: '',
        agreeToTerms: false
      }
    }
  },
  
  mounted() {
    this.resetForms()
  },
  
  beforeUnmount() {
    if (this.toast.timeout) {
      clearTimeout(this.toast.timeout)
    }
  }
}
</script>

<style scoped>
/* Основные стили для полного экрана */
.auth-page {
  height: 100vh;
  width: 100vw;
  background-color: #0E1117;
  overflow: hidden;
  position: fixed;
  top: 0;
  left: 0;
  margin: 0;
  padding: 0;
}

.container {
  height: 100%;
  width: 100%;
  margin: 0;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.main {
  height: 100%;
  width: 100%;
  padding: 0;
  margin: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.auth-container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
  padding: 0;
}

.auth-card {
  background: #1F2937;
  border: 1px solid #374151;
  border-radius: 12px;
  padding: 2.5rem;
  width: 90%;
  max-width: 400px;
  max-height: 90vh;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  transform-origin: center center;
  overflow-y: auto;
}

/* Скрываем скроллбар для auth-card */
.auth-card::-webkit-scrollbar {
  width: 0;
  display: none;
}

.auth-card {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.auth-title {
  color: #E5E7EB;
  font-size: 2rem;
  font-weight: bold;
  text-align: center;
  margin-bottom: 2rem;
  font-family: monospace;
}

.auth-tabs {
  display: flex;
  background: #374151;
  border-radius: 8px;
  padding: 4px;
  margin-bottom: 2rem;
}

.auth-tab {
  flex: 1;
  background: transparent;
  border: none;
  color: #9CA3AF;
  padding: 12px 16px;
  border-radius: 6px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  font-size: 16px;
}

.auth-tab:hover {
  color: #E5E7EB;
  background: rgba(75, 85, 99, 0.5);
}

.auth-tab--active {
  color: #E5E7EB;
  background: #3B82F6;
}

.auth-tab:active {
  transform: scale(0.98);
}

.forms-wrapper {
  position: relative;
  min-height: 300px;
  overflow: hidden;
}

.form-container {
  width: 100%;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
  width: 100%;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-label {
  color: #E5E7EB;
  font-weight: 500;
  font-size: 16px;
}

.form-input {
  background: #374151;
  border: 1px solid #4B5563;
  border-radius: 8px;
  padding: 12px 16px;
  color: #E5E7EB;
  font-size: 16px;
  transition: all 0.3s ease;
}

.form-input:focus {
  outline: none;
  border-color: #3B82F6;
  background: #4B5563;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.15);
}

.form-input::placeholder {
  color: #9CA3AF;
}

.form-options {
  display: flex;
  justify-content: flex-end;
}

.forgot-password {
  color: #3B82F6;
  text-decoration: none;
  font-size: 14px;
  transition: all 0.3s ease;
}

.forgot-password:hover {
  color: #60A5FA;
  text-decoration: underline;
}

.form-agreement {
  margin: 1rem 0;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #9CA3AF;
  font-size: 14px;
  cursor: pointer;
  transition: color 0.3s ease;
}

.checkbox-label:hover {
  color: #E5E7EB;
}

.checkbox-label input[type="checkbox"] {
  display: none;
}

.checkmark {
  width: 18px;
  height: 18px;
  border: 2px solid #4B5563;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  flex-shrink: 0;
}

.checkbox-label input[type="checkbox"]:checked + .checkmark {
  background: #3B82F6;
  border-color: #3B82F6;
  transform: scale(1.05);
}

.checkbox-label input[type="checkbox"]:checked + .checkmark::after {
  content: '✓';
  color: white;
  font-size: 12px;
  font-weight: bold;
}

.terms-link {
  color: #3B82F6;
  text-decoration: none;
  transition: color 0.3s ease;
}

.terms-link:hover {
  text-decoration: underline;
  color: #60A5FA;
}

.auth-submit-btn {
  background: #3B82F6;
  color: #E5E7EB;
  border: none;
  border-radius: 8px;
  padding: 14px 20px;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  margin-top: 1rem;
}

.auth-submit-btn:hover:not(:disabled) {
  background: #2563EB;
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(37, 99, 235, 0.3);
}

.auth-submit-btn:active:not(:disabled) {
  transform: translateY(0);
}

.auth-submit-btn:disabled {
  background: #6B7280;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

/* Анимация увеличения карточки */
.card-scale-enter-active {
  transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.card-scale-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.6, 1);
}

.card-scale-enter-from {
  opacity: 0;
  transform: scale(0.9);
}

.card-scale-leave-to {
  opacity: 0;
  transform: scale(1.1);
}

.card-scale-enter-to,
.card-scale-leave-from {
  opacity: 1;
  transform: scale(1);
}

/* Анимация увеличения блока форм */
.scale-fade-enter-active {
  transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.scale-fade-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.6, 1);
}

.scale-fade-enter-from {
  opacity: 0;
  transform: scale(0.95);
}

.scale-fade-leave-to {
  opacity: 0;
  transform: scale(1.05);
}

.scale-fade-enter-to,
.scale-fade-leave-from {
  opacity: 1;
  transform: scale(1);
}

/* Стили для Toast уведомлений */
.toast {
  position: fixed;
  top: 20px;
  right: 20px;
  background: #1F2937;
  border: 1px solid #374151;
  border-radius: 12px;
  padding: 16px;
  min-width: 300px;
  max-width: 400px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: flex-start;
  gap: 12px;
  z-index: 10000;
  backdrop-filter: blur(10px);
  cursor: pointer;
  transition: all 0.3s ease;
}

.toast:hover {
  transform: translateY(-2px);
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.4);
}

.toast--error {
  border-left: 4px solid #EF4444;
}

.toast--success {
  border-left: 4px solid #10B981;
}

.toast--info {
  border-left: 4px solid #3B82F6;
}

.toast-icon {
  font-size: 20px;
  flex-shrink: 0;
  margin-top: 2px;
}

.toast-content {
  flex: 1;
  min-width: 0;
}

.toast-title {
  color: #E5E7EB;
  font-weight: 600;
  font-size: 16px;
  margin-bottom: 4px;
}

.toast-message {
  color: #9CA3AF;
  font-size: 14px;
  line-height: 1.4;
}

.toast-close {
  background: none;
  border: none;
  color: #9CA3AF;
  font-size: 20px;
  cursor: pointer;
  padding: 0;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 4px;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.toast-close:hover {
  background: #374151;
  color: #E5E7EB;
}

/* Анимация Toast */
.toast-slide-enter-active {
  transition: all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.toast-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 1, 1);
}

.toast-slide-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-slide-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.toast-slide-enter-to {
  opacity: 1;
  transform: translateX(0);
}

/* Адаптация для мобильных устройств */
@media (max-width: 480px) {
  .auth-card {
    padding: 2rem 1.5rem;
    width: 95%;
    max-height: 95vh;
  }
  
  .auth-title {
    font-size: 1.75rem;
    margin-bottom: 1.5rem;
  }
  
  .toast {
    top: 10px;
    right: 10px;
    left: 10px;
    min-width: auto;
    max-width: none;
  }
}

/* Предотвращаем масштабирование на мобильных устройствах при фокусе на инпуты */
@media (max-width: 768px) {
  .form-input {
    font-size: 16px; /* Предотвращает зум на iOS */
  }
}
</style>