package fetchings

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

func FetchAllAdmin(c *fiber.Ctx) error {
	db := middleware.DBConn
	var admins []users.Admin

	// Raw SQL query to fetch all repairmen
	err := db.Raw("SELECT * FROM admins").Scan(&admins).Error
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
		Data:    admins, // Return the raw data
	})
}

// CountAllClients counts all clients (users with type = 'User')
func CountAllAdmin(c *fiber.Ctx) error {
	db := middleware.DBConn
	var count int64

	// Count users where type = 'User'
	err := db.Model(&users.Admin{}).Count(&count).Error
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Request failed",
			Data: errors.ErrorModel{
				Message:   "Failed to count clients from database",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    count, // Return the count
	})
}
