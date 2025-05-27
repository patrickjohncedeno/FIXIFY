package fetchings

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func FetchAllReviews(c *fiber.Ctx) error {
	db := middleware.DBConn 
	var allReviews []users.Review

	err := db.Find(&allReviews).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch reviews",
			Data: errors.ErrorModel{
				Message:   "Database query error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    allReviews,
	})
}

type ReviewBody struct {
	Rating     int    `json:"rating"`
	ReviewText string `json:"review_text"`
}

func ReviewRequest(c *fiber.Ctx) error {
	db := middleware.DBConn

	userIdParam := c.Params("id")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil || userId <= 0 {
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

	var reviews []users.Review
	err = db.Preload("Client").
		Preload("Repairman").
		Preload("Request").
		Where("repairman_id = ?", userId).
		Find(&reviews).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Could not fetch reviews",
			Data: errors.ErrorModel{
				Message:   "Database error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Fetched reviews successfully",
		Data:    reviews,
	})
}
