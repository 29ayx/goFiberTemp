package routes

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"log"
)

var jwtSecret = []byte("my_ultra_secure_secret")

// Login function to authenticate users and return JWT
func Login(c *fiber.Ctx) error {
    var input struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }

    // Parse the login input
    if err := c.BodyParser(&input); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    // Find the user by email
    user, err := FindUserByEmail(input.Email)
    if err != nil || user == nil {
        log.Println("User not found:", input.Email)
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
    }

    // Check if the password matches (no hashing involved, simple string comparison)
    if user.Password != input.Password {
        log.Println("Password mismatch")
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
    }

    // Generate JWT token
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),  // Token expiration in 72 hours
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString(jwtSecret)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
    }

    // Return the token
    return c.JSON(fiber.Map{"token": tokenString})
}
