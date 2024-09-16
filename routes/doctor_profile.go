package routes

import (
	"github.com/gofiber/fiber/v2"
	"myapp/models"
	"myapp/database"
	"log"
)

func CreateDoctorProfile(c *fiber.Ctx) error {

	id := c.Params("id")
	var user models.User

	if result := database.DB.First(&user, id); result.Error !=nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error":"User not found"})
	}

	// if user.Role != "doctor" {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"User dont have doctor role"})
	// }

	profile := new(models.Doctor)
	if err:= c.BodyParser(profile); err != nil {
		log.Println("Body parsing failed:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"Invalid Input"})
	}

	profile.UserID = user.ID

	profile.ProfileOwnerType = "doctor"

	if result := database.DB.Create(&profile); result.Error != nil {
		log.Println("Profile creation failed:", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Profile could not be created"})
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}


func FindDoctorByEmail(email string)(*models.Doctor, error){
	var profile models.Doctor
	if result := database.DB.
	Joins("JOIN users ON users.id = doctors.user_id").
	Where("doctors.email = ? ", email).
	First(&profile); result.Error != nil{
		return nil, result.Error
	}
	return &profile, nil
}

func GetDoctorProfileByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == ""{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Email query parameter is required",
		})
	}

	profile, err := FindDoctorByEmail(email)
	if err != nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":"Doctor account not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}