package routes

import (
	"log"
	"myapp/database"
	"myapp/models"
	"github.com/gofiber/fiber/v2"
)

// Create admin post
func CreateAdminPost(c *fiber.Ctx) error {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Thumbnail string `json:"thumbnail"`
		Description string `json:"description"`
		PostType string `json:"post_type"`
	}

	// Parse the request body
	if err := c.BodyParser(&input); err != nil {
		log.Println("Body parsing failed:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Create the admin post
	post := models.AdminPost{
		Title:   input.Title,
		Content: input.Content,
		Thumbnail: input.Thumbnail,
		Description: input.Description,
		PostType: input.PostType,
		

	}
	database.DB.Create(&post)

	return c.Status(fiber.StatusCreated).JSON(post)
}

func GetAdminPosts(c *fiber.Ctx) error {
	var posts []models.AdminPost
	database.DB.Find(&posts)
	return c.JSON(posts)
}