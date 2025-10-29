<template>
  <div class="home-view">
    <div class="container">
      <header class="header">
        <h1>🚀 Платформа для обучения программированию</h1>
        <p>Выберите язык программирования чтобы начать</p>
      </header>

      <div class="languages-grid">
        <LanguageCard 
          v-for="lang in languages"
          :key="lang.id"
          :language="lang"
          @select="selectLanguage(lang.id)"
        />
      </div>

      <div class="sandbox-section">
        <button class="sandbox-btn" @click="$router.push('/sandbox')">
          🏖️ Режим песочницы
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import LanguageCard from '@/components/common/LanguageCard.vue'
import { LANGUAGES } from '@/utils/constants'
import { useRouter } from 'vue-router'

const router = useRouter()
const languages = LANGUAGES

const selectLanguage = (langId) => {
  router.push(`/language/${langId}`)
}
</script>

<style scoped>
.home-view {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 40px 20px;
}

.header {
  text-align: center;
  margin-bottom: 60px;
  color: white;
}

.header h1 {
  font-size: 2.5rem;
  margin-bottom: 15px;
}

.header p {
  font-size: 1.2rem;
  opacity: 0.9;
}

.languages-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 25px;
  margin-bottom: 40px;
}

.sandbox-section {
  text-align: center;
}

.sandbox-btn {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  border: 2px solid white;
  border-radius: 12px;
  padding: 15px 30px;
  font-size: 1.1rem;
  cursor: pointer;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.sandbox-btn:hover {
  background: white;
  color: #667eea;
}
</style>