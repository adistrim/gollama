package socket

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// allowing all origins for now
		return true
	},
}

type Connection struct {
	ws        *websocket.Conn
	send      chan []byte
	sessionID string
	request   *http.Request
}

type Message struct {
	Content   string `json:"content"`
	SessionID string `json:"session_id,omitempty"`
	Response  string `json:"response,omitempty"`
	Error     string `json:"error,omitempty"`
	IsProcessing  bool   `json:"is_processing,omitempty"`
}

func NewConnection(c *gin.Context) (*Connection, error) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}

	conn := &Connection{
		ws:      ws,
		send:    make(chan []byte, 256),
		request: c.Request,
	}

	return conn, nil
}

func (c *Connection) Context() context.Context {
	return c.request.Context()
}

func (c *Connection) ReadPump(handler func(*Connection, Message)) {
	defer c.ws.Close()

	for {
		var msg Message
		err := c.ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
		
		if msg.Content != "" {
			handler(c, msg)
		}
	}
}

func (c *Connection) WritePump() {
	defer c.ws.Close()

	for message := range c.send {
		c.ws.WriteMessage(websocket.TextMessage, message)
	}
}

func (c *Connection) SendMessage(msg Message) {
	data, _ := json.Marshal(msg)
	select {
	case c.send <- data:
	default:
		close(c.send)
	}
}
