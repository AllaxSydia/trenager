<template>
  <div id="app">
    <h1>Code Trainer</h1>
    <router-view />
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
}
</style>