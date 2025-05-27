package repairmanfeatures

import (
	"fixify_backend/controller"
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"
	"fixify_backend/websocketclient"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func RequestUpdate(c *fiber.Ctx) error {
	db := middleware.DBConn

	type UpdateStatusRequest struct {
		Status string `json:"status"`
	}

	idParam := c.Params("id")
	if idParam == "" {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid Request!",
			Data: errors.ErrorModel{
				Message:   "Request ID is required",
				IsSuccess: false,
				Error:     "Request ID parameter is missing",
			},
		})
	}

	requestId, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid Request!",
			Data: errors.ErrorModel{
				Message:   "Invalid Request ID format",
				IsSuccess: false,
				Error:     "Request ID is not a valid number",
			},
		})
	}

	update := new(UpdateStatusRequest)
	if err := c.BodyParser(update); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request!",
			Data: errors.ErrorModel{
				Message:   "Failed to parse request",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	if update.Status == "" {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid Request!",
			Data: errors.ErrorModel{
				Message:   "Status is required",
				IsSuccess: false,
				Error:     "Status cannot be empty",
			},
		})
	}

	var request users.ServiceRequest
	if err := db.Preload("User").Preload("Repairman").First(&request, requestId).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Request Not Found",
			Data: errors.ErrorModel{
				Message:   "Service request not found",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	request.Status = update.Status
	if err := db.Save(&request).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to Update Request",
			Data: errors.ErrorModel{
				Message:   "Failed to update status",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}
	var notificationDescription string
	var conversationId uint = 0

	repairmanName := "the repairman"
	if request.Repairman.UserId != 0 {
		repairmanName = request.Repairman.First_name
	}

	switch update.Status {
	case "in progress":
		notificationDescription = "Good news! Your service request has been accepted by " + repairmanName + ". They will contact you shortly to schedule the service."

		convID, err := websocketclient.EnsureClientRepairmanConversation(request.UserId, request.RepairmanId)
		if err != nil {
			log.Printf("Failed to ensure client-repairman conversation: %v", err)
		} else {
			conversationId = convID
		}

	case "completed":
		notificationDescription = "Your service request has been marked as completed by " + repairmanName + ". Thank you for using our service!"

	case "canceled":
		notificationDescription = "Unfortunately, your service request was canceled by " + repairmanName + ". Please feel free to request another service."

	default:
		notificationDescription = "The status of your service request has been updated by " + repairmanName + "."
	}

	if err := controller.CreateUserNotification(
		db,
		"Request Response",
		int(request.RequestId),
		int(request.RepairmanId),
		int(request.UserId),
		notificationDescription,
	); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to create notification!",
			Data: errors.ErrorModel{
				Message:   "Notification creation error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Request status updated and notification sent successfully",
		Data: fiber.Map{
			"message":         "Status updated and notification created",
			"isSuccess":       true,
			"error":           "",
			"client_id":       request.UserId,
			"repairman_id":    request.RepairmanId,
			"conversation_id": conversationId, // Only set when in-progress
		},
	})
}
