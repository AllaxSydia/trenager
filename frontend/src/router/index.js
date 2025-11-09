import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'main',
    component: () => import('@/components/Main/Main.vue')
  },
  {
    path: '/courses/python',
    name: 'python',
    component: () => import('@/components/Courses/Python.vue')
  },
  {
    path: '/courses/javascript',
    name: 'javascript',
    component: () => import('@/components/Courses/JavaScript.vue')
  },
  {
    path: '/auth',
    name: 'auth',
    component: () => import('@/components/Auth/AuthPage.vue')
  },
  {
    path: '/rating',
    name: 'rating',
    component: () => import('@/components/Rating/Rating.vue')
  },
  {
    path: '/profile',
    name: 'profile', 
    component: () => import('@/components/Profile/Profile.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Добавим глобальную навигационную защиту
router.beforeEach((to, from, next) => {
  console.log('Navigation:', from.name, '->', to.name)
  next()
})

export default router