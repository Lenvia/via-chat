<template>
  <div class="login-container">
    <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
      <el-form-item label="用户名" prop="username">
        <el-input v-model="formData.username"></el-input>
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input type="password" v-model="formData.password" @keydown.enter.prevent="handleEnter"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSubmit">登录</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { ElForm, ElFormItem, ElInput, ElButton, ElMessage } from 'element-plus'
import router from "@/router";
import app from '@/main.js';

export default {
  name: 'LoginView',
  components: {
    ElForm,
    ElFormItem,
    ElInput,
    ElButton
  },
  setup() {
    const formRef = ref(null)
    const formData = reactive({
      username: '',
      password: ''
    })

    const rules = {
      username: [
        { required: true, message: '请输入用户名', trigger: 'blur' }
      ],
      password: [
        { required: true, message: '请输入密码', trigger: 'blur' },
        { min: 6, message: '密码长度不能少于 6 位', trigger: 'blur' }
      ]
    }


    const handleEnter = (event) => {
      // 处理 Enter 键按下事件
      handleSubmit();
    };

    const handleSubmit = async () => {
      try {
        const form = formRef.value // 直接从 ref 对象中获取组件引用
        // 通过 await 表达式等待校验结果和登录结果
        await form.validate()
        const response = await app.config.globalProperties.$http.post('/login', {
          username: formData.username,
          password: formData.password,
          avatar_id: "1",  // 临时的
        })
        console.log(response.data)

        if (response.data.code === 0) {
          window.sessionStorage.setItem('token', response.data.token)
          await router.push('/home')
        } else if(response.data.code === 5000) {
          ElMessage({
            message: response.data.message,
            type: 'error'
          });
        }
        else{
          ElMessage({
            message: '登录失败',
            type: 'error'
          });
        }
      } catch (error) {
        console.error(error)
      }
    }


    return {
      formRef,  // 将 formRef 渲染到模板中
      formData,
      rules,
      handleSubmit,
      handleEnter,
    }
  },

}
</script>

<style>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}
</style>