package userfeatures

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func UpdateAccount(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Get the user ID from the route parameter
	userIdParam := c.Params("id")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil || userId <= 0 {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid User ID!",
			Data: errors.ErrorModel{
				Message:   "User ID must be a valid number",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Create a variable to hold the data received in the request body
	update := new(users.Repairman)
	// Parse the body of the request to get the update data
	if err := c.BodyParser(update); err != nil {
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

	// Create a map to hold the fields that will be updated
	updates := map[string]interface{}{}

	// Add fields to the map only if they are provided in the request
	if update.First_name != "" {
		updates["first_name"] = update.First_name
	}
	if update.Last_name != "" {
		updates["last_name"] = update.Last_name
	}
	if update.Phone != "" {
		updates["phone"] = update.Phone
	}
	if update.Address != "" {
		updates["address"] = update.Address
	}
	if update.Availability != "" {
		updates["availability"] = update.Availability
	}
	if update.Password != "" {
		updates["password"] = update.Password
	}
	// If no fields were provided for update, return an error
	if len(updates) == 0 {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid Request!",
			Data: errors.ErrorModel{
				Message:   "At least one field must be provided for update",
				IsSuccess: false,
				Error:     "No fields to update",
			},
		})
	}

	// Update only the fields provided
	result := db.Model(&users.Repairman{}).
		Where("user_id = ?", userId).
		Updates(updates)

	if result.Error != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to Update User",
			Data: errors.ErrorModel{
				Message:   "Database update failed",
				IsSuccess: false,
				Error:     result.Error.Error(),
			},
		})
	}

	if result.RowsAffected == 0 {
		return c.JSON(response.ResponseModel{
			RetCode: "404",
			Message: "User Not Found",
			Data: errors.ErrorModel{
				Message:   "No user found with the given ID",
				IsSuccess: false,
				Error:     "No rows affected",
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Account updated successfully",
		Data: errors.ErrorModel{
			Message:   "Account updated successfully",
			IsSuccess: true,
			Error:     "",
		},
	})
}
func UpdateAPassword(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Get the user ID from the route parameter
	userIdParam := c.Params("id")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil || userId <= 0 {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid User ID!",
			Data: errors.ErrorModel{
				Message:   "User ID must be a valid number",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Parse the body of the request
	type Password struct {
		OldPass string `json:"oldpass"`
		NewPass string `json:"newpass"`
	}
	// Create a variable to hold the data received in the request body
	updatepass := new(Password)
	if err := c.BodyParser(updatepass); err != nil {
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
	// Create a map to hold the fields that will be updated
	updates := map[string]interface{}{}

	// Fetch the existing user
	var existingUser users.User
	if err := db.Where("user_id = ?", userId).First(&existingUser).Error; err != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "404",
			Message: "User Not Found",
			Data: errors.ErrorModel{
				Message:   "No user found with the given ID",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}
	if updatepass.OldPass != "" {
		// Compare the provided password with the stored hashed password
		err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(updatepass.OldPass))
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
	}

	// If new password is provided, hash and update
	if updatepass.NewPass != "" {
		hashedPassword, err := middleware.HashPassword(updatepass.NewPass)
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
		updates["password"] = hashedPassword
	}

	// If no update fields provided
	if len(updates) == 0 {
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid Request!",
			Data: errors.ErrorModel{
				Message:   "At least one field must be provided for update",
				IsSuccess: false,
				Error:     "No fields to update",
			},
		})
	}

	// Perform the update
	result := db.Model(&users.User{}).Where("user_id = ?", userId).Updates(updates)
	if result.Error != nil {
		return c.JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to Update User",
			Data: errors.ErrorModel{
				Message:   "Database update failed",
				IsSuccess: false,
				Error:     result.Error.Error(),
			},
		})
	}
	if result.RowsAffected == 0 {
		return c.JSON(response.ResponseModel{
			RetCode: "404",
			Message: "User Not Found",
			Data: errors.ErrorModel{
				Message:   "No user found with the given ID",
				IsSuccess: false,
				Error:     "No rows affected",
			},
		})
	}

	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Account updated successfully",
		Data: errors.ErrorModel{
			Message:   "Account updated successfully",
			IsSuccess: true,
			Error:     "",
		},
	})
}
