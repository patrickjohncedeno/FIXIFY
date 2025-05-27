package fetchings

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

func FetchServices(c *fiber.Ctx) error {
	db := middleware.DBConn
	var request []users.ServiceCategory

	if err := db.Find(&request).Error; err != nil {
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

	// Return the fetched requests if found
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Success",
		Data:    request, // Return the found requests
	})
}

func DeleteServiceCategory(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Get category_id from route params
	categoryID := c.Params("id") // assumes the route is like /categories/:id

	if categoryID == "" {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request",
			Data: errors.ErrorModel{
				Message:   "Category ID is required",
				IsSuccess: false,
				Error:     "Missing category ID in request",
			},
		})
	}

	// Delete the category from the database
	if err := db.Delete(&users.ServiceCategory{}, categoryID).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to delete category",
			Data: errors.ErrorModel{
				Message:   "Database error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Category deleted successfully",
		Data:    nil,
	})
}

func DisableServiceCategory(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Get category_id from route params
	categoryID := c.Params("id")

	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Bad Request",
			Data: errors.ErrorModel{
				Message:   "Category ID is required",
				IsSuccess: false,
				Error:     "Missing category ID in request",
			},
		})
	}

	// Update is_active = false where category_id = :id
	if err := db.Model(&users.ServiceCategory{}).
		Where("category_id = ?", categoryID).
		Update("is_active", false).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to disable category",
			Data: errors.ErrorModel{
				Message:   "Database error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Category disabled successfully",
		Data:    nil,
	})
}
