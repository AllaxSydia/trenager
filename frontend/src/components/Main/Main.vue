<template>
  <div class="main-page">
    <!-- УБРАТЬ Header отсюда -->
    
    <main class="main">
      <div class="container">
        <section class="courses-section">
          <div class="courses-grid">
            <div class="course-row">
              <div class="course-card" 
                   @click="goToCourse('python')"
                   @mousedown="setCourseActive('python', true)"
                   @mouseup="setCourseActive('python', false)"
                   @mouseleave="setCourseActive('python', false)"
                   :class="{ 'course-card--active': activeCourse === 'python' }">
                <div class="progress-overlay" :style="{ width: getProgress('python') + '%' }"></div>
                <img src="@/assets/python-logo.png" alt="Python" class="course-image" draggable="false">
                <h3>Python</h3>
              </div>
              <div class="course-card" 
                   @click="goToCourse('javascript')"
                   @mousedown="setCourseActive('javascript', true)"
                   @mouseup="setCourseActive('javascript', false)"
                   @mouseleave="setCourseActive('javascript', false)"
                   :class="{ 'course-card--active': activeCourse === 'javascript' }">
                <div class="progress-overlay" :style="{ width: getProgress('javascript') + '%' }"></div>
                <img src="@/assets/javascript-logo.png" alt="JavaScript" class="course-image" draggable="false">
                <h3>JavaScript</h3>
              </div>
            </div>
            <div class="course-row">
              <div class="course-card" 
                   @click="goToCourse('java')"
                   @mousedown="setCourseActive('java', true)"
                   @mouseup="setCourseActive('java', false)"
                   @mouseleave="setCourseActive('java', false)"
                   :class="{ 'course-card--active': activeCourse === 'java' }">
                <div class="progress-overlay" :style="{ width: getProgress('java') + '%' }"></div>
                <img src="@/assets/java-logo.png" alt="Java" class="course-image" draggable="false">
                <h3>Java</h3>
              </div>
              <div class="course-card" 
                   @click="goToCourse('golang')"
                   @mousedown="setCourseActive('golang', true)"
                   @mouseup="setCourseActive('golang', false)"
                   @mouseleave="setCourseActive('golang', false)"
                   :class="{ 'course-card--active': activeCourse === 'golang' }">
                <div class="progress-overlay" :style="{ width: getProgress('golang') + '%' }"></div>
                <img src="@/assets/golang-logo.png" alt="Golang" class="course-image" draggable="false">
                <h3>Golang</h3>
              </div>
            </div>
            <div class="course-row">
              <div class="course-card" 
                   @click="goToCourse('csharp')"
                   @mousedown="setCourseActive('csharp', true)"
                   @mouseup="setCourseActive('csharp', false)"
                   @mouseleave="setCourseActive('csharp', false)"
                   :class="{ 'course-card--active': activeCourse === 'csharp' }">
                <div class="progress-overlay" :style="{ width: getProgress('csharp') + '%' }"></div>
                <img src="@/assets/csharp-logo.png" alt="C#" class="course-image" draggable="false">
                <h3>C#</h3>
              </div>
              <div class="course-card" 
                   @click="goToCourse('cpp')"
                   @mousedown="setCourseActive('cpp', true)"
                   @mouseup="setCourseActive('cpp', false)"
                   @mouseleave="setCourseActive('cpp', false)"
                   :class="{ 'course-card--active': activeCourse === 'cpp' }">
                <div class="progress-overlay" :style="{ width: getProgress('cpp') + '%' }"></div>
                <img src="@/assets/cpp-logo.png" alt="C++" class="course-image" draggable="false">
                <h3>C++</h3>
              </div>
            </div>
          </div>
        </section>
      </div>
    </main>
  </div>
</template>

