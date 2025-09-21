package main

import "fmt"

func main() {
	s3Client := s3client.NewS3Client("http://localhost:9000", "us-east-1", "minioadmin", "minioadmin", "facebook-clone-media")

	go s3Client.ListenUploadNotifications("localhost:9092", "minio-events", "feedservice-group", func(key string) {
		// update medias table: set media.status = uploaded
		fmt.Println("Mark media as uploaded in DB:", key)
	})

}
