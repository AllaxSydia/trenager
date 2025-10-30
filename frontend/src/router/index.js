import { createRouter, createWebHashHistory } from 'vue-router'
import routes from './routes'

const router = createRouter({
  history: createWebHashHistory(), // ИЗМЕНИТЕ ЭТУ СТРОКУ
  routes
})

export default router