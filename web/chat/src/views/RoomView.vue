<template>
  <div class="chat-room">
<!--    <ul class="user-info">-->
<!--      <li>User ID: {{ userInfo.uid }}</li>-->
<!--      <li>User Name: {{ userInfo.username }}</li>-->
<!--      <li>User Avatar ID: {{ userInfo.avatar_id }}</li>-->
<!--    </ul>-->
    <div class="chat-messages" ref="msgContainer">
      <div v-for="msg in msgList" :key="msg.id" :class="getMessageClass(msg)">
        <div>{{ msg.content }}</div>
        <div class="message-info">
          <div class="message-sender">{{ getUsername(msg) }}</div>
<!--          <div class="message-time">{{ formatTime(msg.CreatedAt) }}</div>-->
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {defineComponent, ref, onMounted, nextTick} from 'vue';
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
    const msgContainer = ref(null);




    // 获取历史消息
    async function loadHistoryAndBuildWS() {
      const response = await app.config.globalProperties.$http.get('/room/'+room_id);
      console.log(response);
      if (response.status === 200) {
        const data = await response.data;
        userInfo.value = data.user_info;
        msgList.value = data.msg_list;
        console.log(msgList.value)
        msgListCount.value = data.msg_list_count;
        
        console.log(userInfo.value);
        WebSocketConnect(userInfo.value, room_id)
        
      } else {
        ElMessage.error('error')
      }

      // 等待DOM更新并滚动到底部
      nextTick(() => {
        msgContainer.value.scrollTop = msgContainer.value.scrollHeight;
      });
    }
    onMounted(() => {
      loadHistoryAndBuildWS();
    });

    function WebSocketConnect(userInfo, room_id, toUserInfo = null) {
      if ("WebSocket" in window) {
        if (userInfo.uid <= 0) {
          alert('参数错误，请刷新页面重试');
          return false;
        }

        let send_data = JSON.stringify({
          "status": toUserInfo ? 5 : 1,
          "data": {
            "uid": userInfo.uid.toString(),
            "room_id": room_id,
            "avatar_id": userInfo.avatar_id,
            "username": userInfo.username,
            "to_user": toUserInfo
          }
        })

        const ws = new WebSocket(`ws://localhost:8322/ws`); // 连接 WebSocket

        ws.onopen = function () {
          console.log(send_data);
          ws.send(send_data);
          //console.log("send_data 发送数据", send_data)
        };

        ws.onmessage = function (evt) {
          // console.log(evt)
          let received_msg = JSON.parse(evt.data);
          console.log("数据已接收...", received_msg);

          let myDate = new Date();
          let time = myDate.toLocaleDateString() + " " + myDate.toLocaleTimeString()


          let systemInfo;
          switch (received_msg.status) {
            case 1:
              systemInfo =`<li class="systeminfo"><span>`
                  +`【` + received_msg.data.username + `】` + time + " 加入了房间" +`</span></li>`;
              msgContainer.value.innerHTML += systemInfo;
          }
          nextTick(() => {
            msgContainer.value.scrollTop = msgContainer.value.scrollHeight;
          });
        };

        ws.onclose = function () {
          console.log("serve 连接已关闭... ");
          // console.log(c);
        };

        ws.onerror = function (evt) {
          ws.close()
          console.log("触发 onerror", evt)
        }

      } else {
        // 浏览器不支持 WebSocket
        alert("您的浏览器不支持 WebSocket!");
      }
    }

    // Get the sender username for a message
    function getUsername(message) {
      if (message.user_id === userInfo.value.uid) {
        return 'You';
      }
      return message.user_id;
    }

    // Get the class for a message element (based on whether it is from the current user or not)
    function getMessageClass(message) {
      if (message.user_id === userInfo.value.uid) {
        return 'message-sent';
      }
      return 'message-received';
    }

    return {
      userInfo,
      msgList,
      msgListCount,
      msgContainer,
      getUsername,
      getMessageClass,
    };
  },
});
</script>

<style>
.chat-room {
  height: 100%;
}

.chat-messages {
  position: absolute;
  left: 50%;
  width: 50%;
  height: calc(100% - 50px);
  overflow-y: auto;
  transform: translate(-50%, 0);
}

.chat-messages > div {
  margin-bottom: 10px;
}

.message-sent {
  text-align: right;
  background-color: #4caf50;
  color: #fff;
}

.message-received {
  text-align: left;
  background-color: #f5f5f5;
  color: #000;
}

.message-info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
}

.message-sender {
  font-weight: bold;
  margin-right: 10px;
}

.message-time {
  color: #999;
}
</style>