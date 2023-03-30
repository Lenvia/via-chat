<template>
  <div>
    <ul class="user-info">
      <li>User ID: {{ userInfo.uid }}</li>
      <li>User Name: {{ userInfo.username }}</li>
      <li>User Avatar ID: {{ userInfo.avatar_id }}</li>
    </ul>
    <ul class="message-list">
      <li v-for="msg in msgList" :key="msg.ID">
        <div>{{ msg.Content }}</div>
        <div>{{ msg.CreatedAt }}</div>
      </li>
    </ul>
  </div>
</template>

<script>
import {defineComponent, ref, onMounted, watch} from 'vue';
import app from "@/main";
import {ElMessage} from "element-plus";
import router from "@/router";

export default defineComponent({
  name: 'RoomView',
  setup() {
    const userInfo = ref({ uid: '', username: '', avatar_id: '' });
    const msgList = ref([]);
    const msgListCount = ref(0);
    const room_id = router.currentRoute.value.params.room_id;


    // 监听 $route 变化
    // watch(router, (to, from) => {
    //   if (to.params.room_id !== roomId.value) {
    //     roomId.value = to.params.room_id;
    //     // TODO: 根据 roomId 获取房间信息以及聊天记录等
    //     console.log(roomId.value)
    //   }
    // });

    async function loadData() {
      console.log("room_id:" + room_id);
      const response = await app.config.globalProperties.$http.get('/room/'+room_id);
      console.log(response);
      // if (response.status === 200) {
      //   const data = await response.data;
      //   userInfo.value = data.user_info;
      //   msgList.value = data.msg_list;
      //   msgListCount.value = data.msg_list_count;
      // } else {
      //   ElMessage.error('error')
      // }
    }
    onMounted(() => {
      loadData();
    });

    return {
      userInfo,
      msgList,
      msgListCount,
    };
  },
});
</script>