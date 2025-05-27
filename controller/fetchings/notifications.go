package fetchings

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

func FetchAllUserNotifications(c *fiber.Ctx) error {
	db := middleware.DBConn
	var notifications []users.UserNotification

	// Raw SQL query to fetch all users
	err := db.Raw("SELECT * FROM user_notifications").Scan(&notifications).Error
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500", // Internal server error
			Message: "Request failed",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch data from database",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Return the fetched data
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    notifications, // Return the raw data
	})
}

func ParamsNotification(c *fiber.Ctx) error {
	db := middleware.DBConn

	claims, ok := c.Locals("user").(*users.Claims)
	if !ok {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Unauthorized: Invalid or missing token",
			Data: errors.ErrorModel{
				Message:   "Unauthorized access",
				IsSuccess: false,
				Error:     "Invalid user context",
			},
		})
	}
	userID := claims.UserId

	var notifications []users.UserNotification

	err := db.Where("to_user = ?", userID).Find(&notifications).Error
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Request failed",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch data from database",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    notifications,
	})
}
