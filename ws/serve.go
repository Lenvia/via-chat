package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"via-chat/models"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Serve struct {
	ServeInterface
}

func (serve *Serve) RunWs(gin *gin.Context) {
	Run(gin)
}

func (serve *Serve) GetOnlineUserCount() int {
	return GetOnlineUserCount()
}

func (serve *Serve) GetOnlineRoomUserCount(roomId int) int {
	return GetOnlineRoomUserCount(roomId)
}

// 客户端连接详情
type wsClients struct {
	Conn       *websocket.Conn `json:"conn"`
	RemoteAddr string          `json:"remote_addr"`
	Uid        float64         `json:"uid"`
	Username   string          `json:"username"`
	RoomId     string          `json:"room_id"`
	AvatarId   string          `json:"avatar_id"` // AvatarId 头像
	ToUser     interface{}     `json:"to_user"`
}

// client & serve 的消息体
type msg struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

// 变量定义初始化
var (
	wsUpgrader  = websocket.Upgrader{} // 用于将一个普通的 http 请求升级成 websocket 连接。实现了这个升级的过程，可以让客户端和服务端之间建立一个 websocket 的长连接，用于双向数据传输。
	clientMsg   = msg{}
	mutex       = sync.Mutex{}
	rooms       = [roomCount + 1][]wsClients{} // rooms 是一个包含多个 []wsClients 类型的数组，用于存储各个房间在线的用户，数组中的每个元素代表一个房间，每个元素中存储了当前房间在线的所有用户。
	privateChat = []wsClients{}                // privateChat 是一个 []wsClients 类型的切片，用于存储私聊的消息记录
)

// 定义消息类型
const msgTypeOnline = 1        // 上线
const msgTypeOffline = 2       // 离线
const msgTypeSend = 3          // 消息发送
const msgTypeGetOnlineUser = 4 // 获取用户列表
const msgTypePrivateChat = 5   // 私聊

const roomCount = 6 // 房间总数

func Run(gin *gin.Context) {

	// @see https://github.com/gorilla/websocket/issues/523
	// wsUpgrader.CheckOrigin 是用来解决 websocket 跨域问题的，这里设置为返回 true，表示接收来自任何源的请求。
	wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }

	c, _ := wsUpgrader.Upgrade(gin.Writer, gin.Request, nil)

	defer c.Close() // 长连接，在函数退出时需要 defer 关闭。

	mainProcess(c) // 通过 mainProcess() 函数处理接收到的 websocket 消息
}

// mainProcess 主程序，循环读取客户端发送的消息，和发送消息到客户端
func mainProcess(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		serveMsgStr := message

		// 处理心跳响应 , heartbeat为与客户端约定的值
		if string(serveMsgStr) == `heartbeat` {
			c.WriteMessage(websocket.TextMessage, []byte(`{"status":0,"data":"heartbeat ok"}`))
			continue
		}
		// 如果是上线消息，则处理连接请求
		json.Unmarshal(message, &clientMsg)
		// log.Println("来自客户端的消息", clientMsg,c.RemoteAddr())
		if clientMsg.Data == nil {
			return
			//mainProcess(c)
		}

		if err != nil { // 离线通知
			log.Println("ReadMessage error1", err)
			disconnect(c)
			c.Close()
			return
		}

		// 进入房间，建立连接。将该用户加入到房间的在线用户列表中
		if clientMsg.Status == msgTypeOnline {
			handleConnClients(c)
			serveMsgStr = formatServeMsgStr(msgTypeOnline) // 格式化传送给客户端的消息数据
		}

		// 处理私聊消息。先寻找接收者的连接，如果找到则将消息发送给接收者
		if clientMsg.Status == msgTypePrivateChat {
			// 处理私聊
			serveMsgStr = formatServeMsgStr(msgTypePrivateChat)
			toC := findToUserCoonClient()
			if toC != nil {
				toC.(wsClients).Conn.WriteMessage(websocket.TextMessage, serveMsgStr)
			}
		}

		// 处理消息发送。
		if clientMsg.Status == msgTypeSend { // 消息发送
			serveMsgStr = formatServeMsgStr(msgTypeSend)
		}

		// 处理获取在线用户列表的请求。发送在线用户列表给客户端
		if clientMsg.Status == msgTypeGetOnlineUser {
			serveMsgStr = formatServeMsgStr(msgTypeGetOnlineUser)
			c.WriteMessage(websocket.TextMessage, serveMsgStr)
			continue
		}

		//log.Println("serveMsgStr", string(serveMsgStr))

		// 如果收到的消息是发送消息或上线消息，则将消息发送给所有房间内的在线用户
		if clientMsg.Status == msgTypeSend || clientMsg.Status == msgTypeOnline {
			notify(c, string(serveMsgStr))
		}
	}
}

// findToUserCoonClient 获取私聊的用户连接
func findToUserCoonClient() interface{} {
	_, roomIdInt := getRoomId()

	toUserUid := clientMsg.Data.(map[string]interface{})["to_uid"].(string)

	for _, c := range rooms[roomIdInt] {
		stringUid := strconv.FormatFloat(c.Uid, 'f', -1, 64)
		if stringUid == toUserUid {
			return c
		}
	}

	return nil
}

