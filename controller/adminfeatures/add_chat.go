package adminfeatures

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)
func FetchAvailableAdminsForConversation(c *fiber.Ctx) error {
	db := middleware.DBConn
	var availableAdmins []users.Admin

	// Get current admin's ID (assuming it's stored in context)
	user := c.Locals("user").(*users.Claims)
	currentAdminID := user.UserId

	err := db.Raw(`
		SELECT * FROM admins
		WHERE admin_id != ?
		AND admin_id NOT IN (
			SELECT admin1_id FROM conversations WHERE admin2_id = ?
			UNION
			SELECT admin2_id FROM conversations WHERE admin1_id = ?
		)
	`, currentAdminID, currentAdminID, currentAdminID).Scan(&availableAdmins).Error

	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch available admins",
			Data: errors.ErrorModel{
				Message:   "Database query error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    availableAdmins,
	})
}
