package websocketclient

import (
	"fixify_backend/middleware"
	"fixify_backend/model/users"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func EnsureClientRepairmanConversation(clientID uint, repairmanID uint) (uint, error) {
	db := middleware.DBConn

	if db == nil {
		log.Println("Database connection is nil")
		return 0, fmt.Errorf("database connection not established")
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic in EnsureClientRepairmanConversation: %v", r)
		}
	}()

	var conversationID uint
	var existingConvo users.ClientRepairmanConversation
	err := tx.Where("client_id = ? AND repairman_id = ?", clientID, repairmanID).
		First(&existingConvo).Error

	if err == nil {
		// Conversation exists - always create welcome message
		conversationID = existingConvo.ConversationId
		log.Printf("Conversation found - ID: %d, Client: %d, Repairman: %d",
			conversationID, clientID, repairmanID)

		// Create welcome message regardless of existing messages
		welcomeMessage := users.ClientRepairmanMessage{
			ConversationId: conversationID,
			SenderId:       repairmanID,
			Message:        "Hello! I've accepted your service request. Let's discuss the details.",
			CreatedAt:      time.Now(),
		}

		if err := tx.Create(&welcomeMessage).Error; err != nil {
			tx.Rollback()
			log.Printf("Failed to create welcome message: %v", err)
			return 0, fmt.Errorf("failed to create welcome message: %v", err)
		}

		// Send notification to client
		go func() {
			notificationTitle := "New message from repairman"
			notificationBody := welcomeMessage.Message
			data := map[string]string{
				"conversation_id": fmt.Sprintf("%d", conversationID),
				"type":            "new_message",
				"sender_id":       fmt.Sprintf("%d", repairmanID),
			}

			if err := SendPushNotification(
				clientID,
				notificationTitle,
				notificationBody,
				data,
			); err != nil {
				log.Printf("Failed to send notification: %v", err)
			}
		}()

		log.Printf("Created welcome message for conversation %d", conversationID)
	} else if err == gorm.ErrRecordNotFound {
		// Create new conversation
		newConvo := users.ClientRepairmanConversation{
			ClientId:    clientID,
			RepairmanId: repairmanID,
			CreatedAt:   time.Now(),
		}

		if err := tx.Create(&newConvo).Error; err != nil {
			tx.Rollback()
			log.Printf("Failed to create conversation: %v", err)
			return 0, fmt.Errorf("failed to create conversation: %v", err)
		}

		conversationID = newConvo.ConversationId
		log.Printf("Created new conversation - ID: %d, Client: %d, Repairman: %d",
			conversationID, clientID, repairmanID)

		// Create welcome message
		welcomeMessage := users.ClientRepairmanMessage{
			ConversationId: conversationID,
			SenderId:       repairmanID,
			Message:        "Hello! I've accepted your service request. Let's discuss the details.",
			CreatedAt:      time.Now(),
		}

		if err := tx.Create(&welcomeMessage).Error; err != nil {
			tx.Rollback()
			log.Printf("Failed to create welcome message: %v", err)
			return 0, fmt.Errorf("failed to create welcome message: %v", err)
		}

		// Send notification to client
		go func() {
			notificationTitle := "New message from repairman"
			notificationBody := welcomeMessage.Message
			data := map[string]string{
				"conversation_id": fmt.Sprintf("%d", conversationID),
				"type":            "new_message",
				"sender_id":       fmt.Sprintf("%d", repairmanID),
			}

			if err := SendPushNotification(
				clientID,
				notificationTitle,
				notificationBody,
				data,
			); err != nil {
				log.Printf("Failed to send notification: %v", err)
			}
		}()

		log.Printf("Created welcome message for new conversation %d", conversationID)
	} else {
		tx.Rollback()
		log.Printf("Error checking conversation: %v", err)
		return 0, fmt.Errorf("error checking conversation: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return 0, fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully processed conversation for client %d and repairman %d", clientID, repairmanID)
	return conversationID, nil
}
