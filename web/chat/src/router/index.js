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
// router 前置守卫
// router.beforeEach((to, from, next) => {
//   const userToken = window.sessionStorage.getItem('token')
//   if (to.path === '/login') return next()
//   if (!userToken) {
//     next('/login')
//   } else {
//     next()
//   }
// })

export default router