// 处理建立连接的用户
func handleConnClients(c *websocket.Conn) {
	roomId, roomIdInt := getRoomId()

	for cKey, wcl := range rooms[roomIdInt] {
		if wcl.Uid == clientMsg.Data.(map[string]interface{})["uid"].(float64) {
			mutex.Lock()
			// 通知当前用户下线
			wcl.Conn.WriteMessage(websocket.TextMessage, []byte(`{"status":-1,"data":[]}`))
			rooms[roomIdInt] = append(rooms[roomIdInt][:cKey], rooms[roomIdInt][cKey+1:]...)
			wcl.Conn.Close()
			mutex.Unlock()
		}
	}

	mutex.Lock()
	rooms[roomIdInt] = append(rooms[roomIdInt], wsClients{
		Conn:       c,
		RemoteAddr: c.RemoteAddr().String(),
		Uid:        clientMsg.Data.(map[string]interface{})["uid"].(float64),
		Username:   clientMsg.Data.(map[string]interface{})["username"].(string),
		RoomId:     roomId,
		AvatarId:   clientMsg.Data.(map[string]interface{})["avatar_id"].(string),
	})
	mutex.Unlock()
}

// 统一消息发放
func notify(conn *websocket.Conn, msg string) {
	_, roomIdInt := getRoomId()
	for _, con := range rooms[roomIdInt] {
		if con.RemoteAddr != conn.RemoteAddr().String() {
			con.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}

// 离线通知
func disconnect(conn *websocket.Conn) {
	_, roomIdInt := getRoomId()
	for index, con := range rooms[roomIdInt] {
		if con.RemoteAddr == conn.RemoteAddr().String() {
			data := map[string]interface{}{
				"username": con.Username,
				"uid":      con.Uid,
				"time":     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()
			}

			jsonStrServeMsg := msg{
				Status: msgTypeOffline,
				Data:   data,
			}
			serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

			disMsg := string(serveMsgStr)

			mutex.Lock()
			rooms[roomIdInt] = append(rooms[roomIdInt][:index], rooms[roomIdInt][index+1:]...)
			con.Conn.Close()
			mutex.Unlock()
			notify(conn, disMsg)
		}
	}
}

// 格式化传送给客户端的消息数据
func formatServeMsgStr(status int) []byte {

	roomId, roomIdInt := getRoomId()

	data := map[string]interface{}{
		"username": clientMsg.Data.(map[string]interface{})["username"].(string),
		"uid":      clientMsg.Data.(map[string]interface{})["uid"].(float64),
		"room_id":  roomId,
		"time":     time.Now().UnixNano() / 1e6, // 13位  10位 => now.Unix()   获取当前Unix时间戳，并以13位长度的int64格式保存
	}

	// 当普通发送消息或者私聊消息时
	if status == msgTypeSend || status == msgTypePrivateChat {
		data["avatar_id"] = clientMsg.Data.(map[string]interface{})["avatar_id"].(string)
		data["content"] = clientMsg.Data.(map[string]interface{})["content"].(string)

		toUidStr := clientMsg.Data.(map[string]interface{})["to_uid"].(string)
		toUid, _ := strconv.Atoi(toUidStr)

		// 保存消息
		stringUid := strconv.FormatFloat(data["uid"].(float64), 'f', -1, 64)
		intUid, _ := strconv.Atoi(stringUid)

		if _, ok := clientMsg.Data.(map[string]interface{})["image_url"]; ok {
			// 存在图片，调用models包的SaveContent函数，保存聊天内容到数据库中
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"room_id":    data["room_id"],
				"content":    data["content"],
				"image_url":  clientMsg.Data.(map[string]interface{})["image_url"].(string),
			})
		} else {
			models.SaveContent(map[string]interface{}{
				"user_id":    intUid,
				"to_user_id": toUid,
				"room_id":    data["room_id"],
				"content":    data["content"],
			})
		}
	}

	if status == msgTypeGetOnlineUser {
		data["count"] = GetOnlineRoomUserCount(roomIdInt)
		data["list"] = onLineUserList(roomIdInt)
	}

	jsonStrServeMsg := msg{
		Status: status,
		Data:   data,
	}
	serveMsgStr, _ := json.Marshal(jsonStrServeMsg)

	return serveMsgStr
}

func getRoomId() (string, int) {
	roomId := clientMsg.Data.(map[string]interface{})["room_id"].(string) // 两次类型断言
	roomIdInt, _ := strconv.Atoi(roomId)
	return roomId, roomIdInt
}

// 获取在线用户列表
func onLineUserList(roomId int) []wsClients {
	return rooms[roomId]
}

// =======================对外方法=====================================

func GetOnlineUserCount() int {
	num := 0
	for i := 1; i <= roomCount; i++ {
		num = num + GetOnlineRoomUserCount(i)
	}
	return num
}

func GetOnlineRoomUserCount(roomId int) int {
	return len(rooms[roomId])
}
