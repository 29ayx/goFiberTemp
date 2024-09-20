package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"myapp/database"
	"myapp/routes"

	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	jwtware "github.com/gofiber/jwt/v3" // Import JWT middleware as jwtware
)

func main() {
	// Initialize the Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())  // Logs requests
	app.Use(recover.New()) // Recovers from panics
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow CORS from any URL
	}))

	// Connect to the PostgreSQL database
	database.Connect()

	// Public routes
	app.Post("/login", routes.Login)         // Public route for login
	app.Post("/register", routes.CreateUser) // Public route for user registration
	app.Get("/users/email", routes.GetUserByEmail)

	app.Get("/live", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// Protected routes using JWT middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("my_ultra_secure_secret"), // Ensure to use a secure key in production
	}))

	// Protected user routes
	app.Get("/users", routes.GetAllUsers)
	app.Get("/users/:id", routes.GetUser)
	app.Put("/users/:id", routes.UpdateUser)
	app.Delete("/users/:id", routes.DeleteUser)
	app.Post("/users/:id/cars", routes.AssignCarToUser)
	// Protected mother profile routes
	app.Put("/users/:id/role", routes.UpdateUserRole)

	app.Get("/mother-profile/:id", routes.GetMotherProfileByID)
	app.Get("/mother-profile", routes.GetMotherProfileByEmail)
	app.Post("/user/:id/mother-profile", routes.CreateMotherProfile)
	app.Put("/mother-profile/:id", routes.UpdateMotherProfile)
	app.Delete("/mother-profile/:id", routes.DeleteMotherProfile)

	app.Get("/forum/:id", routes.GetForumByID) // Get forum post by ID
	app.Get("/forum", routes.GetForumByEmail)
	app.Get("/forums", routes.GetAllForums)          // Get forum posts by email
	app.Post("/forum", routes.CreateForumPost)       // Create a new forum post
	app.Put("/forum/:id", routes.UpdateForumPost)    // Update a forum post by ID
	app.Put("/forum/:id/like", routes.LikeForumPost) // Increment likes for a forum post
	app.Put("/forum/:id/reply", routes.ReplyToForumPost)
	// Increment replies for a forum post

	// Protected mother journal routes
	app.Get("/mother-journal/:id", routes.GetMotherJournalByID)   // Get a journal entry by ID
	app.Get("/mother-journal", routes.GetMotherJournalByEmail)    // Get journal entries by email or date
	app.Post("/mother-journal", routes.CreateMotherJournal)       // Create a new journal entry
	app.Put("/mother-journal/:id", routes.UpdateMotherJournal)    // Update a journal entry by ID
	app.Delete("/mother-journal/:id", routes.DeleteMotherJournal) // Delete a journal entry by ID

	app.Post("/doctor-profile", routes.CreateDoctorProfile)
	app.Get("/doctor-profile", routes.GetDoctorProfileByEmail)

	// Start the server
	log.Fatal(app.Listen(":8000"))
}
