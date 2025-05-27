package signuplogin

import (
	"fixify_backend/model/users"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// Secret key to sign the token
var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// GenerateJWT generates a new JWT token for a given username
func GenerateJWT(id int) (string, error) {
	// Create the claims
	claims := users.Claims{
		UserId: uint(id),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
			Issuer:    "Fixkify",                             // Issuer of the token
		},
	}

	// Create the token using the claims and secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token and get the string representation
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWT validates the JWT token and returns the username if valid
func ValidateJWT(tokenString string) (string, error) {
	// Parse the token and validate it
	token, err := jwt.ParseWithClaims(tokenString, &users.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	// Check for errors and if the token is valid
	if err != nil {
		return "", err
	}

	// Extract the claims from the token
	if claims, ok := token.Claims.(*users.Claims); ok && token.Valid {
		return fmt.Sprintf("%d", claims.UserId), nil
	} else {
		return "", fmt.Errorf("invalid token")
	}
}

// JWTLogin handles login requests and generates a JWT token
func JWTLogin(c *fiber.Ctx) error {
	logac := new(users.Claims)
	if err := c.BodyParser(&logac); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// Generate JWT Token
	token, err := GenerateJWT(int(logac.UserId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error generating token")
	}

	// Return the token in the response
	return c.JSON(fiber.Map{
		"token": token,
	})
}

// JWTValidate validates a JWT token and returns the username if valid
func JWTValidate(c *fiber.Ctx) error {
	// Get the token from the Authorization header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Authorization token is required")
	}

	// Remove the "Bearer " prefix from the token
	tokenString = tokenString[len("Bearer "):]

	// Validate the token
	email, err := ValidateJWT(tokenString)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid token")
	}

	// Return the username from the token
	return c.JSON(fiber.Map{
		"username": email,
	})
}

// JWTMiddleware is a middleware function that checks the presence and validity of the JWT token
func JWTMiddleware(c *fiber.Ctx) error {
	// Get the token from the Authorization header
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization token is required",
		})
	}

	// Remove the "Bearer " prefix from the token if it's there
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	// Parse and validate the token
	claims := &users.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	// Check for any errors
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired token",
		})
	}

	// Store the claims (user data) in the context for use in the handler
	c.Locals("user", claims)

	return c.Next() // Continue processing the request
}

// ParseJWTClaims extracts and validates claims from a JWT token string
func ParseJWTClaims(tokenString string, claims *users.Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
}
