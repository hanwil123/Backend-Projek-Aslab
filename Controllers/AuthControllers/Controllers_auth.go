package Authcontrollers

import (
	databases "Backend-Projek-Aslab/Databases"
	"Backend-Projek-Aslab/Models/Auth"
	"fmt"
	"math/rand"

	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"
	"github.com/supabase-community/gotrue-go"
	"golang.org/x/crypto/bcrypt"
)

func Register(client *supabase.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		conteks := c.Context()
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot hash password"})
		}
		user, err := client.Auth.SignUp(conteks, supabase.UserCredentials{
			Email:    data["email"],
			Password: string(passwordHash),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Error creating user: %v", err)})
		}

		randomUserID := uint(rand.Intn(10000) + 1)
		userBaru := Auth.UserBaru{
			Userid:   uint64(randomUserID),
			Nama:     data["nama"],
			Email:    data["email"],
			Password: passwordHash,
		}

		databases.UDB.Create(&userBaru)
		return c.JSON(fiber.Map{
			"message": "Sign-up successful, verification email sent",
			"user":    user,
		})
	}
}

func Login(client gotrue.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse request body
		var data map[string]string
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		email, emailExists := data["email"]
		password, passwordExists := data["password"]

		// Validate input
		if !emailExists || !passwordExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email and password are required"})
		}

		// Sign in with email and password
		resp, err := client.SignInWithEmailPassword(email, password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
		}

		// Return successful response
		return c.JSON(fiber.Map{
			"access_token":  resp.AccessToken,
			"refresh_token": resp.RefreshToken,
			"user":          resp.User,
		})
	}
}
