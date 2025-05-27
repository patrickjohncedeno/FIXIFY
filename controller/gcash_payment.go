package controller

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"fixify_backend/middleware"
	"fixify_backend/model/users"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
)

func InitiateXenditGCash(c *fiber.Ctx) error {
	type RequestBody struct {
		Amount float64  `json:"amount"`
		Email  string `json:"email"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// SAFELY get authenticated user from context
	userClaims := c.Locals("user")
	if userClaims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	claims, ok := userClaims.(*users.Claims)
	if !ok || claims == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token claims",
		})
	}

	// Get API Key from environment variable
	apiKey := os.Getenv("XENDIT_API_KEY")
	if apiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Xendit API key not configured",
		})
	}

	client := resty.New()
	refID := "gcash-ref-" + time.Now().Format("20060102150405")

	resp, err := client.R().
		SetBasicAuth(apiKey, "").
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"reference_id":    refID,
			"currency":        "PHP",
			"amount":          body.Amount,
			"checkout_method": "ONE_TIME_PAYMENT",
			"channel_code":    "PH_GCASH",
			"channel_properties": map[string]string{
				"success_redirect_url": os.Getenv("SUCCESS_REDIRECT_URL"),
				"failure_redirect_url": os.Getenv("FAILURE_REDIRECT_URL"),
			},
			"customer": map[string]string{
				"email": body.Email,
			},
		}).
		Post("https://api.xendit.co/ewallets/charges")

	if err != nil || resp.StatusCode() >= 400 {
		errorMsg := "GCash payment failed"
		if err != nil {
			errorMsg = fmt.Sprintf("%s: %v", errorMsg, err)
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":    errorMsg,
			"status":   resp.StatusCode(),
			"response": string(resp.Body()),
		})
	}

	var responseBody struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(resp.Body(), &responseBody); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to parse Xendit response",
		})
	}

	// Save payment to database
	payment := users.GCashPayment{
		PaymentFrom: int(claims.UserId),
		PaymentTo:   0, // Replace with the actual recipient ID if applicable
		TransactionId:     responseBody.ID,
		Amount:      body.Amount,
		GcashID:  int(.GcashID),
		PaymentDate: time.Now(),
	}

	if err := middleware.GetDB().Create(&payment).Error; err != nil {
		fmt.Println("DB Save Error:", err) // Print in logs
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(), // Return full error
		})
	}

	return c.Status(fiber.StatusOK).Send(resp.Body())
}
