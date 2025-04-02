package main

import (
	"chatroom-api/dynamodb"
	//"chatroom-api/database"
	"chatroom-api/router"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")
	dynamodb.InitDB() // 初始化 SQLite
	//dynamodb.CreateAllTables()
	//hub := &websocket.GlobalHub // 获取全局 Hub 实例
	r := router.SetupRouter()
	r.Run(":8080")
}
