package fetchings

import (
	"fixify_backend/middleware"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

func Conversations(c *fiber.Ctx) error {
	db := middleware.DBConn

	var conversations []users.ClientRepairmanConversation

	err := db.Preload("Client").Preload("Repairman").
		Order("created_at DESC").
		Find(&conversations).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch client-repairman conversations",
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
