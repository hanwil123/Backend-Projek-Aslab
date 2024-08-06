package Router

import (
	AwsS3controllers "Backend-Projek-Aslab/Controllers/AwsS3Controllers"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
)

func S3Router(app *fiber.App, s3Client *s3.S3) {
	app.Post("/created-bucket", AwsS3controllers.CreateBucket(s3Client))
	app.Post("/upload", AwsS3controllers.Upload(s3Client))
	app.Get("/getFiles", AwsS3controllers.GetFile(s3Client))
	app.Delete("/deleteFile", AwsS3controllers.DeleteFile(s3Client))
}
