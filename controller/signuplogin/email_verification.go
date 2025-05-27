package signuplogin

import (
	"fixify_backend/middleware"
	errors "fixify_backend/model/error" // Import the errors package
	"fixify_backend/model/response"
	"fixify_backend/model/users"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gopkg.in/gomail.v2"
)

// Function to generate a random verification code
func GenerateVerificationCode() int {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999-100000) + 100000 // Generate a 6-digit code
	return (code)
}

// Function to send the verification email
func SendVerificationEmail(userEmail string, verificationCode int) error {
	from := os.Getenv("FROM")            // Your Gmail address
	password := os.Getenv("APPASS")      // Your Gmail App Passwords
	smtpHost := os.Getenv("SMTPHOST")    // Gmail's SMTP server
	smtpPortStr := os.Getenv("SMTPPORT") // Port for TLS (recommended)

	// Convert the SMTP port string to an integer
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatal("Error converting SMTP port to integer:", err)
	}

	// Convert the verification code to a string
	verificationCodeStr := strconv.Itoa(verificationCode)

	// Create the email message with the verification code
	message := gomail.NewMessage()
	message.SetHeader("From", from)                                                 // Sender's email
	message.SetHeader("To", userEmail)                                              // Recipient's email
	message.SetHeader("Subject", "Email Verification Code")                         // Email subject
	message.SetBody("text/plain", "Your verification code is "+verificationCodeStr) // Email body
	log.Println(verificationCode)
	// Set up the SMTP client
	dialer := gomail.NewDialer(smtpHost, smtpPort, from, password)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		log.Println("Failed to send email:", err)
		return err
	}

	log.Println("Verification email sent successfully!")
	return nil
}

// User signup function with email verification (no password or other fields)
func EmailVer(c *fiber.Ctx) error {
	db := middleware.DBConn
	logac := new(users.EmailVer)

	// Parse incoming request body to the EmailVer struct (only email field)
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

	// Check if the email is already in use
	var existingUserEmail users.EmailVer
	if err := db.Where("email = ?", logac.Email).First(&existingUserEmail).Error; err == nil {
		// Generate the verification code
		code := GenerateVerificationCode()

		// Save the email and verification code in the database
		emailVer := &users.EmailVer{
			Email: logac.Email,
			Code:  strconv.Itoa(code),
		}
		if err := db.Model(&users.EmailVer{}).Where("email = ?", emailVer.Email).Update("code", emailVer.Code).Error; err != nil {
			return c.JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Error saving verification code",
				Data: errors.ErrorModel{
					Message:   "Failed to save verification code!",
					IsSuccess: false,
					Error:     err.Error(),
				},
			})
		}

		// Send the verification email
		err := SendVerificationEmail(logac.Email, code)
		if err != nil {
			return c.JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Failed to send verification email",
				Data: errors.ErrorModel{
					Message:   "Error sending verification email",
					IsSuccess: false,
					Error:     err.Error(),
				},
			})
		}

		// Set a 3-minute timer to automatically delete the verification code from the database
		go func() {
			db1 := middleware.DBConn
			// Wait for 3 minutes (60 seconds)
			time.AfterFunc(5*time.Minute, func() {
				// Delete the email verification code after the delay
				if err := db1.Where("email = ?", logac.Email).Delete(&users.EmailVer{}).Error; err != nil {
					// Log an error if the deletion fails
					// Replace with your preferred logging method
					fmt.Println("Error deleting verification code:", err)
				} else {
					fmt.Println("Verification code deleted successfully.") // Success message
				}
			})
		}()
	} else {
		// Generate the verification code
		code := GenerateVerificationCode()

		// Save the email and verification code in the database
		emailVer := &users.EmailVer{
			Email: logac.Email,
			Code:  strconv.Itoa(code),
		}
		if err := db.Create(emailVer).Error; err != nil {
			return c.JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Error saving verification code",
				Data: errors.ErrorModel{
					Message:   "Failed to save verification code!",
					IsSuccess: false,
					Error:     err.Error(),
				},
			})
		}

		// Send the verification email
		err := SendVerificationEmail(logac.Email, code)
		if err != nil {
			return c.JSON(response.ResponseModel{
				RetCode: "500",
				Message: "Failed to send verification email",
				Data: errors.ErrorModel{
					Message:   "Error sending verification email",
					IsSuccess: false,
					Error:     err.Error(),
				},
			})
		}

		// Set a 3-minute timer to automatically delete the verification code from the database
		go func() {
			db1 := middleware.DBConn
			// Wait for 3 minutes (60 seconds)
			time.AfterFunc(5*time.Minute, func() {
				// Delete the email verification code after the delay
				if err := db1.Where("email = ?", logac.Email).Delete(&users.EmailVer{}).Error; err != nil {
					// Log an error if the deletion fails
					// Replace with your preferred logging method
					fmt.Println("Error deleting verification code:", err)
				} else {
					fmt.Println("Verification code deleted successfully.") // Success message
				}
			})
		}()
	}

	// Return the success response
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Email sign-up successful! A verification email has been sent.",
		Data:    logac,
	})
}

// Function to verify the email verification code
func EmailVerCode(c *fiber.Ctx) error {
	db := middleware.DBConn
	logac := new(users.EmailVer)

	// Parse incoming request body to the EmailVer struct (email and code fields)
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

	// Check if the code matches
	var existingCode users.EmailVer
	if err := db.Where("code = ?", logac.Code).First(&existingCode).Error; err != nil {
		// If no matching record is found, return an error
		return c.JSON(response.ResponseModel{
			RetCode: "400",
			Message: "Invalid Verification Code!",
			Data: errors.ErrorModel{
				Message:   "The provided verification code is invalid",
				IsSuccess: false,
				Error:     err.Error(),
			},
		})
	}

	// Automatically delete the verification code and email from the database after verification
	go func() {
		fmt.Println("Starting the deletion process...") // Debugging message

		// Delete the email verification code after the delay
		if err := db.Where("code = ?", logac.Code).Delete(&users.EmailVer{}).Error; err != nil {
			// Log an error if the deletion fails
			// Replace with your preferred logging method
			fmt.Println("Error deleting verification code:", err)
		} else {
			fmt.Println("Verification code deleted successfully.") // Success message
		}

	}()

	// If the email and code match, return a success response
	return c.JSON(response.ResponseModel{
		RetCode: "200",
		Message: "Verification Successful!",
		Data:    logac,
	})
}
