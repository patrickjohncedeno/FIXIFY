package routes

import (
	"fixify_backend/controller"
	"fixify_backend/controller/adminfeatures"
	"fixify_backend/controller/fetchings"
	"fixify_backend/controller/repairmanfeatures"
	"fixify_backend/controller/signuplogin"
	"fixify_backend/controller/userfeatures"
	"fixify_backend/websocketclient"
	"log"

	"fixify_backend/websocket"

	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
)

// Register all routes
func AppRoutes(app *fiber.App) {
	// SAMPLE ENDPOINT
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello Golang World!")
	})

	//
	token := app.Group("/token", signuplogin.JWTMiddleware) // Routes will be prefixed with /auth

	// -----------------------------
	//  AUTH ROUTES
	// -----------------------------

	app.Post("/signup/admin", signuplogin.AdminSignup)
	app.Post("/signup/user", signuplogin.UserSignup)
	app.Post("/signup/repairman", signuplogin.RepairmanSignup)

	app.Post("/login/admin", signuplogin.AdminLogin)
	app.Post("/login/user", signuplogin.UserLogin)

	app.Post("/login/password/:id", signuplogin.UserLogin) // Still might rename this one

	token.Post("/logout", signuplogin.JWTLogout)

	// -----------------------------
	//  VERIFICATION
	// -----------------------------

	// Email Verification
	app.Post("/verify/email/send", signuplogin.EmailVer)
	app.Post("/verify/email/code", signuplogin.EmailVerCode)
	// app.Post("/verify/email/resend", signuplogin.ResendVerificationCode)

	// Account Verification
	app.Patch("/verify/account/:id", adminfeatures.VerifyUser)

	// Valid ID
	app.Get("/validIDs", fetchings.FetchAllId)
	token.Get("/validID", fetchings.FetchValid)

	// -----------------------------
	//  USERS
	// -----------------------------

	// Fetch all users PUBLIC
	app.Get("/users", fetchings.FetchAllUsers)
	// Fetch single user
	app.Get("/users/:id", fetchings.FetchUser)
	// Fetch all clients (protected)
	app.Get("/clients", fetchings.FetchAllClient)
	// Fetch all repairmen (protected)
	app.Get("/repairmen", fetchings.FetchAllRepairmen)
	//review request
	token.Post("/requests/review/:id", userfeatures.ReviewRequest)

	//PRIVATE
	// Fetch all users
	token.Get("/users", fetchings.FetchAllUsers)
	// Fetch single user
	token.Get("/users/:id", fetchings.FetchUser)
	// Fetch all clients (protected)
	token.Get("/clients", fetchings.FetchAllClient)
	// Fetch all repairmen (protected)
	token.Get("/repairmen", fetchings.FetchAllRepairmen)
	// Fetch all admin
	token.Get("/admins", fetchings.FetchAllAdmin)

	// -----------------------------
	//  SERVICE REQUESTS
	// -----------------------------

	// Fetch all requests
	token.Get("/requests", fetchings.FetchAllRequest)
	app.Get("/requests", fetchings.FetchAllRequest)
	// Send request (protected)
	token.Post("/requests/:id", userfeatures.ServiceRequest)
	// Fetch by status
	token.Get("/requests/completed", fetchings.FetchCompletedRequest)
	token.Get("/requests/canceled", fetchings.FetchCanceledRequest)
	// Update request
	app.Patch("/requests/:id", repairmanfeatures.RequestUpdate)
	// -----------------------------
	// PERCENTAGE
	// -----------------------------

	// Percentage of requests
	// token.Get("/percentage/requests", adminfeatures.PercentageRequests)
	// Percentage of repairmen
	token.Get("/percentage/repairman", adminfeatures.PercentageRepairman)
	// Percentage of clients
	token.Get("/percentage/client", adminfeatures.PercentageClient)

	// -----------------------------
	// COUNTS
	// -----------------------------

	// Count of requests
	token.Get("/count/requests", fetchings.CountAllRequests)
	app.Get("/count/requests", fetchings.CountAllRequests)
	// Count of repairmen
	token.Get("/count/repairman", fetchings.CountAllRepairmen)
	app.Get("/count/repairman", fetchings.CountAllRepairmen)
	// Count of clients
	token.Get("/count/clients", fetchings.CountAllClients)
	app.Get("/count/clients", fetchings.CountAllClients)
	// Count request sent by user
	token.Get("/count/user/requests/:id", adminfeatures.CountUserRequests)
	// Count Admins
	token.Get("/count/admins", fetchings.CountAllAdmin)

	// -----------------------------
	//  ACCOUNT
	// -----------------------------

	// Update account (protected)
	token.Patch("/account/:id", userfeatures.UpdateAccount)
	token.Patch("/account/password/:id", userfeatures.UpdateAPassword)

	//Profile
	// Upload profile picture
	token.Post("/account/profile-picture", controller.UpdateProfilePicture)

	token.Delete("/user/verification/delete-if-rejected", userfeatures.DeleteVerificationIfRejected)

	// -----------------------------
	//  SERVICE CATEGORIES
	// -----------------------------

	//fetch all services
	token.Get("/services", fetchings.FetchServices)
	app.Get("/services", fetchings.FetchServices)
	//delete service
	app.Delete("/services/:id", fetchings.DeleteServiceCategory)
	//Disable service
	app.Patch("/service/disable/:id", fetchings.DisableServiceCategory)
	//Service to offer of repairman
	token.Post("/repairman/services", repairmanfeatures.UpdateRepairmanCategories)
	//admin can add service categories
	token.Post("/admin/services", adminfeatures.AddServiceCategory)
	//admin can update service categories
	app.Patch("/admin/services/:id", adminfeatures.UpdateService)

	// -----------------------------
	// NOTIFICATIONS
	// -----------------------------

	// Fetch all notifications
	app.Get("/notifications", fetchings.FetchAllUserNotifications)
	token.Get("/notifications/user", fetchings.ParamsNotification)

	// -----------------------------
	// REVIEWS
	// -----------------------------

	// Fetch all reviews
	token.Get("/reviews", fetchings.FetchAllReviews)
	token.Get("/ratings/:id", fetchings.ReviewRequest)
	token.Get("/repairmen/top", controller.TopRepairmen)

	// -----------------------------
	// MESSAGES
	// -----------------------------

	app.Get("/messagesclirep", fetchings.FetchClientRepairmanMessages)
	app.Get("/conversationsclirep", fetchings.FetchClientRepairmanConversations)

	app.Get("/conversations", fetchings.Conversations)

	// Add conversation
	token.Get("/conversations/available", adminfeatures.FetchAvailableAdminsForConversation)

	// -----------------------------
	// Upload
	// -----------------------------

	token.Patch("/users/:user_id/profile-picture", userfeatures.PatchProfilePicture)
	token.Post("/user/document/:id", userfeatures.UploadIDCardAndSelfie)
	token.Get("/user/profile-picture", userfeatures.GetProfilePicture)

	// -----------------------------
	// GCASH
	// -----------------------------
	token.Post("/gcash/pay", controller.InitiateXenditGCash)
	token.Post("/gcash/save", controller.SaveGCashInfo)

	// 🧠 Initialize the WebSocket Hub
	go websocket.HubInstance.Run()

	// 🧵 WebSocket Routes
	go websocketclient.HubInstance.Run()

	// 🧵 WebSocket Routes
	app.Use("/ws", websocket.WebSocketHandler)
	app.Get("/ws", fiberws.New(websocket.WebSocketUpgrade))

	app.Use("/ws/client", websocketclient.WebSocketHandler)
	app.Get("/ws/client", fiberws.New(websocketclient.WebSocketUpgrade))

	app.Post("/api/register-fcm-token", func(c *fiber.Ctx) error {
		type request struct {
			UserID   uint   `json:"user_id"`
			FCMToken string `json:"fcm_token"`
		}

		var req request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		if err := websocket.FCMInstance.RegisterToken(req.UserID, req.FCMToken); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to register token",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
		})
	})

	app.Post("/register_fcm_token", func(c *fiber.Ctx) error {
		if websocketclient.FCMInstance == nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "FCM service not initialized. Please contact support.",
			})
		}

		type request struct {
			UserID   uint   `json:"user_id"`
			FCMToken string `json:"fcm_token"`
		}

		var req request
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if err := websocketclient.FCMInstance.RegisterToken(req.UserID, req.FCMToken); err != nil {
			log.Printf("Failed to register FCM token: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to register token",
			})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "FCM token registered successfully",
		})
	})
	// Add this to your routes

}
