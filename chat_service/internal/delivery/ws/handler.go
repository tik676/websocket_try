package ws

import (
	"chat_service/internal/domain"
	"chat_service/internal/usecase"
	"encoding/json"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WsHandler struct {
	uc        *usecase.UseCase
	upgrader  *websocket.Upgrader
	clients   map[*websocket.Conn]bool
	broadcast chan *domain.Message
	mu        sync.Mutex
}

func NewWsHandler(uc *usecase.UseCase, upgrader *websocket.Upgrader) *WsHandler {
	handler := &WsHandler{
		uc:        uc,
		upgrader:  upgrader,
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan *domain.Message),
	}
	go handler.run()

	return handler
}

func (ws *WsHandler) run() {
	for msg := range ws.broadcast {
		ws.mu.Lock()
		for client := range ws.clients {
			if err := client.WriteJSON(msg); err != nil {
				client.Close()
				delete(ws.clients, client)
			}
		}
		ws.mu.Unlock()
	}
}

func (ws *WsHandler) Broadcast(msg []byte) {
	var m domain.Message
	if err := json.Unmarshal(msg, &m); err == nil {
		ws.broadcast <- &m
	}
}

func (ws *WsHandler) HandleWS(c *gin.Context) {
	conn, err := ws.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed connect to websocket"})
		return
	}

	userID := c.GetInt64("user_id")
	username := c.GetString("name")
	role := c.GetString("role")

	ws.mu.Lock()
	ws.clients[conn] = true
	ws.mu.Unlock()

	go func() {
		defer func() {
			ws.mu.Lock()
			delete(ws.clients, conn)
			ws.mu.Unlock()
			conn.Close()
		}()
		for {

			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("read error: %v", err)
				return
			}
			message := domain.Message{
				UserID:   userID,
				Username: username,
				Content:  string(msg),
				Role:     role,
			}

			if saveMsg, err := ws.uc.SendMessage(message); err == nil {
				ws.broadcast <- &saveMsg

			}
		}
	}()
}
