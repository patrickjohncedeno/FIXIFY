package signuplogin

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt" // Import bcrypt for password hashing comparison
	"gorm.io/gorm"
)

// UserLogin handles the login request for a user
func UserLogin(c *fiber.Ctx) error {
	db := middleware.DBConn
	logac := new(users.Repairman)

	// Parse incoming request body to the User struct
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

	// Retrieve the User record from the database based on email
	var user users.User
	if err := db.Where("email = ?", logac.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// If no matching record is found, return a login failed response
			return c.JSON(response.ResponseModel{
				RetCode: "401",
				Message: "Login Failed!",
				Data: errors.ErrorModel{
					Message:   "Invalid email/username or password",
					IsSuccess: false,
					Error:     err.Error(),
				},
			})
		}
		// Handle any other database errors
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error!",
			Data: errors.ErrorModel{
				Message:   "An error occurred while processing your request",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Compare the provided password with the stored hashed password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(logac.Password))
	if err != nil {
		// If the password doesn't match, return a login failed response
		return c.JSON(response.ResponseModel{
			RetCode: "401",
			Message: "Login Failed!",
			Data: errors.ErrorModel{
				Message:   "Invalid email/username or password",
				IsSuccess: false,
				Error:     "Password does not match",
			},
		})
	}

	// If credentials are valid, generate the JWT token
	token, err := GenerateJWT(
		int(user.UserId))
	if err != nil {
		// Handle token generation error
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Internal Server Error!",
			Data: errors.ErrorModel{
				Message:   "Error generating token",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Return success response with the generated token
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Login Successful!",
		Data: map[string]interface{}{
			"user":  user,
			"token": token, // Include the generated token
		},
	})
}
