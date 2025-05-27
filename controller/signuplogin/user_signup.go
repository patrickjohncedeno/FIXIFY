package signuplogin

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
)

func UserSignup(c *fiber.Ctx) error {
	db := middleware.DBConn
	logac := new(users.User)

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

	// Check if the email is already in use
	var existingUserEmail users.User
	if err := db.Where("email = ?", logac.Email).First(&existingUserEmail).Error; err == nil {
		// If an user with the same email exists, return an error
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

	// Check if the phone is already in use
	var existingUserPhone users.User
	if err := db.Where("phone = ?", logac.Phone).First(&existingUserPhone).Error; err == nil {
		//if user with same phone number exists, return an error
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Phone number already in use!",
			Data: errors.ErrorModel{
				Message: "The provided phone number is already associated with an existing account",
			},
		})
	}
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
		Message: "User sign-up successful!",
		Data:    logac,
	})
}

func RepairmanSignup(c *fiber.Ctx) error {
	db := middleware.DBConn
	logac := new(users.Repairman)

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

	// Check if the email is already in use
	var existingUserEmail users.EmailVer
	if err := db.Where("email = ?", existingUserEmail.Email).First(&existingUserEmail).Error; err == nil {
		// If an user with the same email exists, return an error
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

	// Check if the phone is already in use
	var existingUserPhone users.Repairman
	if err := db.Where("phone = ?", logac.Phone).First(&existingUserPhone).Error; err == nil {
		//if user with same phone number exists, return an error
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Phone number already in use!",
			Data: errors.ErrorModel{
				Message: "The provided phone number is already associated with an existing account",
			},
		})
	}
	// Save the repairman in the database (this will insert into the 'users' table or whatever your table is called)
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
		Message: "User sign-up successful!",
		Data:    logac,
	})
}
