<template>
  <div id="app">
    <div class="app-container">
      <h1>Code Trainer</h1>
      <router-view />
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
  mounted() {
    console.log('✅ App.vue mounted')
  }
}
</script>

<style>
#app {
  font-family: Arial, sans-serif;
  padding: 20px;
  width: 100%;
  min-height: 100vh;
  background: #f5f5f5;
}

/* ФИКС: Главный контейнер с ограничением ширины */
.app-container {
  max-width: 1400px;
  margin: 0 auto;
  width: 100%;
}
</style>