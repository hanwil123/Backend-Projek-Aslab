package Router

import (
	Authcontrollers "Backend-Projek-Aslab/Controllers/AuthControllers"

	"github.com/gofiber/fiber/v2"
	"github.com/nedpals/supabase-go"
	"github.com/supabase-community/gotrue-go"
)

func SetupUser(app *fiber.App, client *supabase.Client, clients gotrue.Client) {
	app.Get("/auth/google", Authcontrollers.GoogleAuthSignIn(client))
	app.Get("/auth/callback", Authcontrollers.GoogleAuthCallback(client))

}
