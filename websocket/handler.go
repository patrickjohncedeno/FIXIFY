package websocket

import (
	"encoding/json"
	"fixify_backend/controller/signuplogin"
	"fixify_backend/model/users"
	"log"
	"strings"
	"time"

	"fixify_backend/middleware"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

type IncomingMessage struct {
	To      uint   `json:"to"`
	Content string `json:"content"`
}

// WebSocketHandler - Upgraded WebSocket handler
func WebSocketHandler(c *fiber.Ctx) error {
	// Try to get token from header
	token := c.Get("Authorization")

	if token == "" {
		// Try to get it from query string
		token = c.Query("token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).SendString("Missing token")
		}
	} else {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	claims := &users.Claims{}
	_, err := signuplogin.ParseJWTClaims(token, claims)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}

	c.Locals("admin_id", claims.UserId)
	return c.Next()
}

func WebSocketUpgrade(c *websocket.Conn) {
	db := middleware.DBConn

	// Extract admin_id from context
	adminID := c.Locals("admin_id").(uint)

	// Create a new Client instance for this WebSocket connection
	client := &Client{
		AdminID: adminID,
		Conn:    c,
		Send:    make(chan []byte),
	}

	log.Printf("WebSocket connection established for admin_id: %d", adminID)

	// Register the client with the Hub
	HubInstance.register <- client
	log.Printf("Client registered with admin_id: %d", adminID)

	// Read messages from the WebSocket connection
	go func() {
		log.Printf("Starting to read messages for admin_id: %d", adminID)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("Error reading message from admin_id %d: %v", adminID, err)
				break
			}

			log.Printf("Message from admin_id %d: %s", adminID, msg)

			var incomingMsg IncomingMessage
			if err := json.Unmarshal(msg, &incomingMsg); err != nil {
				log.Printf("Failed to parse message: %v", err)
				return
			}

			toAdminID := incomingMsg.To

			// Check if the recipient admin exists in the database
			var recipient users.Admin
			err = db.First(&recipient, toAdminID).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					log.Printf("Receiver admin with ID %d not found", toAdminID)
					return
				}
				log.Printf("Error checking recipient existence: %v", err)
				return
			}

			// Check for an existing conversation, or create a new one
			var conversation users.Conversation
			err = db.Where("(admin1_id = ? AND admin2_id = ?) OR (admin1_id = ? AND admin2_id = ?)", adminID, toAdminID, toAdminID, adminID).First(&conversation).Error

			if err != nil {
				if err == gorm.ErrRecordNotFound {
					conversation = users.Conversation{
						Admin1ID: adminID,
						Admin2ID: toAdminID,
					}

					if err := db.Create(&conversation).Error; err != nil {
						log.Printf("Failed to create conversation: %v", err)
						return
					}

					log.Printf("Created new conversation with ID: %d", conversation.ID)
				} else {
					log.Printf("Error checking for existing conversation: %v", err)
					return
				}
			} else {
				log.Printf("Found existing conversation with ID: %d", conversation.ID)
			}

			// Create and save the message in the database
			newMessage := users.Message{
				ConversationID: conversation.ID,
				SenderID:       adminID,
				Message:        incomingMsg.Content,
				IsRead:         false,
			}

			if err := db.Create(&newMessage).Error; err != nil {
				log.Printf("Failed to insert message into DB: %v", err)
			}

			if err := db.Model(&conversation).Update("updated_at", time.Now()).Error; err != nil {
				log.Printf("Failed to update conversation timestamp: %v", err)
			}

			// Broadcast the message to the recipient
			HubInstance.broadcast <- Message{
				FromAdminID: adminID,
				ToAdminID:   toAdminID,
				Content:     string(msg),
			}
		}
	}()

	// Send messages to the WebSocket connection
	go func() {
		for msg := range client.Send {
			if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("Error sending message to admin_id %d: %v", adminID, err)
				break
			}
		}
	}()

	// Defer closing the WebSocket connection
	defer func() {
		log.Printf("Closing WebSocket connection for admin_id: %d", adminID)
		HubInstance.unregister <- client
		c.Close()
		log.Printf("Unregistered client with admin_id: %d", adminID)
	}()

	// Keep the main goroutine alive
	select {}
}
