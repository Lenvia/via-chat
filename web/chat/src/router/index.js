import { createRouter, createWebHistory } from 'vue-router'

const Login = () => import(/* webpackChunkName: "Login" */ '../views/LoginView.vue')
const Home = () => import(/* webpackChunkName: "Login" */ '../views/HomeView.vue')

const routes = [

  {
    path: '/login',
    name: 'login',
    component: Login
  },
  {
    path: '/home',
    name: 'home',
    component: Home
  },

]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
