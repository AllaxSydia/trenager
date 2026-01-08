<template>
  <div class="settings-page">
    <div class="settings-container">
      <!-- Хедер -->
      <div class="settings-header">
        <h1>Настройки</h1>
        <p>Управление вашим аккаунтом и настройками платформы</p>
      </div>

      <!-- Основной контент -->
      <div class="settings-content">
        <!-- Левая колонка - настройки аккаунта -->
        <div class="settings-section">
          <div class="section-card">
            <h2>Настройки аккаунта</h2>
            
            <div class="form-group">
              <label class="form-label">Имя пользователя</label>
              <input 
                v-model="user.username" 
                type="text" 
                placeholder="Введите имя пользователя"
                class="form-input"
              >
            </div>
            
            <div class="form-group">
              <label class="form-label">Email</label>
              <input 
                v-model="user.email" 
                type="email" 
                placeholder="Введите email"
                class="form-input"
              >
            </div>
            
            <div class="form-group">
              <label class="form-label">Текущая роль</label>
              <div class="role-display">
                <span :class="['role-badge', user.role]">
                  {{ getRoleName(user.role) }}
                </span>
                <span class="role-description">
                  {{ getRoleDescription(user.role) }}
                </span>
              </div>
            </div>
            
            <div class="form-actions">
              <button @click="saveProfile" class="btn btn-primary">
                Сохранить изменения
              </button>
              <button @click="logout" class="btn btn-outline">
                Выйти
              </button>
            </div>
          </div>
        </div>

        <!-- Правая колонка - настройки платформы -->
        <div class="settings-section">
          <div class="section-card">
            <h2>Настройки интерфейса</h2>
            
            <div class="form-group">
              <label class="form-label">Тема оформления</label>
              <div class="theme-selector">
                <button 
                  v-for="theme in themes" 
                  :key="theme.id"
                  @click="selectTheme(theme.id)"
                  :class="['theme-btn', { 'theme-btn--active': currentTheme === theme.id }]"
                  :style="{ background: theme.bg, color: theme.color }"
                >
                  {{ theme.name }}
                </button>
              </div>
            </div>
            
            <div class="form-group">
              <label class="form-label">Язык интерфейса</label>
              <select v-model="language" class="form-select">
                <option value="ru">Русский</option>
                <option value="en">English</option>
              </select>
            </div>
            
            <div class="form-group">
              <label class="form-label">Уведомления</label>
              <div class="checkbox-group">
                <label class="checkbox-label">
                  <input 
                    v-model="notifications.email" 
                    type="checkbox" 
                    class="checkbox"
                  >
                  <span>Email уведомления</span>
                </label>
                <label class="checkbox-label">
                  <input 
                    v-model="notifications.browser" 
                    type="checkbox" 
                    class="checkbox"
                  >
                  <span>Браузерные уведомления</span>
                </label>
                <label class="checkbox-label">
                  <input 
                    v-model="notifications.achievements" 
                    type="checkbox" 
                    class="checkbox"
                  >
                  <span>Уведомления о достижениях</span>
                </label>
              </div>
            </div>
            
            <div class="form-group">
              <label class="form-label">Автосохранение кода</label>
              <div class="checkbox-group">
                <label class="checkbox-label">
                  <input 
                    v-model="autoSave" 
                    type="checkbox" 
                    class="checkbox"
                  >
                  <span>Сохранять код автоматически</span>
                </label>
                <div v-if="autoSave" class="auto-save-options">
                  <label class="radio-label">
                    <input 
                      v-model="autoSaveInterval" 
                      type="radio" 
                      value="30" 
                      class="radio"
                    >
                    <span>Каждые 30 секунд</span>
                  </label>
                  <label class="radio-label">
                    <input 
                      v-model="autoSaveInterval" 
                      type="radio" 
                      value="60" 
                      class="radio"
                    >
                    <span>Каждую минуту</span>
                  </label>
                  <label class="radio-label">
                    <input 
                      v-model="autoSaveInterval" 
                      type="radio" 
                      value="300" 
                      class="radio"
                    >
                    <span>Каждые 5 минут</span>
                  </label>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Статистика -->
          <div class="section-card">
            <h2>Статистика</h2>
            
            <div class="stats-grid">
              <div class="stat-item">
                <span class="stat-label">Решено задач</span>
                <span class="stat-value">{{ stats.solvedTasks }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">Пройдено тестов</span>
                <span class="stat-value">{{ stats.passedTests }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">Потрачено времени</span>
                <span class="stat-value">{{ formatTime(stats.timeSpent) }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">Рейтинг</span>
                <span class="stat-value">{{ stats.rating }}</span>
              </div>
            </div>
            
            <button @click="resetStats" class="btn btn-sm btn-danger">
              Сбросить статистику
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Settings',
  
  data() {
    return {
      user: {
        username: '',
        email: '',
        role: 'student'
      },
      
      themes: [
        { id: 'dark', name: 'Тёмная', bg: '#1F2937', color: '#F9FAFB' },
        { id: 'light', name: 'Светлая', bg: '#F9FAFB', color: '#1F2937' },
        { id: 'blue', name: 'Синяя', bg: '#1E40AF', color: '#F9FAFB' },
        { id: 'purple', name: 'Фиолетовая', bg: '#6D28D9', color: '#F9FAFB' }
      ],
      currentTheme: 'dark',
      
      language: 'ru',
      
      notifications: {
        email: true,
        browser: false,
        achievements: true
      },
      
      autoSave: true,
      autoSaveInterval: '30',
      
      stats: {
        solvedTasks: 0,
        passedTests: 0,
        timeSpent: 0, // в секундах
        rating: 0
      }
    }
  },
  
  mounted() {
    this.loadUserData()
    this.loadSettings()
    this.loadStats()
  },
  
  methods: {
    loadUserData() {
      const savedUser = JSON.parse(localStorage.getItem('user') || '{}')
      this.user = {
        username: savedUser.username || 'Пользователь',
        email: savedUser.email || '',
        role: savedUser.role || 'student'
      }
    },
    
    loadSettings() {
      const settings = JSON.parse(localStorage.getItem('settings') || '{}')
      this.currentTheme = settings.theme || 'dark'
      this.language = settings.language || 'ru'
      this.notifications = settings.notifications || this.notifications
      this.autoSave = settings.autoSave !== false
      this.autoSaveInterval = settings.autoSaveInterval || '30'
      
      // Применяем тему
      this.applyTheme(this.currentTheme)
    },
    
    loadStats() {
      const stats = JSON.parse(localStorage.getItem('stats') || '{}')
      this.stats = { ...this.stats, ...stats }
    },
    
    applyTheme(themeId) {
      document.documentElement.setAttribute('data-theme', themeId)
      localStorage.setItem('theme', themeId)
    },
    
    selectTheme(themeId) {
      this.currentTheme = themeId
      this.applyTheme(themeId)
      this.saveSettings()
    },
    
    getRoleName(role) {
      const roles = {
        student: 'Студент',
        teacher: 'Учитель',
        admin: 'Администратор'
      }
      return roles[role] || role
    },
    
    getRoleDescription(role) {
      const descriptions = {
        student: 'Можно проходить курсы и решать задачи',
        teacher: 'Можно создавать задачи и просматривать статистику',
        admin: 'Полный доступ ко всем функциям платформы'
      }
      return descriptions[role] || 'Неизвестная роль'
    },
    
    formatTime(seconds) {
      if (!seconds) return '0 мин'
      
      const hours = Math.floor(seconds / 3600)
      const minutes = Math.floor((seconds % 3600) / 60)
      
      if (hours > 0) {
        return `${hours}ч ${minutes}м`
      }
      return `${minutes} минут`
    },
    
    saveProfile() {
      const userData = {
        ...JSON.parse(localStorage.getItem('user') || '{}'),
        username: this.user.username,
        email: this.user.email
      }
      
      localStorage.setItem('user', JSON.stringify(userData))
      alert('Профиль сохранён!')
    },
    
    saveSettings() {
      const settings = {
        theme: this.currentTheme,
        language: this.language,
        notifications: this.notifications,
        autoSave: this.autoSave,
        autoSaveInterval: this.autoSaveInterval
      }
      
      localStorage.setItem('settings', JSON.stringify(settings))
      localStorage.setItem('theme', this.currentTheme)
      
      // Применяем настройки
      this.applyTheme(this.currentTheme)
    },
    
    logout() {
      if (confirm('Вы уверены, что хотите выйти?')) {
        localStorage.removeItem('user')
        localStorage.removeItem('token')
        this.$router.push('/auth')
      }
    },
    
    resetStats() {
      if (confirm('Вы уверены, что хотите сбросить всю статистику?')) {
        localStorage.removeItem('stats')
        this.stats = {
          solvedTasks: 0,
          passedTests: 0,
          timeSpent: 0,
          rating: 0
        }
        alert('Статистика сброшена!')
      }
    }
  },
  
  watch: {
    language() {
      this.saveSettings()
    },
    
    notifications: {
      handler() {
        this.saveSettings()
      },
      deep: true
    },
    
    autoSave() {
      this.saveSettings()
    },
    
    autoSaveInterval() {
      this.saveSettings()
    }
  }
}
</script>

<style scoped>
.settings-page {
  min-height: 100vh;
  background-color: #0E1117;
  color: #E2E8F0;
}

.settings-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 20px;
}

.settings-header {
  text-align: center;
  margin-bottom: 3rem;
}

.settings-header h1 {
  margin: 0;
  font-size: 2.5rem;
  color: #E2E8F0;
  font-weight: 600;
}

.settings-header p {
  color: #9CA3AF;
  margin-top: 10px;
  font-size: 1.1rem;
}

.settings-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
}

