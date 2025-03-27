package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
)

// 客户端结构
type Client struct {
	Conn     *websocket.Conn
	Username string
	RoomID   string
	Send     chan []byte
}

// 房间结构
type Room struct {
	ID      string
	Clients map[*Client]bool
	Lock    sync.Mutex
}

// 中心 Hub：管理所有房间
type Hub struct {
	Rooms map[string]*Room
	Lock  sync.Mutex
}

var GlobalHub = Hub{
	Rooms: make(map[string]*Room),
}

// 添加客户端到房间
func (h *Hub) JoinRoom(roomID string, client *Client) {
	h.Lock.Lock()
	room, exists := h.Rooms[roomID]
	if !exists {
		room = &Room{
			ID:      roomID,
			Clients: make(map[*Client]bool),
		}
		h.Rooms[roomID] = room
	}
	h.Lock.Unlock()

	room.Lock.Lock()
	room.Clients[client] = true
	room.Lock.Unlock()
}

// 从房间移除客户端
func (h *Hub) LeaveRoom(roomID string, client *Client) {
	h.Lock.Lock()
	room, exists := h.Rooms[roomID]
	h.Lock.Unlock()
	if !exists {
		return
	}
	room.Lock.Lock()
	delete(room.Clients, client)
	room.Lock.Unlock()
}

// 向房间广播消息
func (h *Hub) Broadcast(roomID string, message []byte) {
	h.Lock.Lock()
	room, exists := h.Rooms[roomID]
	h.Lock.Unlock()
	if !exists {
		return
	}
	room.Lock.Lock()
	defer room.Lock.Unlock()

	for client := range room.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(room.Clients, client)
		}
	}
}
