<template>
  <div id="app">
    <Header v-if="showHeader" />
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script>
import Header from '@/components/Header/Header.vue'

export default {
  name: 'App',
  components: {
    Header
  },
  computed: {
    showHeader() {
      return this.$route.name !== 'auth'
    }
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  background-color: #0E1117;
  font-family: Arial, sans-serif;
}

/* Добавьте эти стили */
.main-content {
  margin-top: 96px; /* Высота header для десктопа */
}

@media (max-width: 1023px) {
  .main-content {
    margin-top: 80px; /* Высота header для планшетов */
  }
}

@media (max-width: 767px) {
  .main-content {
    margin-top: 70px; /* Высота header для мобильных */
  }
}

@media (max-width: 480px) {
  .main-content {
    margin-top: 60px; /* Высота header для маленьких мобильных */
  }
}
</style>