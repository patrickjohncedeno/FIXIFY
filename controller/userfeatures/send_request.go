package userfeatures

import (
	"fixify_backend/controller"
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type DescriptionBody struct {
	Description string `json:"description"`
}

func ServiceRequest(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Get repairman ID from URL
	repairmanIdParam := c.Params("id")
	repairmanId, err := strconv.Atoi(repairmanIdParam)
	if err != nil || repairmanId <= 0 {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid Repairman ID!",
			Data: errors.ErrorModel{
				Message:   "Repairman ID must be a valid number",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Parse request body
	var body DescriptionBody
	if err := c.BodyParser(&body); err != nil || body.Description == "" {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request!",
			Data: errors.ErrorModel{
				Message:   "Description is required",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Get the authenticated user
	user := c.Locals("user").(*users.Claims)

	// Check if user exists
	var client users.User
	if err := db.First(&client, "user_id = ?", user.UserId).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Client not found!",
			Data: errors.ErrorModel{
				Message:   "Client not found",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Fetch repairman and their category_id
	var repairman users.Repairman
	if err := db.First(&repairman, "user_id = ? AND type = ?", repairmanId, "Repairman").Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Repairman not found!",
			Data: errors.ErrorModel{
				Message:   "Repairman not found",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Create the service request
	request := &users.ServiceRequest{
		UserId:      user.UserId,
		RepairmanId: uint(repairmanId),
		CategoryId:  repairman.CategoryId,
		Description: body.Description,
		Status:      "pending",
	}

	if err := db.Create(&request).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Cannot create request!",
			Data: errors.ErrorModel{
				Message:   "Failed to send request",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}
	// Preload relations for the response
	var fullRequest users.ServiceRequest
	if err := db.Preload("User").
		Preload("Repairman").
		Preload("ServiceCategory").
		First(&fullRequest, request.RequestId).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to fetch service request",
			Data: errors.ErrorModel{
				Message:   "Error fetching full request",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// ...

	// After successfully creating the request and preloading it
	notificationDescription := "You have received a new service request from " + client.First_name + "."

	if err := controller.CreateUserNotification(
		db,
		"Service Request",
		int(fullRequest.RequestId),
		int(user.UserId),
		int(fullRequest.RepairmanId),
		notificationDescription,
	); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Notification creation failed!",
			Data: errors.ErrorModel{
				Message:   "Failed to notify repairman",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    fullRequest,
	})

}
