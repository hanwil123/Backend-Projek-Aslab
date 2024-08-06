package Authcontrollers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"
)

func GoogleAuthSignIn(client *supabase.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		oauthLogin, err := client.Auth.SignInWithProvider(supabase.ProviderSignInOptions{
			Provider:   "google",
			RedirectTo: "http://localhost:3000/auth/callback",
			FlowType:   supabase.PKCE,
			Scopes:     []string{"email"},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Redirect(oauthLogin.URL)
	}
}

func GoogleAuthCallback(client *supabase.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		conteks := c.Context()
		code := c.Query("code")
		if code == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing code in callback"})
		}
		codeVerifier := c.Query("code_verifier")
		if codeVerifier == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing code Verifier"})
		}

		user, err := client.Auth.ExchangeCode(conteks, supabase.ExchangeCodeOpts{
			AuthCode:     code,
			CodeVerifier: codeVerifier,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to exchange code for session"})
		}
		return c.JSON(fiber.Map{"message": "Sign-in successful, verification email sent", "user ": user.AccessToken})
	}
}
