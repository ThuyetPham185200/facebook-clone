package s3client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/segmentio/kafka-go"
)

type S3Client struct {
	Client *s3.Client
	Bucket string
}

// NewS3Client with direct parameters instead of ENV
func NewS3Client(endpoint, region, accessKey, secretKey, bucket string) *S3Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithEndpointResolver(aws.EndpointResolverFunc(
			func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           endpoint,
					SigningRegion: region,
				}, nil
			})),
	)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	return &S3Client{
		Client: s3.NewFromConfig(cfg),
		Bucket: bucket,
	}
}

// Generate a pre-signed PUT URL
func (s *S3Client) GeneratePreSignedURL(objectKey string, expires time.Duration) string {
	ps := s3.NewPresignClient(s.Client)

	req, err := ps.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &objectKey,
	}, s3.WithPresignExpires(expires))
	if err != nil {
		log.Fatalf("failed to sign request: %v", err)
	}
	return req.URL
}

// Upload directly from Go (simulate user upload)
func (s *S3Client) UploadFile(objectKey, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	uploader := manager.NewUploader(s.Client)
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: &s.Bucket,
		Key:    &objectKey,
		Body:   file,
	})
	if err != nil {
		log.Fatalf("failed to upload: %v", err)
	}
	fmt.Println("File uploaded successfully:", objectKey)
}

// Consume notifications from Kafka for uploaded media
func (s *S3Client) ListenUploadNotifications(kafkaBroker, topic, groupID string, handle func(objectKey string)) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   topic,
		GroupID: groupID,
	})
	fmt.Println("Listening for S3 upload notifications...")

	for {
		msg, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("Error reading Kafka message:", err)
			continue
		}

		var event struct {
			EventName string `json:"EventName"`
			Key       string `json:"Key"`
			Bucket    string `json:"Bucket"`
		}
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("Failed to parse event:", err)
			continue
		}

		if event.EventName == "s3:ObjectCreated:Put" && event.Bucket == s.Bucket {
			fmt.Println("Upload completed for object:", event.Key)
			handle(event.Key) // e.g., update medias table in FeedService
		}
	}
}
