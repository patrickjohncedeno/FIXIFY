package fetchings

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

// FetchAllRepairmans fetches all repairmen from the database
func FetchAllRepairmen(c *fiber.Ctx) error {
	db := middleware.DBConn
	var repairman []users.Repairman

	// Raw SQL query to fetch all repairmen
	err := db.Preload("ServiceCategory").Where("type = ?", "Repairman").Find(&repairman).Error
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
		Data:    repairman, // Return the raw data
	})
}

// CountAllRepairmen counts all repairmen (users with type = 'Repairman')
func CountAllRepairmen(c *fiber.Ctx) error {
	db := middleware.DBConn
	var count int64

	// Count users where type = 'Repairman'
	err := db.Model(&users.Repairman{}).Where("type = ?", "Repairman").Count(&count).Error
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Request failed",
			Data: errors.ErrorModel{
				Message:   "Failed to count repairmen from database",
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
