/* eslint-disable sort-imports */

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'
import 'bootswatch/dist/darkly/bootstrap.css'

import Vue from 'vue'
import { BootstrapVue } from 'bootstrap-vue'

import App from './app.vue'

Vue.config.devtools = process.env.NODE_ENV === 'dev'
Vue.use(BootstrapVue)

new Vue({
  components: { App },
  el: '#app',
  name: 'GoLatestVer',
  render: h => h(App),
})
