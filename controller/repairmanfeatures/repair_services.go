package repairmanfeatures

import (
	"encoding/json"
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

func UpdateRepairmanCategories(c *fiber.Ctx) error {
	db := middleware.DBConn

	user := c.Locals("user").(*users.Claims)
	repairmanID := user.UserId

	var body struct {
		CategoryIDs int `json:"category_ids"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid request body",
			Data: errors.ErrorModel{
				Message:   "Could not parse request",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	jsonData, err := json.Marshal(body.CategoryIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to convert category_ids to JSON",
			Data: errors.ErrorModel{
				Message:   "JSON marshal error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	if err := db.Table("users").
		Where("user_id = ?", repairmanID).
		Update("category_id", jsonData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to update categories",
			Data: errors.ErrorModel{
				Message:   "Error updating category_id column",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Service categories updated successfully",
		Data:    body.CategoryIDs,
	})
}
