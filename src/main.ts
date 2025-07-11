import 'bootstrap/dist/css/bootstrap.css' // Bootstrap 5 Styles
import '@fortawesome/fontawesome-free/css/all.css' // All FA free icons

import { createApp, h } from 'vue'

import App from './app.vue'
import router from './router.ts'

const app = createApp({
  name: 'GoLatestVer',
  render() {
    return h(App)
  },
})

app.use(router)
app.mount('#app')
