package userfeatures

import (
	"fixify_backend/controller"
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ReviewBody struct {
	Rating     int    `json:"rating"`
	ReviewText string `json:"review_text"`
}

func ReviewRequest(c *fiber.Ctx) error {
	db := middleware.DBConn
	user := c.Locals("user").(*users.Claims)

	requestIdParam := c.Params("id")
	requestId, err := strconv.Atoi(requestIdParam)
	if err != nil || requestId <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid request ID",
			Data: errors.ErrorModel{
				Message:   "Request ID must be a valid number",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	var body ReviewBody
	if err := c.BodyParser(&body); err != nil || body.Rating <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid review body",
			Data: errors.ErrorModel{
				Message:   "Rating and review text are required",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	var serviceRequest users.ServiceRequest
	if err := db.First(&serviceRequest, "request_id = ? AND user_id = ?", requestId, user.UserId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Service request not found or not yours",
			Data: errors.ErrorModel{
				Message:   "Cannot find request for this user",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	if serviceRequest.ReviewId != 0 {
		return c.Status(fiber.StatusForbidden).JSON(response.ResponseModel{
			RetCode: "403",
			Message: "Request already reviewed",
			Data: errors.ErrorModel{
				Message:   "This request has already been reviewed",
				IsSuccess: false,
			},
		})
	}

	if serviceRequest.Status != "completed" {
		return c.Status(fiber.StatusForbidden).JSON(response.ResponseModel{
			RetCode: "403",
			Message: "You can only review completed requests",
			Data: errors.ErrorModel{
				Message:   "Service request must be completed before review",
				IsSuccess: false,
			},
		})
	}

	review := users.Review{
		RequestId:   uint(requestId),
		ClientId:    user.UserId,
		RepairmanId: serviceRequest.RepairmanId,
		Rating:      float64(body.Rating),
		ReviewText:  body.ReviewText,
		ReviewDate:  users.TimeWithDate(time.Now()),
	}

	if err := db.Create(&review).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Could not save review",
			Data: errors.ErrorModel{
				Message:   "Database error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	if err := db.Model(&users.ServiceRequest{}).
		Where("request_id = ?", requestId).
		Update("review_id", review.ReviewId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to update service request with review",
			Data: errors.ErrorModel{
				Message:   "Could not associate review with the request",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	var client users.User
	if err := db.First(&client, "user_id = ?", user.UserId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Client fetch error",
			Data: errors.ErrorModel{
				Message:   "Failed to find client info",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// === Call centralized notification creator ===
	if err := controller.CreateUserNotification(
		db,
		"Service Review",
		requestId,
		int(user.UserId),
		int(serviceRequest.RepairmanId),
		"You received a new review from "+client.First_name,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to create notification",
			Data: errors.ErrorModel{
				Message:   "Notification creation error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// === Average Rating Update ===
	type AvgRatingResult struct {
		RepairmanID uint    `json:"repairman_id"`
		AvgRating   float64 `json:"avg_rating"`
	}

	var results []AvgRatingResult
	err = db.Table("reviews").
		Select("repairman_id, AVG(rating) as avg_rating").
		Group("repairman_id").
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

	for _, r := range results {
		roundedAvg := math.Round(r.AvgRating*10) / 10 // Round to 1 decimal
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
	}

	var fullReview users.Review
	if err := db.Preload("Client").
		Preload("Repairman").
		Preload("Request").
		First(&fullReview, review.ReviewId).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Could not fetch review with details",
			Data: errors.ErrorModel{
				Message:   "Failed to load related data",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Review submitted and ratings updated successfully",
		Data:    fullReview,
	})
}
