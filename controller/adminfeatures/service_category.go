package adminfeatures

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

func AddServiceCategory(c *fiber.Ctx) error {
	db := middleware.DBConn
	service := new(users.ServiceCategory)

	if err := c.BodyParser(&service); err != nil {
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

	// // Get the authenticated user
	// user := c.Locals("user").(*users.Claims)

	// Create the service request

	if err := db.Create(&service).Error; err != nil {
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

	// Return final response
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    service,
	})
}
func UpdateService(c *fiber.Ctx) error {
	db := middleware.DBConn
	service := new(users.ServiceCategory)

	// Parse URL param ID
	id := c.Params("id")
	if id == "" {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Service ID parameter is required",
			Data: errors.ErrorModel{
				Message:   "Missing ID in URL parameter",
				IsSuccess: false,
				Error:     "ID parameter is empty",
			},
		})
	}

	// Parse incoming JSON body
	if err := c.BodyParser(&service); err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Invalid Request!",
			Data: errors.ErrorModel{
				Message:   "Failed to parse request body",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Fetch existing service by ID
	var existingService users.ServiceCategory
	if err := db.First(&existingService, id).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "404",
			Message: "Service not found",
			Data: errors.ErrorModel{
				Message:   "Service with given ID does not exist",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Update only specified fields
	if err := db.Model(&existingService).Updates(service).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Cannot update service!",
			Data: errors.ErrorModel{
				Message:   "Failed to update service",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Reload updated record
	if err := db.First(&existingService, id).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to reload updated service",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch updated service",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Return updated service
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    existingService,
	})
}
