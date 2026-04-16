import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { initAuth } from './composables/useAuth'

initAuth().then(() => {
  createApp(App).use(router).mount('#app')
})


