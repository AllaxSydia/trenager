<template>
  <nav class="nav">
    <!-- Основные курсы -->
    <div class="nav-dropdown">
      <button class="nav-link dropdown-toggle">
        Курсы
      </button>
      <div class="dropdown-menu">
        <router-link to="/courses/python" class="dropdown-item">
          <span class="lang-icon">🐍</span>
          <span>Python</span>
        </router-link>
        <router-link to="/courses/javascript" class="dropdown-item">
          <span class="lang-icon">🟨</span>
          <span>JavaScript</span>
        </router-link>
        <router-link to="/courses/java" class="dropdown-item">
          <span class="lang-icon">☕</span>
          <span>Java</span>
        </router-link>
        <router-link to="/courses/cpp" class="dropdown-item">
          <span class="lang-icon">⚙️</span>
          <span>C++</span>
        </router-link>
        <router-link to="/courses/go" class="dropdown-item">
          <span class="lang-icon">🐹</span>
          <span>Go</span>
        </router-link>
        <div class="dropdown-divider"></div>
        <router-link to="/" class="dropdown-item">
          <span class="lang-icon">📚</span>
          <span>Все курсы</span>
        </router-link>
      </div>
    </div>
    
    <!-- Рейтинг -->
    <router-link to="/rating" class="nav-link" active-class="nav-link--active">
      <span class="nav-text">Рейтинг</span>
    </router-link>
    
    <!-- Профиль -->
    <router-link to="/profile" class="nav-link" active-class="nav-link--active">
      <span class="nav-text">Профиль</span>
    </router-link>
    
    <!-- Настройки -->
    <router-link to="/settings" class="nav-link" active-class="nav-link--active">
      <span class="nav-text">Настройки</span>
    </router-link>
    
    <!-- Панель учителя -->
    <div v-if="userRole === 'teacher'" class="nav-dropdown teacher-dropdown">
      <button class="nav-link dropdown-toggle teacher-toggle">
        <span class="nav-text">Учитель</span>
      </button>
      <div class="dropdown-menu teacher-menu">
        <router-link to="/teacher/tasks" class="dropdown-item">
          <span>Создание задач</span>
        </router-link>
        <router-link to="/teacher/statistics" class="dropdown-item">
          <span>Статистика</span>
        </router-link>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'HeaderNav',
  props: {
    userRole: {
      type: String,
      default: 'student'
    }
  },
  data() {
    return {
      showCoursesDropdown: false,
      showTeacherDropdown: false
    }
  },
  mounted() {
    console.log('HeaderNav mounted, userRole:', this.userRole)
    
    // Закрываем дропдауны при клике вне меню
    document.addEventListener('click', this.closeDropdowns)
  },
  beforeUnmount() {
    document.removeEventListener('click', this.closeDropdowns)
  },
  watch: {
    userRole(newRole) {
      console.log('HeaderNav userRole changed to:', newRole)
    }
  },
  methods: {
    closeDropdowns(event) {
      // Закрываем дропдауны если клик был вне них
      if (!event.target.closest('.nav-dropdown')) {
        this.showCoursesDropdown = false
        this.showTeacherDropdown = false
      }
    },
    
    toggleCoursesDropdown() {
      this.showCoursesDropdown = !this.showCoursesDropdown
      if (this.showCoursesDropdown) {
        this.showTeacherDropdown = false
      }
    },
    
    toggleTeacherDropdown() {
      this.showTeacherDropdown = !this.showTeacherDropdown
      if (this.showTeacherDropdown) {
        this.showCoursesDropdown = false
      }
    }
  }
}
</script>

<style scoped>
.nav {
  display: flex;
  gap: 20px;
  align-items: center;
}

/* Основные стили ссылок */
.nav-link {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: transparent;
  border: none;
  border-radius: 8px;
  color: #9CA3AF;
  cursor: pointer;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  font-size: 15px;
  font-weight: 500;
  height: 40px;
  padding: 0 16px;
  position: relative;
  text-decoration: none;
  transition: all 0.2s ease;
  white-space: nowrap;
  outline: none;
  user-select: none;
}

