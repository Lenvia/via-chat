import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import antd from './plugin/anti-ui'

import './plugin/anti-ui'
import './plugin/http'

// 创建应用实例
const app = createApp(App)
antd(app)

// 挂载路由
app.use(router)

// 挂载应用
app.mount('#app')