export default [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/HomeView.vue')
  },
  {
    path: '/language/:lang',
    name: 'language',
    component: () => import('@/views/LanguageView.vue'),
    props: true
  },
  {
    path: '/task/:lang/:topic/:taskId',
    name: 'task',
    component: () => import('@/views/TaskView.vue'),
    props: true
  },
  {
    path: '/sandbox',
    name: 'sandbox',
    component: () => import('@/views/SandboxView.vue')
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
]