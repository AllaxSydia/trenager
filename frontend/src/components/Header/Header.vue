<template>
  <div class="header-container">
    <HeaderLogo />
    <HeaderNav />
    <HeaderAuth 
      :isLoggedIn="isLoggedIn"
      :username="username"
      :userAvatar="userAvatar"
      :activeLink="activeLink"
      @navigation="handleNavigation"
      @logout="handleLogout"
      @dropdown-toggle="handleDropdownToggle"
    />
    <HeaderMobile
      :isLoggedIn="isLoggedIn"
      :username="username"
      :userAvatar="userAvatar"
      :activeLink="activeLink"
      @navigation="handleNavigation"
      @logout="handleLogout"
      @mobile-toggle="handleMobileToggle"
    />
  </div>
</template>

<script>
import HeaderLogo from './HeaderLogo.vue'
import HeaderNav from './HeaderNav.vue'
import HeaderAuth from './HeaderAuth.vue'
import HeaderMobile from './HeaderMobile.vue'

export default {
  name: 'Header',
  components: {
    HeaderLogo,
    HeaderNav,
    HeaderAuth,
    HeaderMobile
  },
  data() {
    return {
      isLoggedIn: false,
      username: '',
      userAvatar: '',
      activeLink: '',
      isMobileMenuOpen: false,
      isProfileDropdownOpen: false
    }
  },
  methods: {
    handleNavigation(link) {
      this.activeLink = link
      console.log(`Переход на: ${link}`)
    },
    
    handleLogout() {
      this.isLoggedIn = false
      this.username = ''
      this.userAvatar = ''
      localStorage.removeItem('user')
      console.log('Выход из системы')
      
      // Принудительно обновляем состояние
      this.$forceUpdate()
    },
    
    handleDropdownToggle(isOpen) {
      this.isProfileDropdownOpen = isOpen
      if (isOpen) {
        this.isMobileMenuOpen = false
      }
    },
    
    handleMobileToggle(isOpen) {
      this.isMobileMenuOpen = isOpen
      if (isOpen) {
        this.isProfileDropdownOpen = false
      }
    },
    
    updateUserData(userData) {
      this.isLoggedIn = userData.isLoggedIn
      this.username = userData.username
      this.userAvatar = userData.userAvatar
    },
    
    checkAuthStatus() {
      const savedUser = localStorage.getItem('user')
      console.log('Checking auth status:', !!savedUser)
      if (savedUser) {
        try {
          const userData = JSON.parse(savedUser)
          this.updateUserData(userData)
        } catch (e) {
          console.error('Ошибка при чтении данных пользователя:', e)
          localStorage.removeItem('user')
        }
      } else {
        // Явно сбрасываем статус если пользователь не авторизован
        this.isLoggedIn = false
        this.username = ''
        this.userAvatar = ''
      }
    }
  },
  
  mounted() {
    this.checkAuthStatus()
    
    document.addEventListener('click', (event) => {
      if (!this.$el.contains(event.target)) {
        this.isMobileMenuOpen = false
        this.isProfileDropdownOpen = false
      }
    })
  },
  
  // Добавляем watcher для отслеживания изменений маршрута
  watch: {
    '$route'() {
      // При смене маршрута проверяем статус авторизации
      this.$nextTick(() => {
        this.checkAuthStatus()
      })
    }
  }
}
</script>

<style scoped>
.header-container {
  max-width: 100%;
  margin: 0 auto;
  padding: 0 6px;
  height: 96px;
  display: flex;
  align-items: center;
  border-bottom: 1px solid #374151;
  position: fixed; /* Добавлено */
  top: 0; /* Добавлено */
  left: 0; /* Добавлено */
  right: 0; /* Добавлено */
  background: #0E1117;
  z-index: 1000; /* Добавлено для обеспечения поверх другого контента */
}

@media (max-width: 1023px) {
  .header-container {
    padding: 0 20px;
    height: 80px;
  }
}

@media (max-width: 767px) {
  .header-container {
    padding: 0 15px;
    height: 70px;
    justify-content: space-between;
  }
}

@media (max-width: 480px) {
  .header-container {
    padding: 0 10px;
    height: 60px;
  }
}
</style>