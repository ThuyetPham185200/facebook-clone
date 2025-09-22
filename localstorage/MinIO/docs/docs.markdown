docker-compose up -d
docker logs minio
## MinIO web
http://localhost:9005

## curl to upload image
curl -i -X PUT -T bell.png "http://localhost:9100/facebook-clone-media/8c2e7df0-9034-4246-a693-8095b25a1211.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=minioadmin%2F20250922%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20250922T031012Z&X-Amz-Expires=300&X-Amz-SignedHeaders=host&x-id=PutObject&X-Amz-Signature=c7304a9e5fd4fe3c3ba01abee6cc5d72fc4fb842e2d4b11cc6d94a3a63f313e6"
