package main

import (
	"chatroom-api/dynamodb"
	//"chatroom-api/database"
	"chatroom-api/logger"
	"chatroom-api/router"
	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger() // 初始化日志系统
	log := logger.Log   // 使用自定义 logrus 实例
	//_ = godotenv.Load(".env")
	//dynamodb.InitDB() // 初始化 SQLite
	////dynamodb.CreateAllTables()
	////hub := &websocket.GlobalHub // 获取全局 Hub 实例
	//r := router.SetupRouter()
	//r.Run(":8080")
	log.Info("服务器启动流程开始")

	err := godotenv.Load(".env")
	if err != nil {
		log.Warn("未找到 .env 文件，将使用默认环境变量")
	} else {
		log.Info(".env 文件加载成功")
	}

	log.Info("开始初始化数据库")
	dynamodb.InitDB()
	log.Info("数据库初始化完成")

	r := router.SetupRouter()
	log.Info("启动 HTTP 服务监听 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
