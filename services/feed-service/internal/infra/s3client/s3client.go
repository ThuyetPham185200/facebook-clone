package s3client

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client *s3.Client
	Bucket string
}

// NewS3Client with direct parameters instead of ENV
func NewS3Client(endpoint, region, accessKey, secretKey, bucket string) *S3Client {
	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           endpoint,
			SigningRegion: region,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolver(customResolver),
	)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // ✅ Important for MinIO
	})

	return &S3Client{
		Client: client,
		Bucket: bucket,
	}
}

// Generate a pre-signed PUT URL (safe to copy-paste into curl)
func (s *S3Client) GeneratePreSignedURL(objectKey string, expires time.Duration) string {
	ps := s3.NewPresignClient(s.Client)

	req, err := ps.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &objectKey,
	}, s3.WithPresignExpires(expires))
	if err != nil {
		log.Fatalf("failed to sign request: %v", err)
	}
	// ✅ Wrap with quotes so it's shell-safe
	return req.URL
}

// GeneratePreSignedGetURL generates a pre-signed URL for GET requests
func (s *S3Client) GeneratePreSignedGetURL(objectKey string, expires time.Duration) string {
	ps := s3.NewPresignClient(s.Client)

	req, err := ps.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &s.Bucket,
		Key:    &objectKey,
	}, s3.WithPresignExpires(expires))
	if err != nil {
		log.Fatalf("failed to sign GET request: %v", err)
	}
	return req.URL
}

// Upload directly from Go (simulate user upload)
func UploadWithPresignedURL(url string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	req, err := http.NewRequest("PUT", url, file)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set Content-Type (optional, e.g., image/png)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to upload: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed, status: %s, body: %s", resp.Status, string(body))
	}

	fmt.Println("Upload successful!")
	return nil
}
