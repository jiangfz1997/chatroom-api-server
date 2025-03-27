package main

import (
	"chatroom-api/database"
	"chatroom-api/router"
	"chatroom-api/websocket"
)

func main() {
	database.InitDB()           // 初始化 SQLite
	hub := &websocket.GlobalHub // 获取全局 Hub 实例
	r := router.SetupRouter(hub)
	r.Run(":8080")
}
