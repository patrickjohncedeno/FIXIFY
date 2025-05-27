package main

import (
	"fixify_backend/middleware"
	"fixify_backend/routes"
	"fixify_backend/websocketclient" // Make sure this is the correct package
	"context"
	"encoding/json"
	"fmt"
	"log" // For error logging

	"firebase.google.com/go/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"google.golang.org/api/option"
)

func init() {
	fmt.Println("STARTING SERVER...")
	fmt.Println("INITIALIZE DB CONNECTION...")

	if middleware.ConnectDB() {
		fmt.Println("DB CONNECTION FAILED!")
		log.Fatal("Exiting due to DB connection failure.")
	} else {
		fmt.Println("DB CONNECTION SUCCESSFUL!")
	}

	// Initialize FCM using credentials from .env
	fmt.Println("INITIALIZING FCM...")
	db := middleware.GetDB()

	// Construct Firebase credentials from environment variables
	cred := &struct {
		Type                    string `json:"type"`
		ProjectID               string `json:"project_id"`
		PrivateKeyID            string `json:"private_key_id"`
		PrivateKey              string `json:"private_key"`
		ClientEmail             string `json:"client_email"`
		ClientID                string `json:"client_id"`
		AuthURI                 string `json:"auth_uri"`
		TokenURI                string `json:"token_uri"`
		AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
		ClientX509CertURL       string `json:"client_x509_cert_url"`
		UniverseDomain          string `json:"universe_domain"`
	}{
		Type:                    middleware.GetEnv("FIREBASE_TYPE"),
		ProjectID:               middleware.GetEnv("FIREBASE_PROJECT_ID"),
		PrivateKeyID:            middleware.GetEnv("FIREBASE_PRIVATE_KEY_ID"),
		PrivateKey:              middleware.GetEnv("FIREBASE_PRIVATE_KEY"),
		ClientEmail:             middleware.GetEnv("FIREBASE_CLIENT_EMAIL"),
		ClientID:                middleware.GetEnv("FIREBASE_CLIENT_ID"),
		AuthURI:                 middleware.GetEnv("FIREBASE_AUTH_URI"),
		TokenURI:                middleware.GetEnv("FIREBASE_TOKEN_URI"),
		AuthProviderX509CertURL: middleware.GetEnv("FIREBASE_AUTH_PROVIDER_X509_CERT_URL"),
		ClientX509CertURL:       middleware.GetEnv("FIREBASE_CLIENT_X509_CERT_URL"),
		UniverseDomain:          middleware.GetEnv("FIREBASE_UNIVERSE_DOMAIN"),
	}

	// Marshal the credentials to JSON to properly escape the private key
	credJSON, err := json.Marshal(cred)
	if err != nil {
		log.Fatalf("Failed to marshal Firebase credentials: %v", err)
	}

	// Initialize Firebase app with credentials
	ctx := context.Background()
	firebaseConfig := &firebase.Config{
		ProjectID: cred.ProjectID,
	}
	opt := option.WithCredentialsJSON(credJSON)
	app, err := firebase.NewApp(ctx, firebaseConfig, opt)
	if err != nil {
		log.Fatalf("Failed to initialize Firebase app: %v", err)
	}

	// Pass the Firebase app to InitializeFCM
	err = websocketclient.InitializeFCM(db, app)
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