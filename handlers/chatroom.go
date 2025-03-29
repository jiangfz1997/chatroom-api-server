package handlers

import (
	"chatroom-api/database"
	"chatroom-api/models"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// 生成唯一 room_id
func generateRoomID() string {
	bytes := make([]byte, 6)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// 通过 room_id 查询聊天室信息（用于加入前的校验和展示）
func GetChatroomByRoomID(c *gin.Context) {
	roomID := c.Param("roomId")

	var name string
	var isPrivate bool

	err := database.DB.QueryRow(`
		SELECT name, is_private 
		FROM chatrooms 
		WHERE room_id = ?`, roomID).Scan(&name, &isPrivate)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "聊天室不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        roomID,
		"name":      name,
		"isPrivate": isPrivate,
	})
}

// 创建聊天室的 handler
func CreateChatroom(c *gin.Context) {
	var req models.Chatroom
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	roomID := generateRoomID()

	_, err := database.DB.Exec(`
		INSERT INTO chatrooms (room_id, name, is_private, created_by) 
		VALUES (?, ?, ?, ?)`,
		roomID, req.Name, req.IsPrivate, req.CreatedBy,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建聊天室失败"})
		return
	}

	// 查找用户 ID（用于写入 user_chatroom 表）
	var userID int
	err = database.DB.QueryRow("SELECT id FROM users WHERE username = ?", req.CreatedBy).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 将创建者添加进 user_chatroom 表
	err = models.AddUserToChatroom(userID, roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "聊天室创建成功，但用户加入失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "聊天室创建成功",
		"room_id":   roomID,
		"name":      req.Name,
		"isPrivate": req.IsPrivate,
	})
}

// 加入聊天室的请求体结构
type JoinChatroomRequest struct {
	Username   string `json:"username"`    // 用户名
	ChatroomID string `json:"chatroom_id"` // 聊天室 ID
}

// 退出聊天室请求体结构
type ExitChatroomRequest struct {
	Username   string `json:"username"`    // 用户名
	ChatroomID string `json:"chatroom_id"` // 聊天室 ID
}

// 加入聊天室的 handler
func JoinChatroom(c *gin.Context) {
	var req JoinChatroomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	// 查询用户 ID（通过用户名）
	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = ?", req.Username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		}
		return
	}

	// 检查聊天室是否存在
	var exists int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM chatrooms WHERE room_id = ?", req.ChatroomID).Scan(&exists)
	if err != nil || exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "聊天室不存在"})
		return
	}

	// 插入用户-聊天室记录
	err = models.AddUserToChatroom(userID, req.ChatroomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加入聊天室失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "加入聊天室成功"})
}

// 退出聊天室的 handler
func ExitChatroom(c *gin.Context) {
	var req ExitChatroomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	// 查询用户 ID
	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = ?", req.Username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 删除 user_chatroom 中的记录
	_, err = database.DB.Exec(`
		DELETE FROM user_chatroom
		WHERE user_id = ? AND chatroom_id = ?
	`, userID, req.ChatroomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "退出聊天室失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "退出聊天室成功"})
}

// 查询用户加入了哪些chatroom，呈现在前端页面
func GetUserChatrooms(c *gin.Context) {
	username := c.Param("username")

	// 查询用户ID
	var userID int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 联合查询用户加入的聊天室信息
	rows, err := database.DB.Query(`
		SELECT c.room_id, c.name, c.is_private
		FROM user_chatroom uc
		JOIN chatrooms c ON uc.chatroom_id = c.room_id
		WHERE uc.user_id = ?`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	defer rows.Close()

	// 构造结果
	var rooms []map[string]interface{}
	for rows.Next() {
		var roomID, name string
		var isPrivate bool
		if err := rows.Scan(&roomID, &name, &isPrivate); err != nil {
			continue
		}
		rooms = append(rooms, gin.H{
			"id":        roomID,
			"name":      name,
			"isPrivate": isPrivate,
		})
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}

// 获取聊天室的历史消息（分页）
func GetChatroomMessages(c *gin.Context) {
	roomID := c.Param("roomId")
	before := c.Query("before")               // 时间戳字符串（可选）
	limitStr := c.DefaultQuery("limit", "20") // 限制条数（默认20）
	username := c.Query("username")
	fmt.Println("前端传入的 before 参数是：", before)

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 username 参数"})
		return
	}

	if before == "" {
		before = time.Now().Format("2006-01-02 15:04:05") // 此处有修改
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20 // fallback 默认值
	}

	messages, err := models.GetMessagesWithJoinLimit(roomID, username, before, limit)
	if err != nil {
		fmt.Println("查询消息失败：", err)
		// 这里返回空数组而不是500
		c.JSON(http.StatusOK, gin.H{"messages": []models.Message{}})
		return
	}

	// 即使 messages 是 nil，也要返回空数组，防止前端拿到 null
	if messages == nil {
		messages = []models.Message{}
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}
