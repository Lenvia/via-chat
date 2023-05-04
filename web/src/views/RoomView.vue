<template>
  <div class="chat-room">
    <div class="chat-messages" ref="msgContainer">
      <div v-for="msg in msgList" :key="msg.id" :class="msg.type === 'system' ? 'system-info' : getMessageClass(msg)">

        <div v-if="msg.type !== 'system'" :class="getUsernameClass(msg)">
          <div>{{ getUsername(msg) }}</div>
          <!--<div class="message-time">{{ formatTime(msg.CreatedAt) }}</div>-->
        </div>
        <div>{{ msg.content }}</div>
      </div>
    </div>

    <div class="textarea-wrapper">
      <textarea class="chat-input" ref="sendContent" @keydown.enter.prevent="handleEnter"></textarea>
      <button class="send-button" ref="sendButton">Send</button>
    </div>

  </div>
</template>

<script>
import {defineComponent, ref, onMounted, nextTick, watch} from 'vue';
import app from "@/main";
import {ElMessage} from "element-plus";
import router from "@/router";

export default defineComponent({
  name: 'RoomView',
  setup() {
    var ws;
    const userInfo = ref({ uid: '', username: '', avatar_id: '' });
    const msgList = ref([]);
    const msgListCount = ref(0);
    const room_id = router.currentRoute.value.params.room_id;
    const msgContainer = ref(null);
    const sendButton = ref(null);
    const sendContent = ref(null);

    const handleClick = () => {
      // 获取 textarea 内容
      const content = sendContent.value.value;

      // 在这里执行发送消息的逻辑
      if (content !== '') {
        let myDate = new Date();
        let time = myDate.toLocaleDateString() + myDate.toLocaleTimeString();

        let send_data = JSON.stringify({
          "status": 3,
          "data": {
            "uid": userInfo.value.uid.toString(),
            "room_id": room_id,
            "avatar_id": userInfo.value.avatar_id,
            "username": userInfo.value.username,
            "to_user": null,
            "content": content,
            "to_uid": "0",
          }
        })
        console.log("send_data", send_data)
        ws.send(send_data);
      }
      // 清空输入框
      sendContent.value.value = '';
    };

    const handleEnter = (event) => {
      // 处理 Enter 键按下事件
      if (event.metaKey || event.ctrlKey || event.shiftKey) { // 检查是否按下 Command 或 Ctrl 键
        // 换行逻辑，可以根据需要自行实现
        event.preventDefault(); // 阻止默认的换行行为
        const content = sendContent.value.value;
        sendContent.value.value = content + "\n"; // 在textarea中添加换行符
      } else {
        // 发送消息逻辑
        handleClick();
      }
    };



    // 获取历史消息
    async function loadHistoryAndBuildWS() {
      const token = window.sessionStorage.getItem('token');
      const response = await app.config.globalProperties.$http.get('/room/'+room_id, {
        headers: {
          Authorization: `Bearer ${token}`
        }
      });

      // console.log(response);
      if (response.status === 200) {
        const data = await response.data;
        userInfo.value = data.user_info;
        if(data.msg_list !== null)  // 空消息不赋值
          msgList.value = data.msg_list;
        msgListCount.value = data.msg_list_count;
        

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
      sendButton.value.addEventListener('click', handleClick); // 绑定按钮的点击事件
    });

    function WebSocketConnect(userInfo, room_id, toUserInfo = null) {
        const host = window.location.hostname;
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

      console.log(host);

        ws = new WebSocket(`ws://${host}:8322/ws`); // 连接 WebSocket
        // console.log(ws)

        ws.onopen = function () {
          ws.send(send_data);
          console.log("send_data 发送数据", send_data)
        };

        ws.onmessage = function (evt) {
          // console.log(evt)
          let received_msg = JSON.parse(evt.data);
          console.log("数据已接收...", received_msg);

          let myDate = new Date();
          let time = myDate.toLocaleDateString() + " " + myDate.toLocaleTimeString()


          let systemInfo;
          let newMsg;
          switch (received_msg.status) {
            // WARNING:
            // 对于case 1， 2 不要使用 msgContainer.value.innerHTML 直接操作 DOM！！！
            // 而是通过响应式数据来更新 DOM 可以在这些消息中添加一个额外的属性，例如 type，用于区分普通消息和系统消息。
            // 然后根据此属性在模板中使用不同的样式和显示方式。
            case 1:
              newMsg = {
                type: 'system',
                content: `【${received_msg.data.username}】${time} 加入了房间`,
              };
              msgList.value.push(newMsg);
              break;
            case 2:
              newMsg = {
                type: 'system',
                content: `【${received_msg.data.username}】${time} 离开了房间`,
              };
              msgList.value.push(newMsg);
              break;
            case 3:
              // 因为不是重新请求整个msgList，所以需要做一些小小的转换
              newMsg = {
                "avatar_id": received_msg.data.avatar_id,
                "content": received_msg.data.content,
                "created_at": received_msg.data.created_at,
                "id": received_msg.data.id,
                "image_url": received_msg.data.image_url,
                "room_id": parseInt(received_msg.data.room_id),
                "to_user_id": parseInt(received_msg.data.to_uid),
                "user_id": parseInt(received_msg.data.uid),
                "username": received_msg.data.username,
              };
              msgList.value.push(newMsg);
              break;
            case -1:
              ws.close() // 主动close掉
              console.log("client 连接已关闭...");
              break;
          }

          nextTick(() => {
            msgContainer.value.scrollTop = msgContainer.value.scrollHeight;
          });
        };

        ws.onclose = function () {
          let systemInfo;
          systemInfo =`<li class="systeminfo"><span>`
              +"与服务器连接断开，请刷新页面重试" +`</span></li>`;
          let myDate = new Date();
          let time = myDate.toLocaleDateString() + " " + myDate.toLocaleTimeString()
          console.log("serve 连接已关闭... " + time);
        };

        ws.onerror = function (evt) {
          ws.close()
          console.log("触发 onerror", evt)
        }
    }

    // Get the sender username for a message
    function getUsername(message) {
      if (message.user_id === parseInt(userInfo.value.uid)) {
        return 'You';
      }
      return message.username;
    }

    // Get the class for a message element (based on whether it is from the current user or not)
    function getMessageClass(message) {
      if (message.user_id === parseInt(userInfo.value.uid)) {
        return 'message-sent';
      }
      return 'message-received';
    }

    function getUsernameClass(message){
      if (message.user_id === parseInt(userInfo.value.uid)) {
        return 'username-sent';
      }
      return 'username-received';
    }

    return {
      userInfo,
      msgList,
      msgListCount,
      msgContainer,
      sendButton,
      sendContent,
      handleEnter,
      getUsername,
      getMessageClass,
      getUsernameClass,
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
  top: 5%;
  left: 50%;
  width: 40%;
  height: calc(85% - 50px);
  border: 1px solid #ccc;
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


.message-username {
  font-weight: bold;
  margin-right: 10px;
}

.username-sent{
  font-weight: bold;
  text-align: right;
}
.username-received{
  font-weight: bold;
  text-align: left;
}

.message-time {
  color: #999;
}

.textarea-wrapper {
  position: absolute;
  bottom: 5%;
  left: 50%;
  width: 40%;
  transform: translate(-50%, 0);
  display: flex;
  align-items: center;
  padding: 10px;
  border-top: 1px solid #ccc;
}

.chat-input {
  flex: 1;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

.send-button {
  width: 50px;
  height: 50px;
  margin-left: 10px;
  /*background-image: url('send-button-icon.png');*/
  background-size: cover;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}


</style>