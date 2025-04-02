package database

import (
	"database/sql"
	"fmt"
	"log"
	//"os"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	dsn := "host=localhost user=chat password=123456 dbname=chatroom port=5432 sslmode=disable"
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("❌ PostgreSQL 连接失败：", err)
	}

	// 检查连接
	if err := DB.Ping(); err != nil {
		log.Fatal("❌ 数据库无法连接：", err)
	}

	// 建表语句
	createUsers := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`
	_, _ = DB.Exec(createUsers)

	createRooms := `
	CREATE TABLE IF NOT EXISTS chatrooms (
		id SERIAL PRIMARY KEY,
		room_id TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		is_private BOOLEAN DEFAULT false,
		created_by TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, _ = DB.Exec(createRooms)

	createUserRoom := `
	CREATE TABLE IF NOT EXISTS user_chatroom (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		chatroom_id TEXT NOT NULL,
		joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`
	_, _ = DB.Exec(createUserRoom)

	createMessages := `
	CREATE TABLE IF NOT EXISTS messages (
		id SERIAL PRIMARY KEY,
		room_id TEXT NOT NULL,
		sender TEXT NOT NULL,
		text TEXT NOT NULL,
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, _ = DB.Exec(createMessages)

	fmt.Println("✅ PostgreSQL 数据库初始化完成！")
}
