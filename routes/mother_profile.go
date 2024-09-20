package routes

import (
	"log"
	"myapp/database"
	"myapp/models"

	"github.com/gofiber/fiber/v2"
)

// FindMotherProfileByID finds a mother profile by its ID
func FindMotherProfileByID(id uint) (*models.Profile, error) {
	var profile models.Profile
	if result := database.DB.Where("id = ? AND profile_owner_type = ?", id, "mother").First(&profile); result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}

// FindMotherProfileByEmail finds a mother profile by the associated user's email
func FindMotherProfileByEmail(email string) (*models.Profile, error) {
	var profile models.Profile
	if result := database.DB.
		Joins("JOIN users ON users.id = profiles.user_id").
		Where("users.email = ? AND profiles.profile_owner_type = ?", email, "pregnant").
		First(&profile); result.Error != nil {
		return nil, result.Error
	}
	return &profile, nil
}

// GetMotherProfileByID exposes the FindMotherProfileByID function via API
func GetMotherProfileByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid profile ID",
		})
	}

	profile, err := FindMotherProfileByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Mother profile not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

// GetMotherProfileByEmail exposes the FindMotherProfileByEmail function via API
func GetMotherProfileByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email query parameter is required",
		})
	}

	profile, err := FindMotherProfileByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Mother profile not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

// CreateMotherProfile creates a mother profile for the user
func CreateMotherProfile(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	// Find the user by ID
	if result := database.DB.First(&user, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Ensure the user has the 'mother' role
	if user.Role != "pregnant" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User does not have the mother role"})
	}

	// Parse the mother profile input from the request body
	profile := new(models.Profile)
	if err := c.BodyParser(profile); err != nil {
		log.Println("Body parsing failed:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Associate the profile with the user
	profile.UserID = user.ID
	profile.Email = user.Email
	profile.ProfileOwnerType = "mother"

	// Save the profile to the database
	if result := database.DB.Create(&profile); result.Error != nil {
		log.Println("Profile creation failed:", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Profile could not be created"})
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

// UpdateMotherProfile updates an existing mother profile
func UpdateMotherProfile(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid profile ID",
		})
	}

	var profile models.Profile
	if result := database.DB.First(&profile, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Profile not found"})
	}

	if err := c.BodyParser(&profile); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	database.DB.Save(&profile)
	return c.Status(fiber.StatusOK).JSON(profile)
}

// DeleteMotherProfile deletes a mother profile
func DeleteMotherProfile(c *fiber.Ctx) error {
	id := c.Params("id")
	var profile models.Profile
	if result := database.DB.First(&profile, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Profile not found"})
	}
	database.DB.Delete(&profile)
	return c.SendStatus(fiber.StatusNoContent)
}