<script>
// УБРАТЬ импорт Header отсюда
export default {
  name: 'MainPage',
  data() {
    return {
      courseProgress: {
        python: 25,
        javascript: 50,
        java: 0,
        golang: 75,
        csharp: 10,
        cpp: 100
      },
      activeCourse: ''
    }
  },
  methods: {
    goToCourse(language) {
      this.$router.push(`/courses/${language}`)
    },
    getProgress(language) {
      return this.courseProgress[language] || 0
    },
    setCourseActive(course, active) {
      this.activeCourse = active ? course : ''
    }
  },
  mounted() {
    const savedProgress = localStorage.getItem('courseProgress')
    if (savedProgress) {
      this.courseProgress = JSON.parse(savedProgress)
    }
  }
}
</script>

<style scoped>
.main-page {
  min-height: 100vh;
  background-color: #0E1117;
}

.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 20px;
}

.main {
  padding: 5rem 0;
}

.courses-grid {
  display: flex;
  flex-direction: column;
  gap: 4rem;
  width: 100%;
}

.course-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 3rem;
  width: 100%;
}

.course-card {
  background-color: #303030;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
  text-align: center;
  font-weight: 500;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  min-height: 150px;
  width: 100%;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
  transform: scale(1);
  transition: all 0.1s ease;
}

.course-card--active {
  transform: scale(0.95);
  background-color: #404040;
}

.course-card:hover {
  transform: translateY(-4px) scale(1);
  box-shadow: 0 8px 15px rgba(0,0,0,0.2);
}

.course-card--active:hover {
  transform: scale(0.95);
}

.course-card:hover .progress-overlay {
  background-color: #7a7a7a;
}

.course-card:hover {
  background-color: #3a3a3a;
}

.course-card h3 {
  margin: 0;
  font-size: 64px;
  color: #9CA3AF;
  flex: 1;
  position: relative;
  z-index: 2;
  transition: color 0.3s ease;
}

.course-card:hover h3 {
  color: #ffffff;
}

.course-image {
  width: 100px;
  height: 100px;
  margin-right: 2rem;
  object-fit: contain;
  position: relative;
  z-index: 2;
  transition: transform 0.3s ease;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
  pointer-events: none;
  -webkit-user-drag: none;
  -webkit-touch-callout: none;
}

.course-card:hover .course-image {
  transform: scale(1.05);
}

.progress-overlay {
  position: absolute;
  top: 0;
  left: 0;
  height: 100%;
  background-color: #676767;
  transition: all 0.3s ease;
  z-index: 1;
}

/* Мобильная адаптация */
@media (max-width: 1024px) {
  .course-card h3 {
    font-size: 48px;
  }
  
  .course-image {
    width: 80px;
    height: 80px;
    margin-right: 1.5rem;
  }
  
  .course-card {
    padding: 1.5rem;
    min-height: 130px;
  }
}

@media (max-width: 768px) {
  .main {
    padding: 3rem 0;
  }
  
  .courses-grid {
    gap: 2rem;
  }
  
  .course-row {
    grid-template-columns: 1fr;
    gap: 2rem;
  }
  
  .course-card h3 {
    font-size: 40px;
  }
  
  .course-image {
    width: 70px;
    height: 70px;
    margin-right: 1rem;
  }
  
  .course-card {
    padding: 1.5rem;
    min-height: 120px;
  }
}

@media (max-width: 480px) {
  .container {
    padding: 0 15px;
  }
  
  .main {
    padding: 2rem 0;
  }
  
  .courses-grid {
    gap: 1.5rem;
  }
  
  .course-row {
    gap: 1.5rem;
  }
  
  .course-card h3 {
    font-size: 32px;
  }
  
  .course-image {
    width: 60px;
    height: 60px;
    margin-right: 0.8rem;
  }
  
  .course-card {
    padding: 1.2rem;
    min-height: 100px;
  }
}

@media (max-width: 360px) {
  .course-card h3 {
    font-size: 28px;
  }
  
  .course-image {
    width: 50px;
    height: 50px;
    margin-right: 0.6rem;
  }
  
  .course-card {
    padding: 1rem;
    min-height: 90px;
  }
}
</style>