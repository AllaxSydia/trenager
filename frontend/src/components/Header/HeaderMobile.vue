<template>
  <div>
    <button class="mobile-menu-btn" @click="toggleMobileMenu"
            :class="{ 'mobile-menu-btn--active': isMobileMenuActive }"
            @mousedown="setMobileMenuActive(true)"
            @mouseup="setMobileMenuActive(false)"
            @mouseleave="setMobileMenuActive(false)">
      <span :class="['hamburger', { 'hamburger--active': isMobileMenuOpen }]">
        <span class="hamburger-line"></span>
        <span class="hamburger-line"></span>
        <span class="hamburger-line"></span>
      </span>
    </button>

    <transition name="slide-down">
      <div v-if="isMobileMenuOpen" class="mobile-menu">
        <!-- ИСПРАВЛЕНО: добавлен prevent и прямая навигация -->
        <a href="#" class="mobile-nav-link" 
           @click.prevent="goToRating"
           :class="{ 'mobile-nav-link--active': $route.name === 'rating' }">
          <span class="link-content">Рейтинг</span>
        </a>
        <a href="#" class="mobile-nav-link" 
           @click.prevent="goToProfile"
           :class="{ 'mobile-nav-link--active': $route.name === 'profile' }">
          <span class="link-content">Профиль</span>
        </a>
        
        <div v-if="isLoggedIn" class="mobile-profile-section">
          <div class="mobile-profile-info">
            <span class="mobile-user-avatar">{{ userAvatar }}</span>
            <span class="mobile-username">{{ username }}</span>
          </div>
          
          <button @click="handleLogout" 
                  class="mobile-dropdown-item mobile-dropdown-item--logout"
                  :class="{ 'mobile-dropdown-item--active': isMobileLogoutActive }"
                  @mousedown="setMobileLogoutActive(true)"
                  @mouseup="setMobileLogoutActive(false)"
                  @mouseleave="setMobileLogoutActive(false)">
            <span class="dropdown-icon">🚪</span>Выйти
          </button>
        </div>
        
        <!-- ИСПРАВЛЕНО: добавлен prevent -->
        <a v-else href="#" class="mobile-login-btn" @click.prevent="handleLogin"
           :class="{ 'mobile-login-btn--active': isMobileLoginActive }"
           @mousedown="setMobileLoginActive(true)"
           @mouseup="setMobileLoginActive(false)"
           @mouseleave="setMobileLoginActive(false)">
          Войти
        </a>
      </div>
    </transition>

    <transition name="fade">
      <div v-if="isMobileMenuOpen" class="mobile-overlay" @click="closeMobileMenu"></div>
    </transition>
  </div>
</template>

<script>
export default {
  name: 'HeaderMobile',
  props: {
    isLoggedIn: Boolean,
    username: String,
    userAvatar: String,
    activeLink: String
  },
  data() {
    return {
      isMobileMenuOpen: false,
      isMobileMenuActive: false,
      isMobileLogoutActive: false,
      isMobileLoginActive: false
    }
  },
  methods: {
    toggleMobileMenu() {
      this.isMobileMenuOpen = !this.isMobileMenuOpen
      this.$emit('mobile-toggle', this.isMobileMenuOpen)
    },
    
    closeMobileMenu() {
      this.isMobileMenuOpen = false
      this.$emit('mobile-toggle', false)
    },
    
    // ИСПРАВЛЕНО: прямые методы навигации
    goToRating() {
      this.$router.push('/rating')
      this.closeMobileMenu()
    },
    
    goToProfile() {
      this.$router.push('/profile')
      this.closeMobileMenu()
    },
    
    handleLogin() {
      this.$router.push('/auth')
      this.closeMobileMenu()
    },
    
    handleLogout() {
      this.$emit('logout')
      this.closeMobileMenu()
    },
    
    setMobileMenuActive(active) {
      this.isMobileMenuActive = active
    },
    
    setMobileLogoutActive(active) {
      this.isMobileLogoutActive = active
    },
    
    setMobileLoginActive(active) {
      this.isMobileLoginActive = active
    }
  }
}
</script>

<style scoped>
.mobile-menu-btn {
  display: none;
  background: none;
  border: none;
  cursor: pointer;
  z-index: 1001;
  padding: 8px;
  transition: transform 0.2s;
  transform: scale(1);
}

