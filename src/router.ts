import { createRouter, createWebHistory } from 'vue-router'

import CatalogEntry from './catalog_entry.vue'
import CatalogIndex from './catalog_index.vue'
import Log from './log.vue'

const routes = [
  {
    component: CatalogIndex,
    name: 'index',
    path: '/',
  },
  {
    component: CatalogEntry,
    name: 'entry',
    path: '/:name/:tag',
  },
  {
    component: Log,
    name: 'log',
    path: '/log',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
