package routes

import (
	"github.com/gofiber/fiber/v2"
	"myapp/models"
	"myapp/database"
)

func CreateCar(c *fiber.Ctx) error {
	// Create a new car instance
	car := new(models.Car)

	// Parse the request body into the Car struct
	if err := c.BodyParser(car); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Save the car to the database
	database.DB.Create(&car)
	return c.Status(fiber.StatusOK).JSON(car)
}

func GetAllCars(c *fiber.Ctx) error {
	// Fetch all cars from the database
	var cars []models.Car
	database.DB.Find(&cars)
	return c.Status(fiber.StatusOK).JSON(cars)
}

func GetCar(c *fiber.Ctx) error {
	// Get the car ID from the URL params
	id := c.Params("id")

	// Find the car by ID
	var car models.Car
	if result := database.DB.First(&car, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Car not found"})
	}
	return c.Status(fiber.StatusOK).JSON(car)
}

func UpdateCar(c *fiber.Ctx) error {
	// Get the car ID from the URL params
	id := c.Params("id")

	// Find the car by ID
	var car models.Car
	if result := database.DB.First(&car, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Car not found"})
	}

	// Parse the updated data into the car object
	if err := c.BodyParser(&car); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Save the updated car to the database
	database.DB.Save(&car)
	return c.Status(fiber.StatusOK).JSON(car)
}

func DeleteCar(c *fiber.Ctx) error {
	// Get the car ID from the URL params
	id := c.Params("id")

	// Find the car by ID
	var car models.Car
	if result := database.DB.First(&car, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Car not found"})
	}

	// Delete the car from the database
	database.DB.Delete(&car)
	return c.SendStatus(fiber.StatusNoContent)
}
