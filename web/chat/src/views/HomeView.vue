<template>
  <div class="home">
    <div class="header">
      <div class="toptext">
        <span><b>选择一个聊天室</b></span>
      </div>
      <ul class="topnavlist">
        <li class="userinfo">
          <span><b>{{ user_info.username }}</b></span>
        </li>
      </ul>
    </div>
    <div class="main container">
      <div class="room_list">
        <el-row :gutter="20">
          <el-col v-for="room in rooms" :key="room.id" :xs="12" :sm="12" :md="8">
            <div class="room" @click="$router.push(`/room/${room.id}`)">
              <span class="num">{{ room.num }}</span>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>

<script>
import { ElRow, ElCol } from "element-plus";
import app from "@/main";

export default {
  name: "HomeView",
  components: {
    ElRow,
    ElCol,
  },
  data() {
    return {
      user_info:{
        username:null,
      },
      rooms: [],
    };
  },
  mounted() {
    this.getData();
  },
  methods: {
    async getData() {
      try {
        const res = await app.config.globalProperties.$http.get("/home");
        console.log(res)
        if (res.status === 200) {
          const { data, user_info } = res.data;
          this.rooms = data;
          this.user_info = user_info;

          console.log(this.user_info)
        }
      } catch (error) {
        console.error(error);
      }
    },
  },
};
</script>

<style>
.room {
  width: 150px;
  height: 150px;
  border: 1px solid #ddd;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 36px;
  color: #333;
  cursor: pointer;
}
.num {
  font-weight: bold;
}
</style>