.nav-link::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(255, 255, 255, 0);
  border-radius: 8px;
  transition: background-color 0.2s ease;
  z-index: -1;
}

.nav-link:hover::before {
  background-color: rgba(255, 255, 255, 0.05);
}

.nav-link:hover {
  color: #E5E7EB;
  transform: translateY(-1px);
}

.nav-link--active {
  color: #3B82F6;
  font-weight: 600;
}

.nav-link--active::before {
  background-color: rgba(59, 130, 246, 0.1);
}

.nav-text {
  display: inline-block;
}

/* Стили для дропдаунов */
.nav-dropdown {
  position: relative;
}

.dropdown-toggle {
  cursor: pointer;
}

.dropdown-toggle::after {
  content: '▼';
  font-size: 10px;
  margin-left: 4px;
  opacity: 0.7;
  transition: transform 0.2s ease;
}

.nav-dropdown:hover .dropdown-toggle::after {
  transform: rotate(180deg);
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  min-width: 200px;
  background: #1F2937;
  border-radius: 8px;
  border: 1px solid #374151;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
  opacity: 0;
  visibility: hidden;
  transform: translateY(-10px);
  transition: all 0.2s ease;
  z-index: 1000;
  margin-top: 8px;
  padding: 8px 0;
}

.nav-dropdown:hover .dropdown-menu {
  opacity: 1;
  visibility: visible;
  transform: translateY(0);
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  color: #D1D5DB;
  text-decoration: none;
  transition: all 0.2s ease;
  font-size: 14px;
}

.dropdown-item:hover {
  background: rgba(59, 130, 246, 0.1);
  color: #3B82F6;
}

.dropdown-icon {
  font-size: 16px;
  width: 20px;
  text-align: center;
}

.lang-icon {
  font-size: 18px;
  width: 24px;
  text-align: center;
}

.dropdown-divider {
  height: 1px;
  background: #374151;
  margin: 8px 0;
}

/* Особые стили для панели учителя */
.teacher-dropdown .dropdown-menu {
  min-width: 180px;
}

.teacher-toggle {
  color: #F59E0B;
}

.teacher-toggle:hover {
  color: #FBBF24;
}

.teacher-toggle.nav-link--active {
  color: #F59E0B;
}

.teacher-toggle.nav-link--active::before {
  background-color: rgba(245, 158, 11, 0.1);
}

/* Адаптация для средних экранов */
@media (max-width: 1023px) {
  .nav {
    gap: 16px;
  }
  
  .nav-link {
    padding: 0 14px;
    font-size: 14px;
    height: 36px;
  }
  
  .nav-text {
    display: none; /* Скрываем текст, оставляем только иконки */
  }
  
  .nav-link {
    width: 40px;
    justify-content: center;
    padding: 0;
  }
  
  .nav-icon {
    font-size: 20px;
    margin: 0;
  }
  
  .dropdown-menu {
    left: -80px;
  }
  
  .nav-dropdown .nav-link::after {
    display: none;
  }
}

/* Адаптация для маленьких экранов */
@media (max-width: 767px) {
  .nav {
    display: none; /* На мобильных показываем в гамбургер-меню */
  }
}

/* Для очень больших экранов */
@media (min-width: 1920px) {
  .nav {
    gap: 24px;
  }
  
  .nav-link {
    font-size: 16px;
    height: 44px;
    padding: 0 20px;
  }
  
  .dropdown-menu {
    min-width: 240px;
  }
  
  .dropdown-item {
    padding: 12px 20px;
    font-size: 15px;
  }
}

/* Плавные переходы */
.nav-link:active {
  transform: scale(0.95);
  transition-duration: 0.1s;
}

.dropdown-item:active {
  transform: scale(0.98);
}

/* Анимация появления дропдауна */
@keyframes dropdownAppear {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.dropdown-menu {
  animation: dropdownAppear 0.2s ease;
}
</style>