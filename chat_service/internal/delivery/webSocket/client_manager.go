package websocket

import (
	"chat_service/internal/domain"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientManager struct {
	clients map[*websocket.Conn]bool
	lock    sync.Mutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (cm *ClientManager) AddClient(conn *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	cm.clients[conn] = true
}

func (cm *ClientManager) RemoveClient(conn *websocket.Conn) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	delete(cm.clients, conn)
}

func (cm *ClientManager) Broadcast(message []byte) {
	cm.lock.Lock()
	defer cm.lock.Unlock()

	for conn := range cm.clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			conn.Close()
			delete(cm.clients, conn)
		}
	}
}

func (cm *ClientManager) GetBroadcaster() domain.Broadcaster {
	return cm
}
