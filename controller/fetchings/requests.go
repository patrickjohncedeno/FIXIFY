package fetchings

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

// FetchAllRequest fetches all service requests with preloaded relationships
func FetchAllRequest(c *fiber.Ctx) error {
	db := middleware.DBConn
	var request []users.ServiceRequest

	err := db.Preload("User").
		Preload("Repairman").
		Preload("ServiceCategory").
		Preload("Review").
		Find(&request).Error

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
		Data:    request,
	})
}

func FetchCompletedRequest(c *fiber.Ctx) error {
	db := middleware.DBConn
	var request []users.ServiceRequest

	err := db.Preload("User").
		Preload("Repairman").
		Preload("ServiceCategory").
		Where("status = ?", "completed").
		Find(&request).Error

	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Request failed",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch completed data",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    request,
	})
}

func FetchCanceledRequest(c *fiber.Ctx) error {
	db := middleware.DBConn
	var request []users.ServiceRequest

	err := db.Preload("User").
		Preload("Repairman").
		Preload("ServiceCategory").
		Where("status = ?", "canceled").
		Find(&request).Error

	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Request failed",
			Data: errors.ErrorModel{
				Message:   "Failed to fetch canceled data",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    request,
	})
}

// CountAllRequests counts all service requests
func CountAllRequests(c *fiber.Ctx) error {
	db := middleware.DBConn
	var count int64

	err := db.Model(&users.ServiceRequest{}).Count(&count).Error
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Request failed",
			Data: errors.ErrorModel{
				Message:   "Failed to count service requests",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    count,
	})
}
