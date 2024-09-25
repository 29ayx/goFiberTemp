package routes

import (
	"log"
	"myapp/database"
	"myapp/models"

	"github.com/gofiber/fiber/v2"
)

// FindForumByID finds a forum post by its ID
func FindForumByID(id uint) (*models.ForumPost, error) {
	var forumPost models.ForumPost
	if result := database.DB.First(&forumPost, id); result.Error != nil {
		return nil, result.Error
	}
	return &forumPost, nil
}

// FindForumByEmail finds forum posts by the associated user's email
func FindForumByEmail(email string) ([]models.ForumPost, error) {
	var forumPosts []models.ForumPost
	if result := database.DB.Where("email = ?", email).Find(&forumPosts); result.Error != nil {
		return nil, result.Error
	}
	return forumPosts, nil
}

func FindForumByCategory(category string) ([]models.ForumPost, error) {
	var forumPosts []models.ForumPost
	if result := database.DB.Where("category= ?", category).Find(&forumPosts); result.Error != nil {
		return nil, result.Error
	}

	for i, post := range forumPosts {
		profile, err := FindMotherProfileByEmail(post.Email)
		if err != nil {
			return nil, err
		}
		forumPosts[i].Name = profile.PrefferedName 
	}

	return forumPosts, nil
}


// func FindForumByCategory(category string) ([]models.ForumPost, error) {
// 	var forumPosts []models.ForumPost
// 	if result := database.DB.Where("category= ?", category).Find(&forumPosts); result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return forumPosts, nil
// }

// GetForumByID exposes the FindForumByID function via API
func GetForumByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid forum post ID",
		})
	}

	forumPost, err := FindForumByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Forum post not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(forumPost)
}

// GetAllForums retrieves all forum posts
func GetAllForums(c *fiber.Ctx) error {
	var forumPosts []models.ForumPost

	if result := database.DB.Find(&forumPosts); result.Error != nil {
		log.Println("Error fetching forum posts:", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve forum posts",
		})
	}

	return c.Status(fiber.StatusOK).JSON(forumPosts)
}

// GetForumByEmail exposes the FindForumByEmail function via API
func GetForumByEmail(c *fiber.Ctx) error {
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email query parameter is required",
		})
	}

	forumPosts, err := FindForumByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Forum posts not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(forumPosts)
}

func GetForumByCategory(c *fiber.Ctx) error {
	category := c.Query("category")
	if category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Category query parameter is required",
		})
	}

	forumPosts, err := FindForumByCategory(category)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Forum posts not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(forumPosts)
}

// CreateForumPost creates a new forum post
func CreateForumPost(c *fiber.Ctx) error {
	forumPost := new(models.ForumPost)

	// Parse forum post input from the request body
	if err := c.BodyParser(forumPost); err != nil {
		log.Println("Body parsing failed:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Save the new forum post to the database
	if result := database.DB.Create(&forumPost); result.Error != nil {
		log.Println("Forum post creation failed:", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Forum post could not be created"})
	}

	return c.Status(fiber.StatusOK).JSON(forumPost)
}

// LikeForumPost increments the number of likes for a forum post
func LikeForumPost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid forum post ID"})
	}

	var forumPost models.ForumPost
	if result := database.DB.First(&forumPost, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Forum post not found"})
	}

	forumPost.Likes += 1
	database.DB.Save(&forumPost)

	return c.Status(fiber.StatusOK).JSON(forumPost)
}

// ReplyToForumPost increments the number of replies for a forum post
func ReplyToForumPost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid forum post ID"})
	}

	var forumPost models.ForumPost
	if result := database.DB.First(&forumPost, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Forum post not found"})
	}

	forumPost.Replies += 1
	database.DB.Save(&forumPost)

	return c.Status(fiber.StatusOK).JSON(forumPost)
}

// UpdateForumPost updates an existing forum post by ID
func UpdateForumPost(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid forum post ID"})
	}

	var forumPost models.ForumPost
	if result := database.DB.First(&forumPost, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Forum post not found"})
	}

	// Parse updated fields from request body
	if err := c.BodyParser(&forumPost); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Save updated forum post to the database
	if result := database.DB.Save(&forumPost); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Forum post could not be updated"})
	}

	return c.Status(fiber.StatusOK).JSON(forumPost)
}
