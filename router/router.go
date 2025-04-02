package router

import (
	"chatroom-api/handlers" // handler 层
	"chatroom-api/middleware"
	"github.com/gin-contrib/cors" // 跨域中间件
	"github.com/gin-gonic/gin"
	"time"
)

// SetupRouter 初始化路由表，接收 hub 引用
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 启用 CORS 中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 注册 API 路由
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// 注册 WebSocket 路由：客户端连接 ws://localhost:8080/ws/1?username=aaa
	//r.GET("/ws/:roomId", handlers.ServeWs(hub))

	//以下接口走鉴权
	//r.POST("/chatrooms", handlers.CreateChatroom)
	//
	//r.POST("/chatrooms/join", handlers.JoinChatroom)
	//
	//r.POST("/chatrooms/exit", handlers.ExitChatroom)
	//
	//r.GET("/chatrooms/user/:username", handlers.GetUserChatrooms)
	//
	//r.GET("/chatrooms/:roomId", handlers.GetChatroomByRoomID)
	//
	//r.GET("/messages/:roomId", handlers.GetChatroomMessages)
	//
	//r.GET("/chatrooms/:roomId/enter", handlers.EnterChatRoom)

	// 需要鉴权的接口挂在 auth group 下
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.POST("/chatrooms", handlers.CreateChatroom)
	auth.POST("/chatrooms/join", handlers.JoinChatroom)
	auth.POST("/chatrooms/exit", handlers.ExitChatroom)
	auth.GET("/chatrooms/user/:username", handlers.GetUserChatrooms)
	auth.GET("/chatrooms/:roomId", handlers.GetChatroomByRoomID)
	auth.GET("/messages/:roomId", handlers.GetChatroomMessages)
	auth.GET("/chatrooms/:roomId/enter", handlers.EnterChatRoom)

	return r
}
