import { createRouter, createWebHistory } from 'vue-router'
import Home from '../pages/Home.vue'
import Docs from '../pages/Docs.vue'
import Playground from '../pages/Playground.vue'
import Contributors from '../pages/Contributors.vue'

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/docs', name: 'Docs', component: Docs },
  { path: '/playground', name: 'Playground', component: Playground },
  { path: '/contributors', name: 'Contributors', component: Contributors },
]

const router = createRouter({
  history: createWebHistory('/pam/'),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) {
      return { el: to.hash, behavior: 'smooth' }
    }
    return savedPosition || { top: 0 }
  },
})

export default router
