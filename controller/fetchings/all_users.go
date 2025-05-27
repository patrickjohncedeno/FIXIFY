package fetchings

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

func FetchAllUsers(c *fiber.Ctx) error {
	db := middleware.DBConn
	var userlist []users.Repairman

	// Raw SQL query to fetch all users
	err := db.Raw("SELECT * FROM users").Scan(&userlist).Error
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
		Data:    userlist, // Return the raw data
	})
}

func FetchUser(c *fiber.Ctx) error {
	db := middleware.DBConn
	userId := c.Params("id") // Get ID from URL params
	var user users.Repairman

	// Fetch a single user by ID
	err := db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "404", // Not found
			Message: "User not found",
			Data: errors.ErrorModel{
				Message:   "No user with the given ID",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Return the fetched user
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    user,
	})
}
