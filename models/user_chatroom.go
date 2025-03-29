package models

import (
	"chatroom-api/database"
	"time"
)

// 结构体映射 user_chatroom 表
type UserChatroom struct {
	ID         int       // 主键 ID（自动增长）
	UserID     int       // 用户在 users 表中的 ID
	ChatroomID string    // 聊天室 UUID
	JoinedAt   time.Time // 加入时间
}

// 插入一条用户加入聊天室的记录
func AddUserToChatroom(userID int, chatroomID string) error {
	_, err := database.DB.Exec(
		"INSERT INTO user_chatroom (user_id, chatroom_id) VALUES (?, ?)",
		userID, chatroomID,
	)
	return err
}

// 查询用户加入的聊天室 ID 列表
func GetJoinedChatrooms(userID int) ([]string, error) {
	rows, err := database.DB.Query(
		"SELECT chatroom_id FROM user_chatroom WHERE user_id = ?",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatrooms []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		chatrooms = append(chatrooms, id)
	}
	return chatrooms, nil
}
