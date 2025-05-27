package adminfeatures

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func VerifyUser(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Parse user_id from URL params
	userIDParam := c.Params("id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid user ID",
			Data: errors.ErrorModel{
				Message:   "User ID must be a valid number",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Define request body
	type RequestBody struct {
		Status string `json:"status"` // e.g., "pending", "approved", "rejected"
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid request body",
			Data: errors.ErrorModel{
				Message:   "Failed to parse JSON body",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Check if user verification record exists
	var userVer users.UserVerification
	if err := db.Where("user_id = ?", userID).First(&userVer).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ResponseModel{
			RetCode: "404",
			Message: "User verification record not found",
			Data: errors.ErrorModel{
				Message:   "No user verification entry found for the given user ID",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Update the verification status
	err = db.Model(&users.UserVerification{}).
		Where("user_id = ?", userID).
		Update("status", body.Status).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to update verification status",
			Data: errors.ErrorModel{
				Message:   "Database update error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "User verification status updated successfully",
		Data: fiber.Map{
			"user_id": userID,
			"status":  body.Status,
		},
	})
}
