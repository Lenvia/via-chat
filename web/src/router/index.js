import { createRouter, createWebHistory } from 'vue-router'

const Login = () => import(/* webpackChunkName: "Login" */ '../views/LoginView.vue')
const Home = () => import(/* webpackChunkName: "Home" */ '../views/HomeView.vue')
const Room = () => import(/* webpackChunkName: "Room" */ '../views/RoomView.vue')

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
  {
    path: '/room/:room_id',
    name: 'room',
    component: Room
  }

]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// router 前置守卫
router.beforeEach((to, from, next) => {
  const userToken = window.sessionStorage.getItem('token')
  // console.log(userToken)
  if (!userToken && to.path !== '/login') {
    next('/login')
  } else {
    next()
  }
})

export default router
