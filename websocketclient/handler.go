package websocketclient

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"fixify_backend/controller"
	"fixify_backend/controller/signuplogin"
	"fixify_backend/middleware"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
)

type IncomingMessage struct {
	To      uint   `json:"to"`
	Content string `json:"content"`
}

// WebSocketHandler - Upgraded WebSocket handler for Repairman and Client communication
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

	c.Locals("user_id", claims.UserId)
	return c.Next()
}
func WebSocketUpgrade(c *websocket.Conn) {
	db := middleware.DBConn

	// Extract user_id from context (repairman or client)
	userID := c.Locals("user_id").(uint)

	// Create a new Client instance for this WebSocket connection
	client := &Client{
		UserID: userID,
		Conn:   c,
		Send:   make(chan []byte),
	}

	log.Printf("WebSocket connection established for user_id: %d", userID)

	// Register the client with the Hub
	HubInstance.register <- client
	log.Printf("Client registered with user_id: %d", userID)

	// Read messages from the WebSocket connection
	go func() {
		log.Printf("Starting to read messages for user_id: %d", userID)
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("Error reading message from user_id %d: %v", userID, err)
				break
			}

			log.Printf("Message from user_id %d: %s", userID, msg)

			var incomingMsg IncomingMessage
			if err := json.Unmarshal(msg, &incomingMsg); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			toUserID := incomingMsg.To

			// Check if the recipient user exists
			var recipient users.User
			err = db.First(&recipient, toUserID).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					log.Printf("Receiver user with ID %d not found", toUserID)
					continue
				}
				log.Printf("Error checking recipient existence: %v", err)
				continue
			}

			// Check for or create conversation
			var conversation users.ClientRepairmanConversation
			err = db.Where("(client_id = ? AND repairman_id = ?) OR (client_id = ? AND repairman_id = ?)",
				userID, toUserID, toUserID, userID).First(&conversation).Error

			if err != nil {
				if err == gorm.ErrRecordNotFound {
					conversation = users.ClientRepairmanConversation{
						ClientId:    userID,
						RepairmanId: toUserID,
					}

					if err := db.Create(&conversation).Error; err != nil {
						log.Printf("Failed to create conversation: %v", err)
						continue
					}
					log.Printf("Created new conversation with ID: %d", conversation.ConversationId)
				} else {
					log.Printf("Error checking for existing conversation: %v", err)
					continue
				}
			} else {
				log.Printf("Found existing conversation with ID: %d", conversation.ConversationId)
			}

			// Save the message
			newMessage := users.ClientRepairmanMessage{
				ConversationId: conversation.ConversationId,
				SenderId:       userID,
				Message:        incomingMsg.Content,
				CreatedAt:      time.Now(),
			}

			if err := db.Create(&newMessage).Error; err != nil {
				log.Printf("Failed to insert message into DB. Error: %v, Message: %+v", err, newMessage)
				continue
			}
			log.Printf("Successfully saved message ID %d", newMessage.MessageId)

			// Trigger a user notification
			err = controller.CreateChatNotification(
				db,
				"new_message",
				int(userID),
				int(toUserID),
				incomingMsg.Content,
			)
			if err != nil {
				log.Printf("Failed to create user notification: %v", err)
			}

			// Update conversation timestamp
			if err := db.Model(&conversation).Update("updated_at", time.Now()).Error; err != nil {
				log.Printf("Failed to update conversation timestamp: %v", err)
			}

			// Broadcast message
			messageData := map[string]interface{}{
				"message_id":      newMessage.MessageId,
				"conversation_id": conversation.ConversationId,
				"sender_id":       userID,
				"content":         incomingMsg.Content,
				"created_at":      newMessage.CreatedAt.Format(time.RFC3339),
				"type":            "new_message",
				"is_duplicate":    false,
			}

			jsonData, err := json.Marshal(messageData)
			if err != nil {
				log.Printf("Error marshaling message: %v", err)
				continue
			}

			HubInstance.broadcast <- Message{
				FromUserID: userID,
				ToUserID:   toUserID,
				Content:    string(jsonData),
			}

			// Send push notification
			log.Printf("Attempting FCM to user %d (online: %v)", toUserID, HubInstance.IsOnline(toUserID))
			err = SendPushNotification(
				toUserID,
				"New Message",
				incomingMsg.Content,
				map[string]string{
					"type":            "new_message",
					"conversation_id": strconv.FormatUint(uint64(conversation.ConversationId), 10),
					"sender_id":       strconv.FormatUint(uint64(userID), 10),
					"message_content": incomingMsg.Content,
				},
			)
			if err != nil {
				log.Printf("FCM failed to %d: %v", toUserID, err)
			} else {
				log.Printf("FCM sent to %d", toUserID)
			}
		}
	}()

	// Send messages to the WebSocket connection
	go func() {
		for msg := range client.Send {
			if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("Error sending message to user_id %d: %v", userID, err)
				break
			}
		}
	}()

	// Defer closing the WebSocket connection
	defer func() {
		log.Printf("Closing WebSocket connection for user_id: %d", userID)
		HubInstance.unregister <- client
		c.Close()
		log.Printf("Unregistered client with user_id: %d", userID)
	}()

	select {} // Keep goroutine alive
}
