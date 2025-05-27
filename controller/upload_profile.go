package controller

import (
	"fixify_backend/middleware"
	"fixify_backend/model/users"
	"io/ioutil"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
)

// UpdateProfilePicture uploads and stores image binary in the database
func UpdateProfilePicture(c *fiber.Ctx) error {
	// Get user ID from JWT middleware (set in Locals)
	userIDInterface := c.Locals("user")
	userID, ok := userIDInterface.(int)
	if !ok || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"retCode": "401",
			"message": "Unauthorized: Invalid user ID in token",
			"data":    fiber.Map{"IsSuccess": false},
		})
	}

	// Parse uploaded file
	fileHeader, err := c.FormFile("profile_picture")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"retCode": "400",
			"message": "Failed to parse uploaded file",
			"data":    fiber.Map{"IsSuccess": false, "error": err.Error()},
		})
	}

	// Read file content into byte slice
	fileBytes, err := readFileBytes(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"retCode": "500",
			"message": "Failed to read uploaded file",
			"data":    fiber.Map{"IsSuccess": false, "error": err.Error()},
		})
	}

	// Save to database
	if err := UpdateUserProfilePictureBinary(uint(userID), fileBytes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"retCode": "500",
			"message": "Failed to update profile picture in database",
			"data":    fiber.Map{"IsSuccess": false, "error": err.Error()},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"retCode": "200",
		"message": "Profile picture updated successfully",
		"data":    fiber.Map{"IsSuccess": true},
	})
}

func readFileBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}
func UpdateUserProfilePictureBinary(userID uint, data []byte) error {
	db := middleware.DBConn
	logac := new(users.User)

	return db.Model(&logac).
		Where("user_id = ?", userID).
		Update("profile_picture", data).Error
}
