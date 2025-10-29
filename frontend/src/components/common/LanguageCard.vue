<template>
  <div class="language-card" @click="$emit('select')">
    <div class="icon">{{ language.icon }}</div>
    <h3>{{ language.name }}</h3>
  </div>
</template>

<script setup>
defineProps({
  language: {
    type: Object,
    required: true
  }
})

defineEmits(['select'])
</script>

<style scoped>
.language-card {
  background: white;
  border-radius: 12px;
  padding: 30px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

.language-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 20px rgba(0,0,0,0.15);
}

.icon {
  font-size: 3rem;
  margin-bottom: 15px;
}

h3 {
  color: #333;
  font-size: 1.3rem;
}
</style>