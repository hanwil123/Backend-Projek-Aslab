package AwsS3controllers

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gofiber/fiber/v2"
)

func CreateBucket(s3Client *s3.S3) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bucketName := c.Query("bucket-Tanamin")

		cparams := &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		}
		_, err := s3Client.CreateBucket(cparams)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON("Successfully created bucket %s\n", bucketName)
	}
}
func Upload(s3Client *s3.S3) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bucketName := c.Query("Tanamin-Bucket")
		key := c.Query("key-Tanamin")
		content := c.Query("content")

		_, err := s3Client.PutObject(&s3.PutObjectInput{
			Body:   strings.NewReader(content),
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
		})
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed to upload object %s/%s, %s\n", bucketName, key, err.Error()))
		}
		return c.JSON("Successfully uploaded File %s\n", key)
	}
}

func GetFile(s3Client *s3.S3) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bucketName := c.Query("Tanamin-Bucket")

		// List objects in the bucket
		resp, err := s3Client.ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			return c.Status(500).SendString(fmt.Sprintf("Failed to list objects in bucket %s: %s\n", bucketName, err.Error()))
		}
		// Prepare a slice to store file information
		var files []map[string]interface{}

		// Iterate through the objects and gather information
		for _, item := range resp.Contents {
			file := map[string]interface{}{
				"Key":          *item.Key,
				"LastModified": *item.LastModified,
				"Size":         *item.Size,
				"StorageClass": *item.StorageClass,
			}
			files = append(files, file)
		}

		// Return the list of files
		return c.JSON(fiber.Map{
			"bucket": bucketName,
			"files":  files,
		})
	}
}

func DeleteFile(s3Client *s3.S3) fiber.Handler {
	return func(c *fiber.Ctx) error {
		bucketName := c.Query("Tanamin-Bucket")
		key := c.Query("key-Tanamin")

		if bucketName == "" || key == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Both Tanamin-Bucket and key-Tanamin query parameters are required",
			})
		}

		// Delete all versions if versioning is enabled
		listVersionsInput := &s3.ListObjectVersionsInput{
			Bucket: aws.String(bucketName),
			Prefix: aws.String(key),
		}

		err := s3Client.ListObjectVersionsPages(listVersionsInput,
			func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
				for _, version := range page.Versions {
					_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
						Bucket:    aws.String(bucketName),
						Key:       aws.String(key),
						VersionId: version.VersionId,
					})
					if err != nil {
						fmt.Printf("Error deleting version %s: %v\n", *version.VersionId, err)
					}
				}
				return !lastPage
			})

		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to delete object %s: %s", key, err.Error()),
			})
		}

		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Successfully deleted all versions of file %s from bucket %s", key, bucketName),
		})
	}
}
