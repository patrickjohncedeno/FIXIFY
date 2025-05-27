package main

import (
	"fixify_backend/middleware"
	"fixify_backend/routes"
	"fixify_backend/websocketclient"
	"fmt"
	"log" // For error logging

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	fmt.Println("STARTING SERVER...")
	fmt.Println("INITIALIZE DB CONNECTION...")
	if middleware.ConnectDB() {
		fmt.Println("DB CONNECTION FAILED!")
		// Optionally terminate the program if DB connection fails
		log.Fatal("Exiting due to DB connection failure.")
	} else {
		fmt.Println("DB CONNECTION SUCCESSFUL!")
	}
	// ✅ Initialize FCM here
	fmt.Println("INITIALIZING FCM...")
	db := middleware.GetDB()
	
	// Initialize FCM only once with the correct credentials file
	err := websocketclient.InitializeFCM(db, "firebase_credential_users.json") // Use the correct credentials file
	if err != nil {
		log.Fatalf("Failed to initialize FCM: %v", err)
	}

	fmt.Println("FCM INITIALIZED SUCCESSFULLY!")
}

func main() {
	app := fiber.New(fiber.Config{
		AppName: middleware.GetEnv("PROJ_NAME"),
	})

	// CORS CONFIG (before setting routes)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Use "*" for development, but restrict for production
		AllowMethods: "GET, POST, PUT, DELETE, PATCH",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Sample Endpoint (optional, if no favicon is used)
	app.Get("/favicon.ico", func(c *fiber.Ctx) error {
		return c.SendStatus(204) // No Content
	})

	// Set up the application routes
	routes.AppRoutes(app)

	// LOGGER middleware - It's better to use it before setting up routes for consistency
	app.Use(logger.New())

	// Start Server
	err := app.Listen(fmt.Sprintf(":%s", middleware.GetEnv("PROJ_PORT")))
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
