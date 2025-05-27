package adminfeatures

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// CountUserRequests counts how many requests a single user has sent
func CountUserRequests(c *fiber.Ctx) error {
	db := middleware.DBConn

	userIDParam := c.Params("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid user ID",
			Data: errors.ErrorModel{
				Message:   "User ID must be a valid number",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	var count int64
	err = db.Model(&users.ServiceRequest{}).
		Where("user_id = ?", userID).
		Count(&count).Error

	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to count user requests",
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
		Data: fiber.Map{
			"user_id":       userID,
			"request_count": count,
		},
	})
}
