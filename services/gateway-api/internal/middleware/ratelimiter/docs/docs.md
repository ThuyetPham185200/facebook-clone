rate-limiter/
│
├── ratelimiter.go                       # Entry point
│   ratelimiter-test.go                  # To test rate-limiter module dependently   
│
│── limiter/
│      ├── rate_limiter.go        # interface RateLimiter + wrapper dùng Strategy Algorithm
│      ├── ip_limiter.go          # IPRateLimiter
│      ├── feature_limiter.go     # FeatureRateLimiter
│      ├── token_bucket.go        # TokenBucket state struct
│      ├── leaky_bucket.go        # (option) LeakyBucket state struct
│      └── sliding_window.go      # (option) Sliding Window state struct
│   
│── algorithm/                 # Strategy cho thuật toán rate limiting
│      ├── algorithm.go           # interface RateLimitAlgorithm
│      ├── token_bucket_algo.go   # implement TokenBucketAlgorithm
│      ├── leaky_bucket_algo.go   # implement LeakyBucketAlgorithm
│      └── sliding_window_algo.go # implement SlidingWindowAlgorithm
│   
│── repository/                # Abstraction layer cho Redis
│      ├── redis_repository.go    # interface RedisStateRepository
│      └── redis_repository_impl.go # implement RedisStateRepository
│   
│── factory/
│   └── limiter_factory.go     # Factory tạo limiter + inject algorithm + repository
│   
│── storage/
│      ├── redis_client.go        # Singleton Redis client
│      └── postgres_client.go     # Singleton Postgres client
│   
│── config/
│      └── loader.go              # Load rules từ DB, cập nhật cache
│   
│── utils/
│       ── logger.go              # Logging, metrics
│
└── go.mod
