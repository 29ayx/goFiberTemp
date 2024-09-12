package main

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3"  // Import JWT middleware as jwtware
	"myapp/database"
	"myapp/routes"
)

func main() {
	// Initialize the Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())   // Logs requests
	app.Use(recover.New())  // Recovers from panics

	// Connect to the PostgreSQL database
	database.Connect()

	// Public routes
	app.Post("/login", routes.Login)   // Public route for login
	app.Post("/register", routes.CreateUser) // Public route for user registration
	app.Get("/users/email", routes.GetUserByEmail)

	// Protected routes using JWT middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("my_ultra_secure_secret"),  // Ensure to use a secure key in production
	}))
	

	// Protected user routes
	app.Get("/users", routes.GetAllUsers)
	app.Get("/users/:id", routes.GetUser)
	app.Put("/users/:id", routes.UpdateUser)
	app.Delete("/users/:id", routes.DeleteUser)
	app.Post("/users/:id/cars", routes.AssignCarToUser)

	// Start the server
	log.Fatal(app.Listen(":3000"))
}
