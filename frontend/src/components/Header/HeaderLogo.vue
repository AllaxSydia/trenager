<template>
  <router-link to="/" class="logo-link">
    <div class="logo">
      CodeForge
    </div>
  </router-link>
</template>

<style scoped>
.logo-link {
  text-align: center;
  display: flex;
  align-items: center;
  text-decoration: none;
  flex-shrink: 0;
}

.logo {
  font-size: 36px;
  font-weight: bold;
  width: 330px;
  color: #E5E7EB;
  transition: all 0.2s ease;
  cursor: pointer;
}

.logo-link:hover .logo {
  color: #ffffff;
  transform: scale(1.02);
}

.logo-link:active .logo {
  transform: scale(0.98);
}

@media (max-width: 1023px) {
  .logo {
    font-size: 28px;
    width: 250px;
  }
}

@media (max-width: 767px) {
  .logo {
    font-size: 24px;
    width: auto;
  }
}

@media (max-width: 480px) {
  .logo {
    font-size: 20px;
  }
}
</style>