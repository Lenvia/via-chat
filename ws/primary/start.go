package primary

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"via-chat/ws"
	"via-chat/ws/go_ws"
)

// 定义 serve 的映射关系
var serveMap = map[string]ws.ServeInterface{
	// Serve 和 GoServe，它们是实现了相同接口的不同类，能够根据用户的配置文件，选择合适的 ws.ServeInterface 实例
	"Serve":   &ws.Serve{},
	"GoServe": &go_ws.GoServe{},
}

func Create() ws.ServeInterface {
	// GoServe or Serve
	_type := viper.GetString("app.serve_type") // GoServe
	return serveMap[_type]                     // &go_ws.GoServe{}
}

// Start 启动 websocket
func Start(gin *gin.Context) {
	// 根据配置文件中，`app.serve_type` 键中对应的值创建 serve 实例，并启动服务
	Create().RunWs(gin)
}

// OnlineUserCount 返回在线用户的数量
func OnlineUserCount() int {
	return Create().GetOnlineUserCount()
}

// OnlineRoomUserCount 返回指定房间在线用户的数量
func OnlineRoomUserCount(roomId int) int {
	return Create().GetOnlineRoomUserCount(roomId)
}
