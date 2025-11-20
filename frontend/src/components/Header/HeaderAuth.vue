<template>
  <div class="auth-container">
    <a v-if="!isLoggedIn" href="#" class="login-btn" @click.prevent="handleLogin">
      Войти
    </a>
    
    <div v-else class="profile-dropdown">
      <button class="profile-btn" 
              @click="toggleDropdown"
              :class="{ 'profile-btn--active': isDropdownOpen }">
        <span class="user-avatar">{{ userAvatar }}</span>
        <span class="username">{{ username }}</span>
        <span class="dropdown-arrow" :class="{ 'dropdown-arrow--open': isDropdownOpen }">▼</span>
      </button>
      
      <transition name="dropdown">
        <div v-if="isDropdownOpen" class="dropdown-menu">
          <button @click="handleLogout" 
                  class="dropdown-item dropdown-item--logout"
                  :class="{ 'dropdown-item--active': isLogoutActive }"
                  @mousedown="setLogoutActive(true)"
                  @mouseup="setLogoutActive(false)"
                  @mouseleave="setLogoutActive(false)">
            <span class="dropdown-icon">🚪</span>Выйти
          </button>
        </div>
      </transition>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HeaderAuth',
  props: {
    isLoggedIn: Boolean,
    username: String,
    userAvatar: String,
    activeLink: String
  },
  data() {
    return {
      isDropdownOpen: false,
      isLogoutActive: false
    }
  },
  methods: {
    toggleDropdown() {
      this.isDropdownOpen = !this.isDropdownOpen
      this.$emit('dropdown-toggle', this.isDropdownOpen)
    },
    
    handleLogin() {
      this.$router.push('/auth')
    },
    
    handleLogout() {
      this.isDropdownOpen = false
      this.isLogoutActive = false
      this.$emit('logout')
      
      // Принудительно обновляем компонент после выхода
      this.$nextTick(() => {
        this.$forceUpdate()
      })
    },
    
    setLogoutActive(active) {
      this.isLogoutActive = active
    },
    
    closeDropdown(event) {
      if (this.isDropdownOpen && !this.$el.contains(event.target)) {
        this.isDropdownOpen = false
        this.$emit('dropdown-toggle', false)
      }
    }
  },
  
  mounted() {
    document.addEventListener('click', this.closeDropdown)
  },
  
  beforeUnmount() {
    document.removeEventListener('click', this.closeDropdown)
  }
}
</script>

<style scoped>
/* Стили остаются без изменений */
.auth-container {
  margin-left: auto;
  margin-right: 15px;
  position: relative;
}

.login-btn {
  background: #3B82F6;
  color: #E5E7EB;
  padding: 8px 24px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.2s;
  display: block;
  text-decoration: none;
  cursor: pointer;
  border: none;
  font-family: inherit;
}

.login-btn:hover {
  background: #2563EB;
  transform: translateY(-1px);
}

.login-btn:active {
  transform: scale(0.95);
  background: #1D4ED8;
}

.profile-dropdown {
  position: relative;
}

.profile-btn {
  background: transparent;
  color: #9CA3AF;
  padding: 8px 16px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  border: none;
  cursor: pointer;
  font-family: inherit;
  font-size: 16px;
}

.profile-btn:hover {
  color: #E5E7EB;
  background: #1a1f2e;
}

.profile-btn--active {
  color: #E5E7EB;
  background: #3B82F6;
}

.profile-btn:active {
  transform: scale(0.95);
  background: #2563EB;
}

.user-avatar {
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

.username {
  font-size: 16px;
}

.dropdown-arrow {
  transition: transform 0.2s ease;
  font-size: 12px;
  color: #9CA3AF;
}

.dropdown-arrow--open {
  transform: rotate(180deg);
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 8px;
  background: #1F2937;
  border: 1px solid #374151;
  border-radius: 8px;
  padding: 8px 0;
  min-width: 150px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
  z-index: 1002;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  color: #9CA3AF;
  text-decoration: none;
  transition: all 0.2s;
  background: none;
  border: none;
  width: calc(100% - 16px);
  text-align: left;
  font-family: inherit;
  font-size: 16px;
  cursor: pointer;
  border-radius: 6px;
  margin: 0 8px;
  transform: scale(1);
  transition: all 0.1s ease;
}

.dropdown-item:hover {
  color: #E5E7EB;
  background: #1a1f2e;
}

.dropdown-item--active {
  transform: scale(0.95);
  background: #7F1D1D !important;
}

.dropdown-item--logout {
  color: #EF4444;
}

.dropdown-item--logout:hover {
  background: #7F1D1D;
  color: white;
}

.dropdown-icon {
  width: 20px;
  text-align: center;
  font-size: 16px;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s ease;
}

.dropdown-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}

.dropdown-enter-to {
  opacity: 1;
  transform: translateY(0);
}

.dropdown-leave-from {
  opacity: 1;
  transform: translateY(0);
}

.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 1023px) {
  .login-btn,
  .profile-btn {
    padding: 6px 20px;
    font-size: 14px;
  }

  .user-avatar {
    width: 28px;
    height: 28px;
    font-size: 12px;
  }
}

@media (max-width: 767px) {
  .auth-container {
    display: none;
  }
}
</style>