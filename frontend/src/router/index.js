import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'main',
    component: () => import('@/components/Main/Main.vue')
  },
  // Универсальный маршрут для всех курсов
  {
    path: '/courses/:language',
    name: 'course',
    component: () => import('@/components/Courses/CourseView.vue'),
    props: route => ({ language: route.params.language }),
    beforeEnter: (to, from, next) => {
      // Проверяем поддерживаемые языки (опционально)
      const supportedLanguages = [
        'python', 'javascript', 'typescript', 
        'java', 'cpp', 'csharp', 'go', 'rust',
        'php', 'ruby'
      ]
      
      const requestedLang = to.params.language.toLowerCase()
      
      if (supportedLanguages.includes(requestedLang)) {
        next()
      } else {
        console.warn(`Язык ${requestedLang} не поддерживается`)
        next('/')
      }
    }
  },
  // Старые маршруты для обратной совместимости (редирект)
  {
  path: '/courses/python',
  redirect: { name: 'course', params: { language: 'python' } }
  },
  {
    path: '/courses/javascript',
    redirect: { name: 'course', params: { language: 'javascript' } }
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
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('@/components/Settings/Settings.vue')
  },
  // Маршруты для учителей
  {
    path: '/teacher/statistics',
    name: 'teacher-statistics',
    component: () => import('@/components/Admin/Statistics.vue'),
    meta: { requiresAuth: true, requiresTeacher: true }
  },
  {
    path: '/teacher/tasks',
    name: 'teacher-tasks',
    component: () => import('@/components/Admin/Tasks.vue'), // У вас уже есть
    meta: { requiresAuth: true, requiresTeacher: true }
  },
  {
    path: '/teacher/tasks/create',
    name: 'teacher-create-task',
    component: () => import('@/components/Admin/Tasks.vue'), // Тот же компонент
    meta: { requiresAuth: true, requiresTeacher: true }
  },
  {
    path: '/teacher/tasks/:id/edit',
    name: 'teacher-edit-task',
    component: () => import('@/components/Admin/Tasks.vue'), // Тот же компонент
    meta: { requiresAuth: true, requiresTeacher: true },
    props: route => ({ editTaskId: route.params.id }) // Передаем ID для редактирования
  },
  // 404 страница
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: () => import('@/components/NotFound/NotFound.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Глобальная навигационная защита
router.beforeEach((to, from, next) => {
  console.log('Navigation:', from.name, '->', to.name)
  
  // Проверяем требуется ли авторизация
  if (to.meta.requiresAuth) {
    const user = JSON.parse(localStorage.getItem('user') || '{}')
    
    if (!user.isLoggedIn) {
      next('/auth')
      return
    }
    
    if (to.meta.requiresTeacher && user.role !== 'teacher') {
      alert('Только учителя имеют доступ к этой странице!')
      next('/')
      return
    }
  }
  
  next()
})

export default router