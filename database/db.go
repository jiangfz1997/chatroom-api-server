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
}
