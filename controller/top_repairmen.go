package controller

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"
	"math"

	"github.com/gofiber/fiber/v2"
)

// UpdateAllRepairmanAverageRatings calculates average ratings from reviews and updates each repairman's average_rating
func TopRepairmen(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Struct to hold aggregated and joined data
	type AvgRatingResult struct {
		RepairmanID    uint    `json:"repairman_id"`
		AvgRating      float64 `json:"avg_rating"`
		FirstName      string  `json:"first_name"`
		LastName       string  `json:"last_name"`
		ProfilePicture string  `json:"profile_picture"`
		ServiceName    string  `json:"service_name"`
	}

	var results []AvgRatingResult

	// Query: Join reviews with users, repairmen, and service_categories
	err := db.Table("reviews").
		Select(`
			reviews.repairman_id,
			AVG(reviews.rating) AS avg_rating,
			users.first_name,
			users.last_name,
			service_categories.category_name AS service_name
		`).
		Joins("JOIN users ON users.user_id = reviews.repairman_id").
		Joins("JOIN service_categories ON service_categories.category_id = users.category_id").
		Where("users.type = ?", "Repairman").
		Group("reviews.repairman_id, users.first_name, users.last_name, service_categories.category_name").
		Order("avg_rating DESC").
		Scan(&results).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch average ratings",
			Data: errors.ErrorModel{
				Message:   "Database query error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Update average_rating in the repairman table
	for i, r := range results {
		roundedAvg := math.Round(r.AvgRating*10) / 10

		err := db.Model(&users.Repairman{}).
			Where("user_id = ? AND type = ?", r.RepairmanID, "Repairman").
			Update("average_rating", roundedAvg).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Failed to update average rating for user",
				Data: errors.ErrorModel{
					Message:   "Database update error",
					IsSuccess: false,
					Error:     err.Error(),
				},
			})
		}

		// Update result to reflect the rounded value
		results[i].AvgRating = roundedAvg
	}

	// Return success response with enriched data
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Average ratings updated successfully",
		Data:    results,
	})
}
