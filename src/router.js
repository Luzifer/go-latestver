import CatalogEntry from './catalog_entry.vue'
import CatalogIndex from './catalog_index.vue'
import Log from './log.vue'
import VueRouter from 'vue-router'

const router = new VueRouter({
  mode: 'history',
  routes: [
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
  ],
})

export default router
