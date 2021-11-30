/* eslint-disable sort-imports */

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import 'bootswatch/dist/darkly/bootstrap.css'

import { library } from '@fortawesome/fontawesome-svg-core'
import { fab } from '@fortawesome/free-brands-svg-icons'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

library.add(fab, fas)

import Vue from 'vue'
import { BootstrapVue } from 'bootstrap-vue'
import VueRouter from 'vue-router'

import App from './app.vue'
import router from './router.js'

Vue.config.devtools = process.env.NODE_ENV === 'dev'
Vue.component('FontAwesomeIcon', FontAwesomeIcon)
Vue.use(BootstrapVue)
Vue.use(VueRouter)

new Vue({
  components: { App },
  el: '#app',
  name: 'GoLatestVer',
  render: h => h(App),
  router,
})
