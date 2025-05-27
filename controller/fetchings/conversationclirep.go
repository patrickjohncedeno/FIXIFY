package fetchings

import (
	"fixify_backend/middleware"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FetchClientRepairmanMessages(c *fiber.Ctx) error {
	db := middleware.DBConn
	conversationID := c.Query("conversation_id")

	if conversationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "conversation_id is required",
			Data:    fiber.Map{"success": false},
		})
	}

	var messages []users.ClientRepairmanMessage

	if err := db.Preload("Conversation.Client", func(db *gorm.DB) *gorm.DB {
		return db.Select("user_id, first_name, last_name")
	}).
		Preload("Conversation.Repairman", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id, first_name, last_name")
		}).
		Preload("Sender", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id, first_name, last_name")
		}).
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch messages",
			Data: fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data: fiber.Map{
			"type":          "client_repairman_messages",
			"message_count": len(messages),
			"messages":      messages,
		},
	})
}
func FetchClientRepairmanConversations(c *fiber.Ctx) error {
	db := middleware.DBConn
	clientID := c.Query("client_id")
	repairmanID := c.Query("repairman_id")

	if clientID == "" && repairmanID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "client_id or repairman_id is required",
			Data:    fiber.Map{"success": false},
		})
	}

	var conversations []users.ClientRepairmanConversation

	query := db.Preload("Client", func(db *gorm.DB) *gorm.DB {
		return db.Select("user_id, first_name, last_name")
	}).
		Preload("Repairman", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id, first_name, last_name")
		}).
		Order("created_at DESC")

	if clientID != "" {
		query = query.Where("client_id = ?", clientID)
	} else {
		query = query.Where("repairman_id = ?", repairmanID)
	}

	if err := query.Find(&conversations).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch conversations",
			Data: fiber.Map{
				"success": false,
				"error":   err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data: fiber.Map{
			"type":               "client_repairman_conversations",
			"conversation_count": len(conversations),
			"conversations":      conversations,
		},
	})
}
