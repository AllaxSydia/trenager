<template>
  <div class="rating-page">
    <!-- УБРАТЬ Header отсюда -->
    <main class="main">
      <div class="container">
        <h1>Страница рейтинга</h1>
        <p>Здесь будет рейтинг пользователей</p>
      </div>
    </main>
  </div>
</template>

<script>
// УБРАТЬ импорт Header отсюда
export default {
  name: 'Rating'
}
</script>

<style scoped>
.rating-page {
  min-height: 100vh;
  background-color: #0E1117;
  color: white;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 20px;
}
</style>