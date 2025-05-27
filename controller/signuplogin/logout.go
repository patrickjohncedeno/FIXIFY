package signuplogin

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

// JWTLogout handles logging out by clearing the JWT token from the client's cookies
func JWTLogout(c *fiber.Ctx) error {

	// Clear the token cookie by setting it to expire in the past
	c.Cookie(&fiber.Cookie{
		Name:    "token",                    // The name of your cookie (make sure this matches the cookie you're using for the token)
		Value:   "",                         // Empty value to clear it
		Expires: time.Now().Add(-time.Hour), // Expire the cookie immediately
		Path:    "/",
		// Optional: You can set a message in the cookie (removed unkeyed field to fix the error)
	})
	// Return a success response to the user
	return c.JSON(fiber.Map{
		"message": "Logged out successfully", // Inform the client about the successful logout
	})
}