.settings-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.section-card {
  background-color: #303030;
  border-radius: 16px;
  padding: 1.25rem;
  border: 1px solid #404040;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.35);
}

.section-card h2 {
  margin-top: 0;
  margin-bottom: 1.5rem;
  color: #F8FAFC;
  font-size: 1.3rem;
  font-weight: 600;
}

.form-group {
  margin-bottom: 1.25rem;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  color: #E2E8F0;
  font-weight: 500;
}

.form-input,
.form-select {
  width: 100%;
  background: #1E1E1E;
  border: 1px solid #404040;
  color: #E2E8F0;
  border-radius: 8px;
  padding: 0.75rem;
  font-size: 0.95rem;
  transition: all 0.2s;
  font-family: inherit;
}

.form-input::placeholder {
  color: #6B7280;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: #3B82F6;
  background: #252525;
}

.role-display {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.75rem;
  background: #1E1E1E;
  border-radius: 12px;
  border: 1px solid #404040;
}

.role-badge {
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 0.9rem;
  font-weight: 600;
  white-space: nowrap;
}

.role-badge.student {
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
}

.role-badge.teacher {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: white;
}

.role-badge.admin {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
}

.role-description {
  color: #9CA3AF;
  font-size: 0.9rem;
  flex: 1;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #404040;
}

