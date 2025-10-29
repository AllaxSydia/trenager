<template>
  <div class="language-view">
    <div class="container">
      <header class="header">
        <button class="back-btn" @click="$router.push('/')">← Назад</button>
        <div class="language-info">
          <span class="language-icon">{{ getLanguageIcon($route.params.lang) }}</span>
          <h1>{{ getLanguageName($route.params.lang) }}</h1>
          <p>Изучайте {{ getLanguageName($route.params.lang) }} через практические задачи</p>
        </div>
      </header>

      <div class="content">
        <div class="topics-container">
          <div 
            v-for="topic in availableTopics" 
            :key="topic.key"
          >
            <TopicSection
              :language="$route.params.lang"
              :topic="topic.key"
              :tasks="topic.tasks"
            />
          </div>
        </div>

        <div v-if="availableTopics.length === 0" class="no-tasks">
          <h3>📝 Задачи в разработке</h3>
          <p>Задачи для языка "{{ getLanguageName($route.params.lang) }}" скоро будут добавлены</p>
          <div style="margin-top: 20px; text-align: left;">
            <h4>Доступные языки с задачами:</h4>
            <ul>
              <li v-for="lang in languagesWithTasks" :key="lang.id">
                {{ lang.icon }} {{ lang.name }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import TopicSection from '@/components/tasks/TopicSection.vue'
import { LANGUAGES, TASK_DATA, TOPICS } from '@/utils/constants'
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

// Доступные темы с задачами для текущего языка
const availableTopics = computed(() => {
  const languageTasks = TASK_DATA[route.params.lang] || {}
  const topics = []

  // Проходим по всем возможным темам
  Object.entries(TOPICS).forEach(([topicKey, topicInfo]) => {
    const tasks = languageTasks[topicKey] || []
    if (tasks.length > 0) {
      topics.push({
        key: topicKey,
        name: topicInfo.name,
        tasks: tasks
      })
    }
  })

  // Сортируем по порядку
  return topics.sort((a, b) => TOPICS[a.key].order - TOPICS[b.key].order)
})

// Языки которые имеют задачи
const languagesWithTasks = computed(() => {
  return LANGUAGES.filter(lang => TASK_DATA[lang.id])
})

const getLanguageIcon = (langId) => {
  return LANGUAGES.find(l => l.id === langId)?.icon || '💻'
}

const getLanguageName = (langId) => {
  return LANGUAGES.find(l => l.id === langId)?.name || langId
}
</script>

<style scoped>
.language-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #74b9ff 0%, #0984e3 100%);
  padding: 30px 20px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  margin-bottom: 40px;
}

.back-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: none;
  border-radius: 8px;
  padding: 10px 15px;
  cursor: pointer;
  backdrop-filter: blur(10px);
  margin-bottom: 20px;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.language-info {
  text-align: center;
  color: white;
}

.language-icon {
  font-size: 4rem;
  display: block;
  margin-bottom: 10px;
}

.language-info h1 {
  font-size: 2.5rem;
  margin-bottom: 10px;
}

.language-info p {
  font-size: 1.2rem;
  opacity: 0.9;
}

.content {
  background: white;
  border-radius: 15px;
  padding: 40px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
}

.topics-container {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.no-tasks {
  text-align: center;
  padding: 60px 40px;
  color: #666;
}

.no-tasks h3 {
  font-size: 1.5rem;
  margin-bottom: 15px;
}

.no-tasks p {
  font-size: 1.1rem;
  opacity: 0.8;
}

.no-tasks ul {
  list-style: none;
  padding: 0;
  margin-top: 10px;
}

.no-tasks li {
  padding: 5px 0;
  font-size: 1.1rem;
}
</style>