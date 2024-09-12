package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"myapp/routes"
	"myapp/database"
	"myapp/models"
)

// Color and Emoji Constants for status display
const (
	greenCheck = "\u2705" // ✅
	redCross   = "\u274C" // ❌
	colorReset = "\033[0m"
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
)

// Setup Fiber app with routes
func setupApp() *fiber.App {
	app := fiber.New()

	// Connect to the database (ensure database connection is working)
	database.Connect()

	// Register the routes
	app.Post("/register", routes.CreateUser)
	app.Post("/login", routes.Login)
	app.Get("/users/email", routes.GetUserByEmail)  // Endpoint to find user by email
	app.Put("/users/:id", routes.UpdateUser)
	app.Delete("/users/:id", routes.DeleteUser)

	return app
}

// TestUserCRUD tests the full flow of API endpoints (Create, Login, Get by Email, Update, and Delete)
func TestUserCRUD(t *testing.T) {
	app := setupApp()

	// Step 1: Create User
	user := models.User{
		Name:     "Test User",
		Email:    "test_user@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	// Decode created user and get the userID
	var createdUser models.User
	json.NewDecoder(resp.Body).Decode(&createdUser)
	userID := createdUser.ID

	// Assert the user creation was successful
	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		t.Logf("%s %sUser Created | Test Passed OK%s\n", greenCheck, colorGreen, colorReset)
	} else {
		t.Fatalf("%s %sUser Created | Test Failed%s\n", redCross, colorRed, colorReset)
	}

	// Step 2: Login User and Get JWT Token
	loginData := map[string]string{
		"email":    "test_user@example.com",
		"password": "password123",
	}
	loginBody, _ := json.Marshal(loginData)
	req = httptest.NewRequest("POST", "/login", bytes.NewReader(loginBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req, -1)

	// Extract token from login response
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	token := result["token"].(string)

	// Assert that token is returned
	if assert.NotEmpty(t, token) {
		t.Logf("%s %sUser Login | Test Passed OK%s\n", greenCheck, colorGreen, colorReset)
	} else {
		t.Fatalf("%s %sUser Login | Test Failed%s\n", redCross, colorRed, colorReset)
	}

	// Step 3: Get User by Email (using token)
	req = httptest.NewRequest("GET", "/users/email?email=test_user@example.com", nil)
	req.Header.Set("Authorization", "Bearer "+token) // Use the token from login
	resp, _ = app.Test(req, -1)

	// Decode user from the response
	var foundUser models.User
	json.NewDecoder(resp.Body).Decode(&foundUser)
	userID = foundUser.ID // Extract userID for later operations

	// Assert that user details are fetched successfully
	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		t.Logf("%s %sGet User by Email | Test Passed OK%s\n", greenCheck, colorGreen, colorReset)
	} else {
		t.Fatalf("%s %sGet User by Email | Test Failed%s\n", redCross, colorRed, colorReset)
	}

	// Step 4: Update User Details (using token and userID)
	updatedUser := models.User{
		Name: "Updated User",
	}
	updateBody, _ := json.Marshal(updatedUser)
	req = httptest.NewRequest("PUT", fmt.Sprintf("/users/%d", userID), bytes.NewReader(updateBody)) // Use the correct userID
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, _ = app.Test(req, -1)

	// Assert the update was successful
	if assert.Equal(t, http.StatusOK, resp.StatusCode) {
		t.Logf("%s %sUpdate User | Test Passed OK%s\n", greenCheck, colorGreen, colorReset)
	} else {
		t.Fatalf("%s %sUpdate User | Test Failed%s\n", redCross, colorRed, colorReset)
	}

	// Step 5: Delete User (using token and userID)
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/users/%d", userID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, _ = app.Test(req, -1)

	// Assert the deletion was successful
	if assert.Equal(t, http.StatusNoContent, resp.StatusCode) {
		t.Logf("%s %sDelete User | Test Passed OK%s\n", greenCheck, colorGreen, colorReset)
	} else {
		t.Fatalf("%s %sDelete User | Test Failed%s\n", redCross, colorRed, colorReset)
	}
}