/* Темы */
.theme-selector {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.theme-btn {
  padding: 0.75rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
  border: 2px solid transparent;
  font-size: 0.9rem;
}

.theme-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.theme-btn--active {
  border-color: #3B82F6 !important;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.3);
}

/* Чекбоксы и радио */
.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.checkbox-label,
.radio-label {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  cursor: pointer;
  color: #E2E8F0;
}

.checkbox,
.radio {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.auto-save-options {
  margin-left: 28px;
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* Статистика */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.stat-item {
  background: #1E1E1E;
  border: 1px solid #404040;
  border-radius: 12px;
  padding: 1rem;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-label {
  color: #9CA3AF;
  font-size: 0.85rem;
  margin-bottom: 0.5rem;
  text-align: center;
}

.stat-value {
  color: #3B82F6;
  font-size: 1.5rem;
  font-weight: 700;
}

/* Кнопки */
.btn {
  padding: 0.7rem 1.1rem;
  border: none;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  font-size: 0.85rem;
}

.btn-primary {
  background: #3B82F6;
  color: white;
  box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
}

.btn-primary:hover {
  background: #2563EB;
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(37, 99, 235, 0.4);
}

.btn-primary:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.4);
}

.btn-outline {
  background: #303030;
  color: #E2E8F0;
  border: 1px solid #404040;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.btn-outline:hover {
  background: #404040;
  color: #E2E8F0;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(255, 255, 255, 0.1);
}

.btn-outline:active {
  transform: translateY(0);
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
}

.btn-danger {
  background: #EF4444;
  color: white;
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.3);
}

.btn-danger:hover {
  background: #DC2626;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.4);
}

.btn-danger:active {
  transform: translateY(0);
}

.btn-sm {
  padding: 0.45rem 0.9rem;
  font-size: 0.7rem;
  border-radius: 8px;
}

/* Адаптация */
@media (max-width: 768px) {
  .settings-content {
    grid-template-columns: 1fr;
  }
  
  .settings-header h1 {
    font-size: 2rem;
  }
  
  .theme-selector {
    grid-template-columns: 1fr;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .settings-container {
    padding: 1rem 10px;
  }
  
  .section-card {
    padding: 1rem;
    border-radius: 12px;
  }
  
  .settings-header h1 {
    font-size: 2rem;
  }
}
</style>