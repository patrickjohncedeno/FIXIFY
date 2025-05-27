package websocket

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	AdminID uint
	Conn    *websocket.Conn
	Send    chan []byte
}

type Hub struct {
	clients    map[uint]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	lock       sync.RWMutex
}

type Message struct {
	FromAdminID uint   `json:"from"`
	ToAdminID   uint   `json:"to"`
	Content     string `json:"content"`
}

var HubInstance = &Hub{
	clients:    make(map[uint]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	broadcast:  make(chan Message),
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.lock.Lock()
			h.clients[client.AdminID] = client
			h.lock.Unlock()

		case client := <-h.unregister:
			h.lock.Lock()
			if _, ok := h.clients[client.AdminID]; ok {
				delete(h.clients, client.AdminID)
				close(client.Send)
			}
			h.lock.Unlock()

		case msg := <-h.broadcast:
			h.lock.RLock()
			if receiver, ok := h.clients[msg.ToAdminID]; ok {
				data := []byte(msg.Content)
				receiver.Send <- data
			}
			h.lock.RUnlock()
		}
	}
}
func (Message) TableName() string {
	return "messages" // Explicitly define the table name here
}
