import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { public: true },
  },
  {
    path: '/',
    name: 'Dashboard',
    component: () => import('../views/Dashboard.vue'),
  },
  {
    path: '/articles',
    name: 'ArticleList',
    component: () => import('../views/ArticleList.vue'),
  },
  {
    path: '/articles/new',
    name: 'ArticleCreate',
    component: () => import('../views/ArticleEdit.vue'),
  },
  {
    path: '/articles/:id/edit',
    name: 'ArticleEdit',
    component: () => import('../views/ArticleEdit.vue'),
  },
  {
    path: '/categories',
    name: 'Categories',
    component: () => import('../views/CategoryMgmt.vue'),
  },
  {
    path: '/tags',
    name: 'Tags',
    component: () => import('../views/TagMgmt.vue'),
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('../views/Settings.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Auth guard
router.beforeEach(async (to) => {
  if (to.meta.public) return true

  const { useAuthStore } = await import('../stores/auth')
  const auth = useAuthStore()

  if (!auth.authenticated) {
    const isAuth = await auth.checkAuth()
    if (!isAuth) {
      return { name: 'Login' }
    }
  }

  return true
})

export default router
