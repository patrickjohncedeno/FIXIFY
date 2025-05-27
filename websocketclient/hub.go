package websocketclient

import (
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	UserID uint
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	clients    map[uint]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan Message
	lock       sync.RWMutex
}

type Message struct {
	FromUserID uint   `json:"from"`
	ToUserID   uint   `json:"to"`
	Content    string `json:"content"`
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
			h.clients[client.UserID] = client
			h.lock.Unlock()

		case client := <-h.unregister:
			h.lock.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.lock.Unlock()

		case msg := <-h.broadcast:
			h.lock.RLock()
			if receiver, ok := h.clients[msg.ToUserID]; ok {
				data := []byte(msg.Content)
				receiver.Send <- data
			}
			h.lock.RUnlock()
		}
	}
}

func (Message) TableName() string {
	return "client_repairman_messages" // Explicitly define the table name here
}

func (h *Hub) IsOnline(userID uint) bool {
	_, ok := h.clients[userID]
	return ok
}
