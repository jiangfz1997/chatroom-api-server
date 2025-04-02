package handlers

import (
	"chatroom-api/dynamodb"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type JoinChatroomRequest struct {
	Username   string `json:"username"`    // 用户名
	ChatroomID string `json:"chatroom_id"` // 聊天室 ID
}

// 退出聊天室请求体结构
type ExitChatroomRequest struct {
	Username   string `json:"username"`    // 用户名
	ChatroomID string `json:"chatroom_id"` // 聊天室 ID
}

func CreateChatroom(c *gin.Context) {
	log.Println("🔥 CreateChatroom 被触发")
	var req dynamodb.Chatroom
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	// 检查用户是否存在（你也可以放后面用户加入时验证）
	_, err := dynamodb.GetUserByUsername(req.CreatedBy)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 用 UUID 作为 room_id
	//roomID := uuid.New().String()
	roomID := generateRoomID()
	chatroom := dynamodb.Chatroom{
		RoomID:    roomID,
		Name:      req.Name,
		IsPrivate: req.IsPrivate,
		CreatedBy: req.CreatedBy,
		CreatedAt: time.Now().Format(time.RFC3339),
		Users:     []string{req.CreatedBy}, // 把创建者直接加入聊天室
	}

	if err := dynamodb.CreateChatroom(chatroom); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建聊天室失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "聊天室创建成功",
		"room_id":   roomID,
		"name":      chatroom.Name,
		"isPrivate": chatroom.IsPrivate,
	})
}

func JoinChatroom(c *gin.Context) {
	var req JoinChatroomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	// 验证用户是否存在
	_, err := dynamodb.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 检查聊天室是否存在
	_, err = dynamodb.GetChatroom(req.ChatroomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "聊天室不存在"})
		return
	}

	// 加入聊天室
	err = dynamodb.AddUserToChatroom(req.Username, req.ChatroomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "加入聊天室失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "加入聊天室成功"})
}

func ExitChatroom(c *gin.Context) {
	var req ExitChatroomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数格式错误"})
		return
	}

	// 用户是否存在
	_, err := dynamodb.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 移除用户
	err = dynamodb.RemoveUserFromChatroom(req.Username, req.ChatroomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "退出聊天室失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "退出聊天室成功"})
}
func GetUserChatrooms(c *gin.Context) {
	username := c.Param("username")

	// 检查用户是否存在
	_, err := dynamodb.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	chatrooms, err := dynamodb.GetChatroomsByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}

	// 构造返回
	var rooms []map[string]interface{}
	for _, room := range chatrooms {
		rooms = append(rooms, gin.H{
			"id":        room.RoomID,
			"name":      room.Name,
			"isPrivate": room.IsPrivate,
		})
	}

	c.JSON(http.StatusOK, gin.H{"rooms": rooms})
}
func GetChatroomMessages(c *gin.Context) {
	roomID := c.Param("roomId")
	before := c.Query("before")
	limitStr := c.DefaultQuery("limit", "20")
	username := c.Query("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 username 参数"})
		return
	}

	if before == "" {
		before = time.Now().Format(time.RFC3339)
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}

	messages, err := dynamodb.GetMessagesBefore(roomID, before, limit)
	if err != nil {
		fmt.Println("查询消息失败：", err)
		c.JSON(http.StatusOK, gin.H{"messages": []dynamodb.Message{}})
		return
	}

	if messages == nil {
		messages = []dynamodb.Message{}
	}

	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func EnterChatRoom(c *gin.Context) {
	roomID := c.Param("roomId")
	username := c.Query("username")

	if roomID == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roomId 和 username 是必须的"})
		return
	}

	// TODO: 后续可加负载均衡调度逻辑，这里先写死
	//wsHost := "ws://host.docker.internal:8081"
	// 构造返回的 WebSocket 地址
	wsHost := getNextWsHost()
	wsURL := fmt.Sprintf("%s/ws/%s?username=%s", wsHost, roomID, username)

	c.JSON(http.StatusOK, gin.H{
		"room_id": roomID,
		"ws_url":  wsURL,
	})
}

// For development only
var wsIndex = 0
var ports = []int{8081, 8081} // TODO: need to get from env

func getNextWsHost() string {
	port := ports[wsIndex%len(ports)]
	wsIndex++
	return fmt.Sprintf("ws://10.0.0.23:%d", port)
}

func GetChatroomByRoomID(c *gin.Context) {
	roomID := c.Param("roomId")

	chatroom, err := dynamodb.GetChatroom(roomID)
	if err != nil {
		log.Println("❌ 查询聊天室失败:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "聊天室不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        chatroom.RoomID,
		"name":      chatroom.Name,
		"isPrivate": chatroom.IsPrivate,
	})
}
