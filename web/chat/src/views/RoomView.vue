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
import { defineComponent, ref, onMounted } from 'vue';

export default defineComponent({
  name: 'RoomView',
  setup() {
    const room_id = 1; // 你的房间ID
    const userInfo = ref({ uid: '', username: '', avatar_id: '' });
    const msgList = ref([]);
    const msgListCount = ref(0);

    async function loadData() {
      const res = await fetch(`/api/chat/${room_id}`);
      const data = await res.json();
      userInfo.value = data.user_info;
      msgList.value = data.msg_list;
      msgListCount.value = data.msg_list_count;
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