package main

import (
	"feedservice/internal/infra/s3client"
	"fmt"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// GenerateRandomObjectKey generates a random object key with the same file extension
func GenerateRandomObjectKey(originalFilename string) string {
	ext := filepath.Ext(originalFilename) // e.g. ".png"
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}
func main() {
	// User-provided config
	endpoint := "http://localhost:9100"
	bucket := "facebook-clone-media"
	region := "us-east-1"
	accessKey := "minioadmin"
	secretKey := "minioadmin"

	// Init client
	s3client := s3client.NewS3Client(endpoint, region, accessKey, secretKey, bucket)

	randomKey := GenerateRandomObjectKey("bell.png")
	url := s3client.GeneratePreSignedURL(randomKey, 5*time.Minute)
	fmt.Println("Presigned URL:", url)

	// Example: upload directly (optional)
	//s3client.UploadFile("direct-upload.png", "../assets/bell.png")
}
