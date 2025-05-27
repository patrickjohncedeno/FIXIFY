package controller

import (
	"fixify_backend/middleware"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

func SaveGCashInfo(c *fiber.Ctx) error {
	type GCashInput struct {
		Name   string `json:"name"`
		Number string `json:"number"`
	}

	var input GCashInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	userClaims := c.Locals("user")
	if userClaims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	claims, ok := userClaims.(*users.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	db := middleware.GetDB()

	gcashInfo := users.Gcash{
		UserID:      claims.UserId,
		GcashName:   input.Name,
		GcashNumber: input.Number,
	}

	// Upsert (update if existing, insert if not)
	if err := db.
		Where("user_id = ?", claims.UserId).
		Assign(gcashInfo).
		FirstOrCreate(&gcashInfo).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save GCash info",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "GCash info saved successfully",
	})
}
