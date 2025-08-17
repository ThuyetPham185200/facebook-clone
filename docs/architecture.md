facebook-clone/
│
├── services/
│   ├── gateway-api/                # API Gateway (entry point)
│   │   ├── cmd/
│   │   │   └── main.go
│   │   ├── internal/
│   │   │   ├── handler/
│   │   │   ├── service/
│   │   │   ├── repository/
│   │   │   └── ...
│   │   ├── pkg/
│   │   └── go.mod
│   │
│   ├── user-service/               # User Service
│   ├── follow-service/             # Follow Service
│   ├── post-service/               # Post Service
│   ├── feed-service/               # Feed Service
│   ├── like-service/               # Like Service
│   ├── comment-service/            # Comment Service
│   ├── notification-service/       # Notification Service
│   ├── content-delivery-service/   # Media/CDN handling
│   └── auth-service/               # Authentication & token
│
├── shared/                         # Share code dùng chung cho tất cả service
│   ├── configs/
│   ├── logger/
│   ├── middleware/
│   ├── models/
│   ├── utils/
│   ├── proto/                      # Nếu dùng gRPC
│   └── events/                     # Event schemas cho Kafka/NATS
│
├── deployments/                    # Docker-compose, k8s manifests
│   ├── docker-compose.yml
│   ├── k8s/
│   └── ...
│
├── docs/
│   ├── architecture-diagram.png
│   ├── api-design.md
│   └── ...
│
├── go.work                         # Nếu dùng Go workspace
└── README.md
