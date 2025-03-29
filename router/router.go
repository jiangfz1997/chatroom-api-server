package router

import (
	"chatroom-api/handlers"       // handler 层
	ws "chatroom-api/websocket"   // 自定义 websocket 包，取别名 ws
	"github.com/gin-contrib/cors" // 跨域中间件
	"github.com/gin-gonic/gin"
	"time"
)

// SetupRouter 初始化路由表，接收 hub 引用
func SetupRouter(hub *ws.Hub) *gin.Engine {
	r := gin.Default()

	// 启用 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 注册 API 路由
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// 注册 WebSocket 路由：客户端连接 ws://localhost:8080/ws/1?username=aaa
	r.GET("/ws/:roomId", handlers.ServeWs(hub))

	r.POST("/chatrooms", handlers.CreateChatroom)

	r.POST("/chatrooms/join", handlers.JoinChatroom)

	r.POST("/chatrooms/exit", handlers.ExitChatroom)

	r.GET("/chatrooms/user/:username", handlers.GetUserChatrooms)

	r.GET("/chatrooms/:roomId", handlers.GetChatroomByRoomID)

	r.GET("/messages/:roomId", handlers.GetChatroomMessages)

	return r
}
