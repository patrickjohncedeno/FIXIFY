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

func AdminLogin(c *fiber.Ctx) error {
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

	// Get email or username and password from the parsed request body
	input := logac.Email // This could be either the email or username
	password := logac.Password

	// Retrieve the admin record from the database based on either email or username
	var admin users.Admin
	// We now check both email and username
	if err := db.Where("email = ? OR username = ?", input, input).First(&admin).Error; err != nil {
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
		// If there's some other error (such as a database error), return it
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
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
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
		int(admin.AdminId))
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

	// If we get here, the email/username and password match, return success
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Login Successful!",
		Data: map[string]interface{}{
			"token": token,
			"admin": admin, // Include the generated token
		},
	})
}
