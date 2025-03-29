package database

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

var DB *sql.DB

func InitDB() {
	dir, _ := os.Getwd()
	fmt.Println("当前工作目录：", dir)

	var err error

	// 连接数据库，没有文件会自动创建
	DB, err = sql.Open("sqlite", "chatroom.db")
	if err != nil {
		log.Fatal("数据库连接失败：", err)
	}

	// 创建 users 表（如果不存在）
	createTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`
	if _, err := DB.Exec(createTable); err != nil {
		log.Fatal("建表失败：", err)
	}
	//打印初始化成功的日志
	fmt.Println("数据库初始化成功！chatroom.db 已准备好")

	// 创建 chatrooms 表（如果不存在）
	createRoomTable := `
	CREATE TABLE IF NOT EXISTS chatrooms (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		room_id TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		is_private BOOLEAN DEFAULT false,
		created_by TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := DB.Exec(createRoomTable); err != nil {
		log.Fatal("建表失败（chatrooms）:", err)
	}

	// 创建 user_chatroom 表（用于记录用户加入的聊天室）
	userChatroomTable := `
	CREATE TABLE IF NOT EXISTS user_chatroom (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		chatroom_id TEXT NOT NULL,
		joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	if _, err := DB.Exec(userChatroomTable); err != nil {
		log.Fatal("user_chatroom 表创建失败：", err)
	}

	// 创建 messages 表（用于保存聊天室消息）
	messageTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		room_id TEXT NOT NULL,
		sender TEXT NOT NULL,
		text TEXT NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := DB.Exec(messageTable); err != nil {
		log.Fatal("messages 表创建失败：", err)
	}

}
