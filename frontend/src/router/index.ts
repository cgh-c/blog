import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue'),
  },
  {
    path: '/article/:slug',
    name: 'Article',
    component: () => import('../views/Article.vue'),
  },
  {
    path: '/category/:id',
    name: 'Category',
    component: () => import('../views/Category.vue'),
  },
  {
    path: '/archive',
    name: 'Archive',
    component: () => import('../views/Archive.vue'),
  },
  {
    path: '/about',
    name: 'About',
    component: () => import('../views/About.vue'),
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior() {
    return { top: 0 }
  },
})

export default router
