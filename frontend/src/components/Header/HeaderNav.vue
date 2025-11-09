<template>
  <nav class="nav">
    <router-link to="/rating" class="nav-link" active-class="nav-link--active">
      Рейтинг
    </router-link>
    <router-link to="/profile" class="nav-link" active-class="nav-link--active">
      Профиль
    </router-link>
  </nav>
</template>

<script>
export default {
  name: 'HeaderNav'
}
</script>

<style scoped>
.nav {
  display: flex;
  gap: 32px;
  align-items: center;
}

.nav-link {
  width: 200px;
  font-size: 18px;
  text-align: center;
  transition: all 0.2s;
  color: #9CA3AF;
  text-decoration: none;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
}

.nav-link:hover {
  color: #E5E7EB;
  background: #1a1f2e;
}

.nav-link--active {
  color: #E5E7EB;
  background: #3B82F6;
}

.nav-link:active {
  transform: scale(0.95);
  background: #2563EB;
}

@media (max-width: 1023px) {
  .nav {
    gap: 20px;
  }

  .nav-link {
    width: 150px;
    font-size: 16px;
  }
}

@media (max-width: 767px) {
  .nav {
    display: none;
  }
}
</style>