.mobile-menu-btn--active {
  transform: scale(0.9);
}

.hamburger {
  display: flex;
  flex-direction: column;
  width: 24px;
  height: 18px;
  position: relative;
  transition: all 0.3s;
}

.hamburger-line {
  height: 2px;
  width: 100%;
  background: #E5E7EB;
  border-radius: 2px;
  transition: all 0.3s;
  transform-origin: center;
}

.hamburger-line:nth-child(1) {
  margin-bottom: 6px;
}

.hamburger-line:nth-child(2) {
  margin-bottom: 6px;
}

.hamburger--active .hamburger-line:nth-child(1) {
  transform: rotate(45deg) translate(6px, 6px);
}

.hamburger--active .hamburger-line:nth-child(2) {
  opacity: 0;
}

.hamburger--active .hamburger-line:nth-child(3) {
  transform: rotate(-45deg) translate(6px, -6px);
}

.mobile-menu {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #0E1117;
  border-bottom: 1px solid #374151;
  padding: 12px;
  flex-direction: column;
  gap: 4px;
  z-index: 1001;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.mobile-nav-link {
  display: block;
  text-decoration: none;
  color: #9CA3AF;
  transition: all 0.2s ease;
  border-radius: 8px;
  padding: 0;
  cursor: pointer;
}

.link-content {
  display: block;
  padding: 12px 16px;
  border-radius: 8px;
  transition: all 0.2s ease;
  transform: scale(1);
}

.mobile-nav-link:hover .link-content {
  background: #1a1f2e;
  color: #E5E7EB;
  transform: scale(1.02);
}

.mobile-nav-link--active .link-content {
  background: #3B82F6;
  color: #E5E7EB;
}

.mobile-nav-link:active .link-content {
  transform: scale(0.95);
  background: #2563EB;
}

.mobile-profile-section {
  margin-top: 16px;
  border-top: 1px solid #374151;
  padding-top: 16px;
}

.mobile-profile-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 16px 12px 16px;
  color: #E5E7EB;
  font-weight: 500;
}

.mobile-user-avatar {
  background: #3B82F6;
  color: #E5E7EB;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
}

.mobile-username {
  font-size: 16px;
}

.mobile-dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  color: #9CA3AF;
  text-decoration: none;
  transition: all 0.2s;
  border-radius: 8px;
  font-size: 16px;
  background: none;
  border: none;
  width: 100%;
  text-align: left;
  cursor: pointer;
  font-family: inherit;
  margin: 2px 0;
  transform: scale(1);
  transition: all 0.1s ease;
}

.mobile-dropdown-item:hover {
  background: #1a1f2e;
  color: #E5E7EB;
}

.mobile-dropdown-item--active {
  transform: scale(0.95);
  background: #7F1D1D !important;
}

.mobile-dropdown-item--logout {
  color: #EF4444;
}

.mobile-dropdown-item--logout:hover {
  background: #7F1D1D;
  color: white;
}

.mobile-login-btn {
  background: #3B82F6;
  color: #E5E7EB;
  padding: 16px;
  border-radius: 8px;
  font-weight: 500;
  text-align: center;
  text-decoration: none;
  margin-top: 8px;
  transition: all 0.2s;
  display: block;
  cursor: pointer;
  transform: scale(1);
  transition: all 0.1s ease;
}

.mobile-login-btn:hover {
  background: #2563EB;
  transform: translateY(-1px) scale(1);
}

.mobile-login-btn--active {
  transform: scale(0.95);
  background: #1D4ED8;
}

.mobile-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 999;
}

.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease;
}

.slide-down-enter-from {
  opacity: 0;
  transform: translateY(-20px);
}

.slide-down-enter-to {
  opacity: 1;
  transform: translateY(0);
}

.slide-down-leave-from {
  opacity: 1;
  transform: translateY(0);
}

.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.fade-enter-to,
.fade-leave-from {
  opacity: 1;
}

@media (max-width: 767px) {
  .mobile-menu-btn {
    display: block;
  }
}

@media (max-width: 480px) {
  .mobile-menu {
    padding: 8px;
  }

  .link-content {
    padding: 10px 14px;
    font-size: 16px;
  }

  .mobile-login-btn {
    padding: 14px;
    font-size: 14px;
  }

  .mobile-user-avatar {
    width: 28px;
    height: 28px;
    font-size: 12px;
  }
}
</style>