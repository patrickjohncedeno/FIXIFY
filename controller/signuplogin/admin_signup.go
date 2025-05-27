package signuplogin

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

func AdminSignup(c *fiber.Ctx) error {
	db := middleware.DBConn
	logac := new(users.Admin)

	// Parse incoming request body to the Admin struct
	if err := c.BodyParser(logac); err != nil {
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

	// Log received admin data (for debugging)
	print("Admin Sign-up", logac.Username, "\n", logac.Password, "\n", logac.Email, "\n")

	// Check if the email is already in use
	var existingAdminByEmail users.Admin
	if err := db.Where("email = ?", logac.Email).First(&existingAdminByEmail).Error; err == nil {
		// If an admin with the same email exists, return an error
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Email already in use!",
			Data: errors.ErrorModel{
				Message:   "The provided email is already associated with an existing account",
				IsSuccess: false,
				Error:     "Email already in use",
			},
		})
	}

	// Hash the admin password before saving it to the database
	hashedPassword, err := middleware.HashPassword(logac.Password)
	if err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to hash password",
			Data: errors.ErrorModel{
				Message:   "Error hashing password",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}
	// Store the hashed password
	logac.Password = hashedPassword

	// Save the admin in the database (this will insert into the 'admins' table or whatever your table is called)
	if err := db.Create(logac).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Cannot sign up!",
			Data: errors.ErrorModel{
				Message:   "Failed to sign up!",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Return the success response with the admin data
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Admin sign-up successful!",
		Data:    logac,
	})
}
