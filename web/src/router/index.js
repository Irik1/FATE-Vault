import { createRouter, createWebHistory } from 'vue-router'
import CharactersList from '../views/CharactersList.vue'
import CharacterDetail from '../views/CharacterDetail.vue'
import Stunts from '../views/Stunts.vue'
import User from '../views/User.vue'
import Login from '../views/Login.vue'
import Register from '../views/Register.vue'
import { user } from '../composables/useAuth'

const routes = [
  {
    path: '/',
    redirect: '/characters'
  },
  {
    path: '/characters',
    name: 'Characters',
    component: CharactersList
  },
  {
    path: '/characters/new',
    name: 'CharacterNew',
    component: CharacterDetail,
    props: true
  },
  {
    path: '/characters/:id',
    name: 'CharacterDetail',
    component: CharacterDetail,
    props: true
  },
  {
    path: '/stunts',
    name: 'Stunts',
    component: Stunts
  },
  {
    path: '/user',
    name: 'User',
    component: User
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { guestOnly: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { guestOnly: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, _from, next) => {
  if (to.meta.guestOnly && user.value) {
    next({ name: 'User' })
    return
  }
  next()
})

export default router


