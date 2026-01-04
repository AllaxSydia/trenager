<template>
  <div class="settings-page">
    <div class="settings-container">
      <!-- Хедер -->
      <div class="settings-header">
        <h1>⚙️ Настройки</h1>
        <p>Управление вашим аккаунтом и настройками платформы</p>
      </div>

      <!-- Основной контент -->
      <div class="settings-content">
        <!-- Левая колонка - настройки аккаунта -->
        <div class="settings-section">
          <div class="section-card">
            <h2>👤 Настройки аккаунта</h2>
            
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
                💾 Сохранить изменения
              </button>
              <button @click="logout" class="btn btn-outline">
                🚪 Выйти
              </button>
            </div>
          </div>
        </div>

        <!-- Правая колонка - настройки платформы -->
        <div class="settings-section">
          <div class="section-card">
            <h2>🎨 Настройки интерфейса</h2>
            
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
            <h2>📊 Статистика</h2>
            
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
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
  color: #e2e8f0;
  padding: 20px;
}

.settings-container {
  max-width: 1200px;
  margin: 0 auto;
}

.settings-header {
  text-align: center;
  margin-bottom: 40px;
  padding: 20px;
  background: rgba(30, 41, 59, 0.8);
  border-radius: 12px;
  border: 1px solid #334155;
}

.settings-header h1 {
  margin: 0;
  font-size: 2.5rem;
  background: linear-gradient(135deg, #60a5fa, #3b82f6);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.settings-header p {
  color: #94a3b8;
  margin-top: 10px;
  font-size: 1.1rem;
}

.settings-content {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 30px;
}

.settings-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-card {
  background: rgba(30, 41, 59, 0.8);
  border-radius: 12px;
  padding: 25px;
  border: 1px solid #334155;
}

.section-card h2 {
  margin-top: 0;
  margin-bottom: 25px;
  color: #f8fafc;
  font-size: 1.3rem;
}

.form-group {
  margin-bottom: 20px;
}

.form-label {
  display: block;
  margin-bottom: 8px;
  color: #cbd5e1;
  font-weight: 500;
}

.form-input,
.form-select {
  width: 100%;
  background: #0f172a;
  border: 1px solid #475569;
  color: #e2e8f0;
  border-radius: 6px;
  padding: 12px;
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: #3b82f6;
}

.role-display {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 12px;
  background: #0f172a;
  border-radius: 6px;
  border: 1px solid #475569;
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
  color: #94a3b8;
  font-size: 0.9rem;
  flex: 1;
}

.form-actions {
  display: flex;
  gap: 15px;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #334155;
}

/* Темы */
.theme-selector {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

.theme-btn {
  padding: 12px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
  border: 2px solid transparent;
}

.theme-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.theme-btn--active {
  border-color: #3b82f6 !important;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.3);
}

/* Чекбоксы и радио */
.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.checkbox-label,
.radio-label {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  color: #cbd5e1;
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
  gap: 15px;
  margin-bottom: 20px;
}

.stat-item {
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 15px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-label {
  color: #94a3b8;
  font-size: 0.85rem;
  margin-bottom: 5px;
  text-align: center;
}

.stat-value {
  color: #3b82f6;
  font-size: 1.5rem;
  font-weight: 700;
}

/* Кнопки */
.btn {
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 1rem;
}

.btn-primary {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
  color: white;
}

.btn-primary:hover {
  background: linear-gradient(135deg, #2563eb, #1d4ed8);
  transform: translateY(-1px);
}

.btn-outline {
  background: transparent;
  border: 1px solid #475569;
  color: #94a3b8;
}

.btn-outline:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #e2e8f0;
}

.btn-danger {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
}

.btn-danger:hover {
  background: linear-gradient(135deg, #dc2626, #b91c1c);
}

.btn-sm {
  padding: 8px 16px;
  font-size: 0.9rem;
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
  .settings-page {
    padding: 10px;
  }
  
  .section-card {
    padding: 20px;
  }
  
  .settings-header {
    padding: 15px;
  }
}
</style>