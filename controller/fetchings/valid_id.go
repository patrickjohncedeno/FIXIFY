package fetchings

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

func FetchAllId(c *fiber.Ctx) error {
	db := middleware.DBConn
	var validID []users.UserVerification

	err := db.Preload("User").Find(&validID).Error
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
		Data:    validID, // ✅ Correct variable
	})
}

func FetchValid(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Extract claims from context
	claimsInterface := c.Locals("user")
	if claimsInterface == nil {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Unauthorized",
			Data: errors.ErrorModel{
				Message:   "Token claims not found",
				IsSuccess: false,
				Error:     "Missing token data",
			},
		})
	}

	claims, ok := claimsInterface.(*users.Claims)
	if !ok {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid token claims",
			Data: errors.ErrorModel{
				Message:   "Could not parse token claims",
				IsSuccess: false,
				Error:     "Type assertion failed",
			},
		})
	}

	// Get user ID from claims
	userID := claims.UserId

	// Step 2: Query the UserVerification using user_id
	var userVerification users.UserVerification
	err := db.Preload("User").Where("user_id = ?", userID).First(&userVerification).Error
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "404",
			Message: "User verification not found",
			Data: errors.ErrorModel{
				Message:   "No verification record found for this user",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    userVerification,
	})
}
