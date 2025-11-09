<template>
  <aside class="sidebar" :class="{ 'sidebar--mobile-hidden': isMobile }">
    <div class="sidebar-content">
      <div class="sidebar-header">
        <h2>{{ title }}</h2>
        <div class="progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: progress + '%' }"></div>
          </div>
          <span>{{ progress }}% завершено</span>
        </div>
      </div>
      
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
  </aside>
</template>

<script>
export default {
  name: 'LessonSidebar',
  props: {
    lessons: Array,
    currentLesson: Object,
    progress: Number,
    title: String,
    isMobile: Boolean
  },
  emits: ['select-lesson']
}
</script>

<style scoped>
.sidebar {
  background: #303030;
  border: 1px solid #404040;
  border-radius: 16px;
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.4);
  position: sticky;
  top: 20px;
  height: fit-content;
  max-height: calc(100vh - 140px);
  overflow: hidden;
  margin-left: 0;
}

.sidebar--mobile-hidden {
  display: none;
}

.sidebar-content {
  height: 100%;
  overflow-y: auto;
  max-height: calc(100vh - 140px);
  padding-right: 8px;
  margin-right: -8px;
}

.sidebar-content::-webkit-scrollbar {
  width: 12px;
}

.sidebar-content::-webkit-scrollbar-track {
  background: #252525;
  border-radius: 8px;
  margin: 10px 0;
  border: 1px solid #404040;
  margin-left: 4px;
}

.sidebar-content::-webkit-scrollbar-thumb {
  background: linear-gradient(135deg, #3B82F6, #8B5CF6);
  border-radius: 8px;
  border: 3px solid #252525;
  min-height: 60px;
}

.sidebar-content::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(135deg, #2563EB, #7C3AED);
}

.sidebar-content {
  scrollbar-width: thin;
  scrollbar-color: #3B82F6 #252525;
}

.sidebar-header {
  padding: 1.25rem;
  border-bottom: 1px solid #404040;
  position: sticky;
  top: 0;
  background: #303030;
  z-index: 10;
  border-radius: 16px 16px 0 0;
}

.sidebar-header h2 {
  margin: 0 0 0.75rem 0;
  color: #F8FAFC;
  font-size: 1.2rem;
  font-weight: 600;
}

.progress-bar {
  background: #404040;
  border-radius: 8px;
  height: 6px;
  margin-bottom: 0.5rem;
  overflow: hidden;
}

.progress-fill {
  background: #3B82F6;
  height: 100%;
  transition: width 0.3s ease;
}

.progress span {
  font-size: 0.8rem;
  color: #94A3B8;
  font-weight: 500;
}

.lessons-list {
  padding: 0.875rem;
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

.lesson-item--completed.lesson-item--active {
  background: rgba(30, 41, 59, 0.9);
  border: 2px solid #3B82F6;
  box-shadow: 0 4px 20px rgba(59, 130, 246, 0.3);
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

.lesson-item:hover .lesson-icon {
  background: rgba(59, 130, 246, 0.2);
  transform: scale(1.03);
}

.lesson-item:active .lesson-icon {
  transform: scale(0.98);
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

.lesson-item--completed.lesson-item--active .lesson-icon {
  background: #3B82F6;
  color: white;
  border-color: #3B82F6;
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

@media (max-width: 1024px) {
  .sidebar {
    display: none;
  }
}

@media (max-width: 768px) {
  .sidebar-header {
    padding: 1rem;
  }
  
  .lessons-list {
    padding: 0.75rem;
  }
  
  .lesson-item {
    padding: 0.875rem;
  }
}
</style>