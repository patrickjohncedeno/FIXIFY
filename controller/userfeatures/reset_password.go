package userfeatures

import (
	"fixify_backend/controller/signuplogin"
	"fixify_backend/middleware"
	errors "fixify_backend/model/error"
	"fixify_backend/model/response"
	"fixify_backend/model/users"

	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ResetPassword(c *fiber.Ctx) error {
	db := middleware.DBConn

	// Parse user ID
	userIdParam := c.Params("id")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil || userId <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid User ID!",
			Data: errors.ErrorModel{
				Message:   "User ID must be a valid number",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Request body: expects email and new password
	var body struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid request body!",
			Data: errors.ErrorModel{
				Message:   "Failed to parse body",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Check if the user exists and email matches
	var user users.User
	if err := db.Where("user_id = ? AND email = ?", userId, body.Email).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.ResponseModel{
			RetCode: "404",
			Message: "User not found",
			Data: errors.ErrorModel{
				Message:   "User ID and email mismatch",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Generate verification code
	code := signuplogin.GenerateVerificationCode()

	// Store the verification code in EmailVer table (upsert style)
	emailVer := &users.EmailVer{
		Email: body.Email,
		Code:  strconv.Itoa(code),
	}
	db.Where("email = ?", body.Email).Delete(&users.EmailVer{})
	if err := db.Create(emailVer).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Failed to store verification code",
			Data: errors.ErrorModel{
				Message:   "DB insert failed",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Send the email with the code
	if err := signuplogin.SendVerificationEmail(body.Email, code); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.ResponseModel{
			RetCode: "500",
			Message: "Email sending failed",
			Data: errors.ErrorModel{
				Message:   "Could not send verification email",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Auto-delete code after 5 minutes
	go func(email string) {
		time.Sleep(5 * time.Minute)
		db1 := middleware.DBConn
		db1.Where("email = ?", email).Delete(&users.EmailVer{})
	}(body.Email)

	// Return success
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Verification code sent to your email. Please verify to reset your password.",
		Data:    body.Email,
	})
}
