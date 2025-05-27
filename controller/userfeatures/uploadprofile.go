package userfeatures

import (
	"encoding/base64"
	"io/ioutil"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"

	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"
)

// PatchProfilePicture handles updating an existing user's profile picture as BYTEA (PATCH)
func PatchProfilePicture(c *fiber.Ctx) error {
	db := middleware.DBConn
	// Retrieve user claims from the token (ensure consistency with "user" key)
	claims, ok := c.Locals("user").(*users.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Unauthorized",
			Data: errors.ErrorModel{
				Message:   "User not authenticated",
				IsSuccess: false,
				Error:     "Missing user claims",
			},
		})
	}

	// Get uploaded file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "File is required",
			Data: errors.ErrorModel{
				Message:   "No file uploaded",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	fileBytes, err := readFileBytes(fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to read file",
			Data: errors.ErrorModel{
				Message:   "Could not read uploaded file",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Patch the user's profile picture (only update the profile_picture field)
	if err := db.Model(&users.User{}).
		Where("user_id = ?", claims.UserId).
		Update("profile_picture", fileBytes).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to update profile picture",
			Data: errors.ErrorModel{
				Message:   "Database error",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Profile picture updated successfully",
		Data:    nil,
	})
}

// GetProfilePicture returns the profile picture in base64 format
func GetProfilePicture(c *fiber.Ctx) error {
	db := middleware.DBConn
	// Retrieve user claims from the token (ensure consistency with "user" key)
	claims, ok := c.Locals("user").(*users.Claims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Unauthorized",
			Data: errors.ErrorModel{
				Message:   "User not authenticated",
				IsSuccess: false,
				Error:     "Missing user claims",
			},
		})
	}

	var user users.User
	if err := db.First(&user, "user_id = ?", claims.UserId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ResponseModel{
			RetCode: "404",
			Message: "User not found",
			Data: errors.ErrorModel{
				Message:   "No user with that ID",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	if len(user.Profile_picture) == 0 {
		return c.JSON(response.ResponseModel{
			RetCode: "204",
			Message: "No profile picture uploaded",
			Data:    nil,
		})
	}

	encoded := base64.StdEncoding.EncodeToString(user.Profile_picture)
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Profile picture retrieved",
		Data: fiber.Map{
			"image_base64": encoded,
		},
	})
}

// Helper to read multipart file into []byte
func readFileBytes(fh *multipart.FileHeader) ([]byte, error) {
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}
	