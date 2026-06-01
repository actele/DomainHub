import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: '/dashboard' },
    { path: '/login', name: 'login', component: () => import('../views/Login.vue') },
    { path: '/register', name: 'register', component: () => import('../views/Register.vue') },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('../views/Dashboard.vue'),
      meta: { requiresAuth: true },
      redirect: '/dashboard/domains',
      children: [
        { path: 'domains', name: 'domains', component: () => import('../views/Domains.vue') },
        { path: 'domains/:id/records', name: 'domain-records', component: () => import('../views/DomainRecords.vue') },
        { path: 'provider-keys', name: 'provider-keys', component: () => import('../views/ProviderKeys.vue') },
        { path: 'admin/users', name: 'admin-users', component: () => import('../views/admin/Users.vue'), meta: { requiresAdmin: true } },
        { path: 'admin/providers', name: 'admin-providers', component: () => import('../views/admin/Providers.vue'), meta: { requiresAdmin: true } }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  if (to.matched.some(r => r.meta.requiresAuth)) {
    if (!token) {
      next({ path: '/login', query: { redirect: to.fullPath } })
    } else {
      if (to.matched.some(r => r.meta.requiresAdmin) && localStorage.getItem('role') !== 'admin') {
        next({ path: '/dashboard/domains' })
      } else {
        next()
      }
    }
  } else {
    token && to.path === '/login' ? next({ path: '/dashboard' }) : next()
  }
})

export default router
