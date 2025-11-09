<template>
  <div v-if="isMobile && showSidebar" class="sidebar-overlay" @click="$emit('close')">
    <div class="sidebar-mobile" :class="{ 'sidebar-mobile--active': showSidebar }" @click.stop>
      <div class="sidebar-mobile-header">
        <h2>{{ title }}</h2>
        <button class="close-btn" @click="$emit('close')">✕</button>
      </div>
      <div class="sidebar-mobile-content">
        <div class="lessons-list">
          <div
            v-for="lesson in lessons" 
            :key="lesson.id"
            class="lesson-item"
            :class="{
              'lesson-item--active': currentLesson.id === lesson.id,
              'lesson-item--completed': lesson.completed
            }"
            @click="$emit('select-lesson', lesson)"
          >
            <div class="lesson-icon">
              <span v-if="lesson.completed">✓</span>
              <span v-else>{{ lesson.id }}</span>
            </div>
            <div class="lesson-info">
              <h3>{{ lesson.title }}</h3>
              <p>{{ lesson.description }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'MobileSidebar',
  props: {
    lessons: Array,
    currentLesson: Object,
    title: String,
    isMobile: Boolean,
    showSidebar: Boolean
  },
  emits: ['select-lesson', 'close']
}
</script>

<style scoped>
.sidebar-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
  opacity: 0;
  animation: overlayFadeIn 0.3s ease forwards;
}

@keyframes overlayFadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.sidebar-mobile {
  background: #303030;
  border-radius: 16px;
  width: 100%;
  max-width: 400px;
  max-height: 80vh;
  overflow: hidden;
  border: 1px solid #404040;
  transform: translateY(-20px) scale(0.95);
  opacity: 0;
  animation: sidebarSlideIn 0.3s ease 0.1s forwards;
}

@keyframes sidebarSlideIn {
  from {
    transform: translateY(-20px) scale(0.95);
    opacity: 0;
  }
  to {
    transform: translateY(0) scale(1);
    opacity: 1;
  }
}

.sidebar-mobile-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem;
  border-bottom: 1px solid #404040;
  background: #303030;
}

.sidebar-mobile-header h2 {
  margin: 0;
  color: #F8FAFC;
  font-size: 1.2rem;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  color: #94A3B8;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 0;
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  transition: all 0.15s ease;
}

.close-btn:active {
  background: rgba(148, 163, 184, 0.1);
  transform: scale(0.95);
}

.sidebar-mobile-content {
  max-height: calc(80vh - 80px);
  overflow-y: auto;
  padding: 1rem;
}

.lessons-list {
  display: flex;
  flex-direction: column;
  gap: 0.625rem;
}

.lesson-item {
  background: rgba(14, 17, 23, 0.75);
  border-radius: 14px;
  padding: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: flex-start;
  gap: 0.875rem;
  position: relative;
  overflow: hidden;
}

.lesson-item::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, #3B82F6, #8B5CF6);
  transform: scaleX(0);
  transition: transform 0.3s ease;
}

.lesson-item:hover {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.25);
  border-color: rgba(59, 130, 246, 0.3);
}

.lesson-item:hover::before {
  transform: scaleX(1);
}

.lesson-item:active {
  transform: translateY(0);
  transition: all 0.1s ease;
}

.lesson-item--active {
  background: rgba(30, 41, 59, 0.9);
  border: 2px solid #3B82F6;
  box-shadow: 0 4px 20px rgba(59, 130, 246, 0.3);
  transform: translateY(-2px);
}

.lesson-item--active::before {
  transform: scaleX(1);
  background: linear-gradient(90deg, #3B82F6, #60A5FA);
  height: 3px;
}

.lesson-item--completed {
  background: rgba(16, 44, 87, 0.75);
  border-color: #1E40AF;
}

.lesson-item--completed::before {
  background: linear-gradient(90deg, #1E40AF, #3730A3);
}

.lesson-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: rgba(59, 130, 246, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  font-weight: 600;
  flex-shrink: 0;
  transition: all 0.3s ease;
  border: 1px solid rgba(59, 130, 246, 0.2);
  color: #3B82F6;
}

.lesson-item--active .lesson-icon {
  background: #3B82F6;
  color: white;
  border-color: #3B82F6;
  transform: scale(1.05);
}

.lesson-item--completed .lesson-icon {
  background: #1E40AF;
  color: white;
  border-color: #1E40AF;
}

.lesson-info {
  flex: 1;
  min-width: 0;
}

.lesson-info h3 {
  margin: 0 0 0.375rem 0;
  font-size: 0.9rem;
  color: #F8FAFC;
  font-weight: 600;
  line-height: 1.3;
}

.lesson-info p {
  margin: 0;
  font-size: 0.8rem;
  color: #94A3B8;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-clamp: 2;
}

@media (max-width: 768px) {
  .sidebar-mobile {
    max-width: 100%;
    max-height: 90vh;
  }
  
  .sidebar-mobile-content {
    max-height: calc(90vh - 80px);
  }
  
  .sidebar-mobile-header {
    padding: 1rem;
  }
  
  .lesson-item {
    padding: 0.875rem;
  }
}
</style>