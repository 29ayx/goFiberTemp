package routes

import (
	"log"
	"myapp/database"
	"myapp/models"

	"github.com/gofiber/fiber/v2"
)

// FindUserByEmail checks if a user with the given email exists
func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	if result := database.DB.Where("email = ?", email).First(&user); result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}


// GetUserByEmail exposes the FindUserByEmail function via an API endpoint
func GetUserByEmail(c *fiber.Ctx) error {
	// Get the email from the query parameter
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email query parameter is required",
		})
	}

	// Find the user by email
	user, err := FindUserByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Return the user details
	return c.Status(fiber.StatusOK).JSON(user)
}

func GetUserPublicDetailByEmail(c *fiber.Ctx) error {
	// Get the email from the query parameter
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email query parameter is required",
		})
	}

	// Find the user by email
	user, err := FindUserByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Return the user details
	return c.Status(fiber.StatusOK).JSON(user)
}

func CreateUser(c *fiber.Ctx) error {
	user := new(models.User)

	// Parse user input from the request body
	if err := c.BodyParser(user); err != nil {
		log.Println("Body parsing failed:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Check if a user with the given email already exists
	existingUser, _ := FindUserByEmail(user.Email)
	if existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
	}

	// Save the new user to the database
	if result := database.DB.Create(&user); result.Error != nil {
		log.Println("User creation failed:", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "User could not be created"})
	}

	// After user creation, check if the role is 'pregnant' to create the profile
	if user.Role == "pregnant" {
		if err := createMotherProfile(user); err != nil {
			log.Println("Profile creation failed:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Profile creation failed"})
		}
	}

	if user.Role == "doctor" {
		if err := createDoctorProfile(user); err != nil {
			log.Println("Profile creation failed:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Profile creation failed"})
		}
	}



	return c.Status(fiber.StatusOK).JSON(user)
}


func createDoctorProfile(user *models.User) error {
	// Initialize a new profile struct
	doctor := &models.Doctor{
		Email:            user.Email,
		FirstName:	user.FirstName,
		LastName: user.LastName,
		Phone: user.Phone,
		AccStatus : "pending",
		ProfileOwnerType: "doctor",
	}

	// Save the profile to the database
	if result := database.DB.Create(doctor); result.Error != nil {
		log.Println("Profile creation failed:", result.Error)
		return result.Error
	}

	return nil
}
func createMotherProfile(user *models.User) error {
	// Initialize a new profile struct
	profile := &models.Profile{
		Email:            user.Email,
		PrefferedName:	user.FirstName,
		Phone: user.Phone,
		ProfileOwnerType: "mother",
	}

	// Save the profile to the database
	if result := database.DB.Create(profile); result.Error != nil {
		log.Println("Profile creation failed:", result.Error)
		return result.Error
	}

	return nil
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	database.DB.Preload("Cars").Find(&users)
	return c.Status(fiber.StatusOK).JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	log.Println("Fetching user with ID:", id) // Add this line for debugging

	if result := database.DB.First(&user, id); result.Error != nil {
		log.Println("Error fetching user:", result.Error) // Add this line for debugging
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	log.Println("User found:", user) // Add this line for debugging
	return c.Status(fiber.StatusOK).JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if result := database.DB.First(&user, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	database.DB.Save(&user)
	return c.Status(fiber.StatusOK).JSON(user)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	if result := database.DB.First(&user, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	database.DB.Delete(&user)
	return c.SendStatus(fiber.StatusNoContent)
}

func AssignCarToUser(c *fiber.Ctx) error {
	userId := c.Params("id")

	// Find the user by ID
	var user models.User
	if result := database.DB.First(&user, userId); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Parse the request body into the Car struct
	car := new(models.Car)
	if err := c.BodyParser(car); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Set the user ID to the car
	car.UserID = user.ID

	// Save the car to the database
	database.DB.Create(&car)
	return c.Status(fiber.StatusOK).JSON(car)
}

// UpdateUserRole assigns a role to the user
func UpdateUserRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Find the user by ID
	if result := database.DB.First(&user, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Parse the role from the request body
	role := struct {
		Role string `json:"role"`
	}{}

	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse role"})
	}

	// Update the user's role
	user.Role = role.Role
	database.DB.Save(&user)

	return c.Status(fiber.StatusOK).JSON(user)
}
