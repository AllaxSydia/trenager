<template>
  <div class="not-found">
    <div class="container">
      <h1>404</h1>
      <h2>Страница не найдена</h2>
      <p>Запрошенная страница не существует или была перемещена.</p>
      <router-link to="/" class="btn btn-primary">
        Вернуться на главную
      </router-link>
    </div>
  </div>
</template>

<script>
export default {
  name: 'NotFound'
}
</script>

<style scoped>
.not-found {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 70vh;
  text-align: center;
  padding: 20px;
}

.container {
  max-width: 500px;
  padding: 2rem;
  background: rgba(30, 41, 59, 0.8);
  border-radius: 12px;
  border: 1px solid #334155;
}

h1 {
  font-size: 6rem;
  color: #3B82F6;
  margin: 0;
  line-height: 1;
}

h2 {
  font-size: 2rem;
  margin: 1rem 0;
  color: #f8fafc;
}

p {
  color: #94a3b8;
  margin-bottom: 2rem;
}

.btn-primary {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  background: #3B82F6;
  color: white;
  text-decoration: none;
  border-radius: 6px;
  font-weight: 500;
  transition: background 0.2s;
}

.btn-primary:hover {
  background: #2563EB;
}
</style>