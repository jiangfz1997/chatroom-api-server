package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

// 心跳超时设置
const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.LeaveRoom(c.RoomID, c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("读消息错误:", err)
			break
		}

		var incoming map[string]string
		if err := json.Unmarshal(message, &incoming); err != nil {
			log.Println("解析前端消息失败:", err)
			continue
		}

		msg := map[string]string{
			"sender": c.Username,
			"text":   incoming["text"],
		}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Println("JSON 编码失败:", err)
			continue
		}

		hub.Broadcast(c.RoomID, jsonMsg)
	}
}

// 将消息发送到客户端
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// channel 关闭
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := c.Conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("写消息错误:", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
