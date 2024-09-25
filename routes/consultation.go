package routes

import (
	"log"
	"myapp/database"
	"myapp/models"
	"github.com/gofiber/fiber/v2"
)

func CreateAppointment(c *fiber.Ctx) error {
	//Creates appointment for the user
	var consultation models.AppointmentForm
	if err := c.BodyParser(&consultation); err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	if result := database.DB.Create(&consultation); result.Error != nil {
		log.Println(result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create consultation",
		})
	}
	return c.JSON(consultation)
}

func GetAppointmentById(c *fiber.Ctx) error {
	//returns the appointment based on the id also provide the patient profile
	id := c.Params("id")
	var consultation models.AppointmentForm
	if result := database.DB.First(&consultation, id); result.Error != nil {
		log.Println(result.Error)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Consultation not found",
		})
	}

	// Call FindMotherProfileByEmail using the PatientEmail
	motherProfile, err := FindMotherProfileByEmail(consultation.PatientEmail)
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Mother profile not found",
		})
	}

	// Return the consultation data and the mother profile
	return c.JSON(fiber.Map{
		"appointment": consultation,
		"mother_profile": motherProfile,
	})
}


func GetAppointmentsBySpecialist(c *fiber.Ctx) error {
	//Shows all the apppoinments based on the specialist
	specialist := c.Params("specialist")

	// Check if the 'specialist' parameter is provided
	if specialist == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Specialist parameter is missing",
		})
	}

	var appointments []models.AppointmentForm

	if result := database.DB.Where("specialist = ?", specialist).Find(&appointments); result.Error != nil {
		log.Println(result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch appointments",
		})
	}

	if len(appointments) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "No appointments found for this specialist",
		})
	}

	return c.JSON(appointments)
}
