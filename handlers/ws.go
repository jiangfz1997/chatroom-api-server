package handlers

import (
	ws "chatroom-api/websocket" // 自定义 WebSocket 模块（hub.go 中定义），取别名为 ws
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// WebSocket 升级器：将 HTTP 请求升级为 WebSocket 连接
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许跨域（开发用，生产环境应限制）
	},
}

// ServeWs 是 WebSocket 的入口 handler，供路由 /ws/:roomId 使用
func ServeWs(hub *ws.Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取聊天室 ID 和用户名
		roomID := c.Param("roomId")
		username := c.Query("username") // 前端通过 ws://xxx/ws/1?username=xxx 传递用户名

		// 将 HTTP 升级为 WebSocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "WebSocket 升级失败"})
			return
		}

		// 创建客户端连接结构体
		client := &ws.Client{
			Conn:     conn,
			Username: username,
			RoomID:   roomID,
			Send:     make(chan []byte, 256),
		}

		// 将客户端注册到 Hub 中指定聊天室
		hub.JoinRoom(roomID, client)
		// 加入打印日志
		log.Printf("用户 %s 加入了房间 %s", username, roomID)
		// 添加在注册后立即发送欢迎语
		welcome := map[string]string{
			"sender": "Admin",
			"text":   "Welcome!",
		}
		msg, _ := json.Marshal(welcome)
		client.Conn.WriteMessage(websocket.TextMessage, msg)

		// 启动读写协程
		go client.ReadPump(hub)
		go client.WritePump()
	}
}
