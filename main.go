package main

import (
	databases "Backend-Projek-Aslab/Databases"
	"Backend-Projek-Aslab/Router"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
	"github.com/supabase-community/gotrue-go"
)

func main() {
	errs := godotenv.Load(".env")
	if errs != nil {
		log.Fatalf("Error loading .env file")
	}
	backBlazeApiKey := os.Getenv("BACKBLAZE_API_KEY_ID")
	backBlazeAppKey := os.Getenv("BACKBLAZE_APPLICATION_KEY")
	backBlazeEndpoint := os.Getenv("BACKBLAZE_ENDPOINT_URL")
	backBlazeRegion := os.Getenv("BACKBLAZE_REGION")

	S3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(backBlazeApiKey, backBlazeAppKey, ""),
		Endpoint:         aws.String(backBlazeEndpoint),
		Region:           aws.String(backBlazeRegion),
		S3ForcePathStyle: aws.Bool(true),
	}
	s3Session, errr := session.NewSession(S3Config)
	if errr != nil {
		panic(errr)
	}
	s3Client := s3.New(s3Session)

	client := supabase.CreateClient(os.Getenv("SUPABASE_API_URL"), os.Getenv("SUPABASE_API_KEY"), true)
	clients := gotrue.New(os.Getenv("SUPABASE_API_URL"), os.Getenv("SUPABASE_API_KEY"))
	databases.ConnectUserAuth()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: false,
		AllowOrigins:     "*",
	}))
	Router.SetupUser(app, client, clients)
	Router.S3Router(app, s3Client)
	var port = envPortOr("3000")
	err := app.Listen(port)
	if err != nil {
		panic(err)
	}
}
func envPortOr(port string) string {
	// If `PORT` variable in environment exists, return it
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}
	// Otherwise, return the value of `port` variable from function argument
	return ":" + port
}
