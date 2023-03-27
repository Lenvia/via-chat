import { createRouter, createWebHistory } from 'vue-router'

const Login = () => import(/* webpackChunkName: "Login" */ '../views/LoginView.vue')

const routes = [

  {
    path: '/login',
    name: 'login',
    component: Login
  },

]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
