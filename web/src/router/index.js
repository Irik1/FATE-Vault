import { createRouter, createWebHistory } from 'vue-router'
import CharactersList from '../views/CharactersList.vue'
import CharacterDetail from '../views/CharacterDetail.vue'
import Stunts from '../views/Stunts.vue'
import User from '../views/User.vue'

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
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router


