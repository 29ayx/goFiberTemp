package routes

import (
    "github.com/gofiber/fiber/v2"
    "myapp/models"
    "myapp/database"
    "log"
)

// FindMotherJournalByID finds a journal entry by its ID
func FindMotherJournalByID(id uint) (*models.DailyEntry, error) {
    var journal models.DailyEntry
    if result := database.DB.First(&journal, id); result.Error != nil {
        return nil, result.Error
    }
    return &journal, nil
}

// FindMotherJournalByEmail finds journal entries by the associated user's email
func FindMotherJournalByEmail(email string) ([]models.DailyEntry, error) {
    var journals []models.DailyEntry
    if result := database.DB.Where("user_email = ?", email).Find(&journals); result.Error != nil {
        return nil, result.Error
    }
    return journals, nil
}

// FindMotherJournalByDate finds journal entries by date
func FindMotherJournalByDate(entryDate string) ([]models.DailyEntry, error) {
    var journals []models.DailyEntry
    if result := database.DB.Where("entry_date = ?", entryDate).Find(&journals); result.Error != nil {
        return nil, result.Error
    }
    return journals, nil
}

// GetMotherJournalByID exposes the FindMotherJournalByID function via API
func GetMotherJournalByID(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid journal entry ID",
        })
    }

    journal, err := FindMotherJournalByID(uint(id))
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Journal entry not found",
        })
    }

    return c.Status(fiber.StatusOK).JSON(journal)
}

// GetMotherJournalByEmail exposes the FindMotherJournalByEmail function via API
func GetMotherJournalByEmail(c *fiber.Ctx) error {
    email := c.Query("email")
    if email == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Email query parameter is required",
        })
    }

    journals, err := FindMotherJournalByEmail(email)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Journal entries not found",
        })
    }

    return c.Status(fiber.StatusOK).JSON(journals)
}

// GetMotherJournalByDate exposes the FindMotherJournalByDate function via API
func GetMotherJournalByDate(c *fiber.Ctx) error {
    date := c.Query("date")
    if date == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Date query parameter is required",
        })
    }

    journals, err := FindMotherJournalByDate(date)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
            "error": "Journal entries not found",
        })
    }

    return c.Status(fiber.StatusOK).JSON(journals)
}

// CreateMotherJournal creates a new journal entry
func CreateMotherJournal(c *fiber.Ctx) error {
    journal := new(models.DailyEntry)

    // Parse journal input from the request body
    if err := c.BodyParser(journal); err != nil {
        log.Println("Body parsing failed:", err)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    // Save the new journal entry to the database
    if result := database.DB.Create(&journal); result.Error != nil {
        log.Println("Journal creation failed:", result.Error)
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Journal entry could not be created"})
    }

    return c.Status(fiber.StatusOK).JSON(journal)
}

// UpdateMotherJournal updates an existing journal entry by ID
func UpdateMotherJournal(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid journal entry ID",
        })
    }

    var journal models.DailyEntry
    if result := database.DB.First(&journal, id); result.Error != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Journal entry not found"})
    }

    // Parse updated fields from request body
    if err := c.BodyParser(&journal); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
    }

    // Save updated journal entry to the database
    if result := database.DB.Save(&journal); result.Error != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Journal entry could not be updated"})
    }

    return c.Status(fiber.StatusOK).JSON(journal)
}

// DeleteMotherJournal deletes a journal entry by ID
func DeleteMotherJournal(c *fiber.Ctx) error {
    id := c.Params("id")
    var journal models.DailyEntry
    if result := database.DB.First(&journal, id); result.Error != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Journal entry not found"})
    }
    database.DB.Delete(&journal)
    return c.SendStatus(fiber.StatusNoContent)
}
