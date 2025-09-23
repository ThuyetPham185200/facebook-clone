package workerpocessor

import (
	"context"
	"feedservice/internal/core/followserviceclient"
	"feedservice/internal/infra/redisclient"
	"feedservice/internal/model"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type FanoutWorker struct {
	BaseWorkerProcessor
	redisclient         *redisclient.RedisClient
	followserviceclient *followserviceclient.FollowServiceClient
	newposteventqueue   *chan model.NewPostEvent
}

func NewFanoutWorker(newposteventqueue_ *chan model.NewPostEvent, redisclient_ *redisclient.RedisClient,
	followserviceclient_ *followserviceclient.FollowServiceClient) *FanoutWorker {
	s := &FanoutWorker{
		redisclient:         redisclient_,
		followserviceclient: followserviceclient_,
		newposteventqueue:   newposteventqueue_,
	}
	s.Init(s)
	return s
}

func (s *FanoutWorker) RunningTask() error {
	if s.newposteventqueue == nil {
		return fmt.Errorf("newposteventqueue is nil")
	}

	for {
		select {
		case newPostEvent := <-*s.newposteventqueue:
			ctx := context.Background()
			score := float64(time.Now().Unix())

			// 1️⃣ Cache mapping: post_id -> user_id
			err := s.redisclient.GetClient().HSet(ctx,
				"post_authors",
				newPostEvent.PostID,
				newPostEvent.UserID,
			).Err()
			if err != nil {
				log.Printf("[FanoutWorker] failed to cache author for post %s: %v", newPostEvent.PostID, err)
				continue
			}

			// 2️⃣ Add to author’s own posts
			authorPostsKey := fmt.Sprintf("user:%s:posts", newPostEvent.UserID)
			if err := s.redisclient.GetClient().ZAdd(ctx, authorPostsKey, redis.Z{
				Score:  score,
				Member: newPostEvent.PostID,
			}).Err(); err != nil {
				log.Printf("[FanoutWorker] failed to add post %s to author %s posts: %v",
					newPostEvent.PostID, newPostEvent.UserID, err)
				continue
			}

			// 3️⃣ Fetch followers from FollowService
			followers, err := s.followserviceclient.GetFollowers(newPostEvent.UserID)
			if err != nil {
				log.Printf("[FanoutWorker] failed to fetch followers for user=%s: %v", newPostEvent.UserID, err)
				continue
			}

			// 4️⃣ Fanout to followers’ feeds
			for _, followerID := range followers {
				feedKey := fmt.Sprintf("user:%s:feed", followerID)
				if err := s.redisclient.GetClient().ZAdd(ctx, feedKey, redis.Z{
					Score:  score,
					Member: newPostEvent.PostID,
				}).Err(); err != nil {
					log.Printf("[FanoutWorker] failed to add post %s to feed of user %s: %v",
						newPostEvent.PostID, followerID, err)
					continue
				}
				log.Printf("[FanoutWorker] added post %s to feed of user %s",
					newPostEvent.PostID, followerID)
			}

		case <-time.After(1 * time.Second):
			// idle wait
		}
	}
}
