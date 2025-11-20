<template>
  <div class="mobile-header" v-if="isMobile">
    <div class="mobile-header-content">
      <h1 class="mobile-title">{{ title }}</h1>
      <button class="mobile-menu-btn" @click="$emit('toggle-sidebar')">
        <span>📚 Все задания</span>
      </button>
    </div>
    <div class="mobile-progress">
      <div class="progress-bar">
        <div class="progress-fill" :style="{ width: progress + '%' }"></div>
      </div>
      <span>{{ progress }}% завершено</span>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MobileHeader',
  props: {
    title: String,
    progress: Number,
    isMobile: Boolean
  },
  emits: ['toggle-sidebar']
}
</script>

<style scoped>
.mobile-header {
  display: none;
  background: #303030;
  border-radius: 16px;
  padding: 1rem;
  margin-bottom: 1rem;
  border: 1px solid #404040;
}

.mobile-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.mobile-title {
  margin: 0;
  color: #F8FAFC;
  font-size: 1.4rem;
  font-weight: 600;
}

.mobile-menu-btn {
  display: none;
  background: #3B82F6;
  color: white;
  border: none;
  border-radius: 12px;
  padding: 0.75rem 1.5rem;
  font-size: 0.9rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.mobile-menu-btn:active {
  transform: translateY(2px);
  box-shadow: 0 2px 6px rgba(59, 130, 246, 0.4);
}

.mobile-progress {
  text-align: center;
}

.mobile-progress .progress-bar {
  background: #404040;
  border-radius: 8px;
  height: 6px;
  margin-bottom: 0.5rem;
  overflow: hidden;
}

.mobile-progress .progress-fill {
  background: #3B82F6;
  height: 100%;
  transition: width 0.3s ease;
}

.mobile-progress span {
  font-size: 0.8rem;
  color: #94A3B8;
  font-weight: 500;
}

@media (max-width: 1024px) {
  .mobile-header {
    display: block;
  }
  
  .mobile-menu-btn {
    display: block;
  }
}

@media (max-width: 768px) {
  .mobile-header {
    padding: 0.75rem;
    border-radius: 12px;
  }
  
  .mobile-title {
    font-size: 1.2rem;
  }
}
</style>