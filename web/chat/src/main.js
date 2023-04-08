import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'

import './assets/css/style.css'

// 创建应用实例
const app = createApp(App)

const getUrl = () => {
    const host = window.location.hostname;
    const protocol = window.location.protocol;
    return `${protocol}//${host}:8322`;
};


const Url = getUrl()
axios.defaults.baseURL = Url
app.config.globalProperties.$http = axios
app.config.globalProperties.$http.defaults.withCredentials = true
app.provide('Url', Url)

// 挂载路由
app.use(router)


// 挂载应用
app.mount('#app')

export default app;