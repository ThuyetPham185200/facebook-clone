package workerpocessor

import (
	"context"
	"feedservice/internal/infra/redisclient"
	"feedservice/internal/model"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type FanoutWorker struct {
	BaseWorkerProcessor
	redisclient       *redisclient.RedisClient
	newposteventqueue *chan model.NewPostEvent
}

func NewFanoutWorker(newposteventqueue_ *chan model.NewPostEvent, redisclient_ *redisclient.RedisClient) *FanoutWorker {
	s := &FanoutWorker{
		redisclient:       redisclient_,
		newposteventqueue: newposteventqueue_,
	}
	s.Init(s) // 🔑 rất quan trọng: gắn HttpServer vào BaseServerProcessor
	return s
}

func (s *FanoutWorker) RunningTask() error {
	if s.newposteventqueue == nil {
		return fmt.Errorf("newposteventqueue is nil")
	}

	for {
		select {
		case newPostEvent := <-*s.newposteventqueue:
			// Redis key per user
			key := fmt.Sprintf("user:%s:posts", newPostEvent.UserID)
			score := float64(time.Now().Unix()) // use timestamp for ordering

			// Use raw Redis client
			err := s.redisclient.GetClient().ZAdd(context.Background(), key, redis.Z{
				Score:  score,
				Member: newPostEvent.PostID,
			}).Err()
			if err != nil {
				log.Printf("[FanoutWorker] failed to add post to redis: %v", err)
				continue
			}

			log.Printf("[FanoutWorker] added post %s to user %s feed", newPostEvent.PostID, newPostEvent.UserID)

		case <-time.After(1 * time.Second):
			// idle wait to prevent busy loop
		}
	}
}
