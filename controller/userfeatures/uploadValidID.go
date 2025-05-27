package userfeatures

import (
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v2"

	"fixify_backend/middleware"
	"fixify_backend/model/users"
)

// UploadIDCardAndSelfie uploads both the user's ID card and selfie with ID to UserVerification
// UploadIDCardAndSelfie uploads both the user's ID card and selfie with ID to UserVerification
func UploadIDCardAndSelfie(c *fiber.Ctx) error {
	db := middleware.DBConn
	claims, ok := c.Locals("user").(*users.Claims)
	if !ok {
		return unauthorizedResponse()
	}

	// Retrieve both files from the request
	idCardFile, err := c.FormFile("valid_id") // Expecting 'id_card' field in the request
	if err != nil {
		return badRequestResponse("ID card is required", err)
	}

	selfieFile, err := c.FormFile("selfie") // Expecting 'selfie' field in the request
	if err != nil {
		return badRequestResponse("Selfie with ID is required", err)
	}

	backIdFile, err := c.FormFile("back_id") // Expecting 'selfie' field in the request
	if err != nil {
		return badRequestResponse("Back ID is required", err)
	}

	// Read the bytes of the files
	idCardBytes, err := readFileBytes(idCardFile)
	if err != nil {
		return serverErrorResponse("Failed to read ID card file", err)
	}

	selfieBytes, err := readFileBytes(selfieFile)
	if err != nil {
		return serverErrorResponse("Failed to read selfie file", err)
	}

	backIdByte, err := readFileBytes(backIdFile)
	if err != nil {
		return serverErrorResponse("Failed to read Back ID file", err)
	}

	// Check if the user already has a verification record
	var verification users.UserVerification
	if err := db.Where("user_id = ?", claims.UserId).First(&verification).Error; err != nil {
		// Create a new record if it doesn't exist
		verification = users.UserVerification{
			UserId:      claims.UserId,
			ValidId:     idCardBytes,
			Selfie:      selfieBytes,
			BackId:      backIdByte,
			Status:      "pending",
			SubmittedAt: time.Now(),
		}
		if err := db.Create(&verification).Error; err != nil {
			return serverErrorResponse("Failed to create verification record", err)
		}
	} else {
		// Update the existing record if it exists
		verification.ValidId = idCardBytes
		verification.Selfie = selfieBytes
		verification.SubmittedAt = time.Now()
		if err := db.Save(&verification).Error; err != nil {
			return serverErrorResponse("Failed to update verification record", err)
		}
	}

	// Encode the uploaded ID card and selfie to base64 for response
	encodedIDCard := base64.StdEncoding.EncodeToString(idCardBytes)
	encodedSelfie := base64.StdEncoding.EncodeToString(selfieBytes)

	// Prepare response
	response := fiber.Map{
		"message": "ID card and selfie uploaded successfully",
		"id_card": encodedIDCard,
		"selfie":  encodedSelfie,
	}

	return c.JSON(response)
}

// GetUploadedDocuments retrieves both the ID card and selfie with ID for the user
// GetUploadedDocuments retrieves both the ID card and selfie with ID for the user
func GetUploadedDocuments(c *fiber.Ctx) error {
	db := middleware.DBConn
	claims, ok := c.Locals("user").(*users.Claims)
	if !ok {
		return unauthorizedResponse()
	}

	// Retrieve the user's verification record
	var verification users.UserVerification
	if err := db.Where("user_id = ?", claims.UserId).First(&verification).Error; err != nil {
		return notFoundResponse("No verification data found", err)
	}

	// Check if the ID card or selfie is available
	var response fiber.Map
	if len(verification.ValidId) == 0 && len(verification.Selfie) == 0 {
		return successResponse("No ID card or selfie with ID uploaded yet", nil)
	}

	// Prepare base64-encoded data for both ID card and selfie
	if len(verification.ValidId) > 0 {
		encodedIDCard := base64.StdEncoding.EncodeToString(verification.ValidId)
		response["id_card"] = encodedIDCard
	}

	if len(verification.Selfie) > 0 {
		encodedSelfie := base64.StdEncoding.EncodeToString(verification.Selfie)
		response["selfie"] = encodedSelfie
	}

	response["message"] = "ID card and selfie retrieved successfully"
	return c.JSON(response)
}

func unauthorizedResponse() error {
	return fiber.NewError(fiber.StatusUnauthorized, "Unauthorized")
}

func badRequestResponse(msg string, err error) error {
	return fiber.NewError(fiber.StatusBadRequest, msg+": "+err.Error())
}

func serverErrorResponse(msg string, err error) error {
	return fiber.NewError(fiber.StatusInternalServerError, msg+": "+err.Error())
}

func notFoundResponse(msg string, err error) error {
	return fiber.NewError(fiber.StatusNotFound, msg+": "+err.Error())
}

func successResponse(message string, data interface{}) error {
	return fiber.NewError(fiber.StatusOK, message)
}

func DeleteVerificationIfRejected(c *fiber.Ctx) error {
	db := middleware.DBConn
	claims, ok := c.Locals("user").(*users.Claims)
	if !ok {
		return unauthorizedResponse()
	}

	var verification users.UserVerification
	if err := db.Where("user_id = ?", claims.UserId).First(&verification).Error; err != nil {
		return notFoundResponse("Verification record not found", err)
	}

	// Only delete if status is "rejected"
	if verification.Status != "rejected" {
		return badRequestResponse("Verification status is not rejected", nil)
	}

	// ✅ SAFE DELETE with WHERE clause
	if err := db.Where("user_id = ?", claims.UserId).Delete(&users.UserVerification{}).Error; err != nil {
		return serverErrorResponse("Failed to delete verification record", err)
	}

	return c.JSON(fiber.Map{
		"message": "Rejected verification record deleted successfully",
	})
}
