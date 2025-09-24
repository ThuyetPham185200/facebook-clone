package scrollingfeedmanager

import (
	"context"
	"feedservice/internal/core/userserviceclient"
	"feedservice/internal/infra/redisclient"
	"feedservice/internal/infra/store"
	"fmt"
	"log"
	"time"
)

type SrollingFeedManager struct {
	MediaStore        *store.MediaStore
	PostStore         *store.PostStore
	PostMediaStore    *store.PostMediaStore
	userserviceclient *userserviceclient.UserService
	redisclient       *redisclient.RedisClient
}

func NewSrollingFeedManager(
	MediaStore_ *store.MediaStore,
	PostStore_ *store.PostStore,
	PostMediaStore_ *store.PostMediaStore,
	userserviceclient_ *userserviceclient.UserService,
	redisclient_ *redisclient.RedisClient) *SrollingFeedManager {
	return &SrollingFeedManager{
		MediaStore:        MediaStore_,
		PostStore:         PostStore_,
		PostMediaStore:    PostMediaStore_,
		userserviceclient: userserviceclient_,
		redisclient:       redisclient_,
	}
}

type FeedItem struct {
	PostID         string    `json:"post_id"`
	Author         UserBrief `json:"author"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
	Media          []Media   `json:"media"`
	Stats          PostStats `json:"stats"`
	ViewerReaction *string   `json:"viewer_reaction,omitempty"`
}

type UserBrief struct {
	UserID   string `json:"user_id"`
	Username string `json:"name"`
	Avatar   string `json:"avatar_url,omitempty"`
}

type Media struct {
	MediaID string `json:"media_id"`
	Type    string `json:"type"`
	URL     string `json:"url"`
}

type PostStats struct {
	Likes    int `json:"likes"`
	Comments int `json:"comments"`
}

type FeedResponse struct {
	Feed       []FeedItem `json:"feed"`
	NextOffset int64      `json:"next_offset"`
}

func (s *SrollingFeedManager) ScrollingFeed(userID string, offset, limit int64) (FeedResponse, error) {
	ctx := context.Background()
	feedKey := fmt.Sprintf("user:%s:feed", userID)

	// Get posts (most recent first)
	postIDs, err := s.redisclient.GetClient().
		ZRevRange(ctx, feedKey, offset, offset+limit-1).
		Result()
	if err != nil {
		return FeedResponse{}, fmt.Errorf("[ScrollingFeed] failed to fetch feed for user %s: %w", userID, err)
	}

	feed := make([]FeedItem, 0, len(postIDs))

	for _, postID := range postIDs {
		// 1. Fetch post content
		post, err := s.PostStore.GetPostByID(postID)
		if err != nil {
			log.Printf("[ScrollingFeed] failed to fetch post %s: %v", postID, err)
			continue
		}
		if post == nil {
			continue
		}

		// 2. Fetch media linked to this post
		mediaIDs, err := s.PostMediaStore.GetMediaByPostID(postID)
		if err != nil {
			log.Printf("[ScrollingFeed] failed to fetch media for post %s: %v", postID, err)
			continue
		}
		mediaList := []Media{}
		for _, mid := range mediaIDs {
			m, err := s.MediaStore.GetMediaByID(mid)
			if err == nil {
				mediaList = append(mediaList, Media{
					MediaID: m.MediaID,
					Type:    m.MediaType,
					URL:     m.URL,
				})
			}
		}

		// 3. Fetch author profile
		author, err := s.userserviceclient.GetUserProfile(post.UserID)
		if err != nil {
			log.Printf("[ScrollingFeed] failed to fetch author for user %s: %v", post.UserID, err)
		}

		authorBrief := UserBrief{
			UserID:   post.UserID,
			Username: author.Username,
			Avatar:   author.AvatarURL, // will be gen by objects3 and fetch from s3
		}

		// 4. Fetch stats (stub now, real impl later)
		stats := PostStats{
			Likes:    0, // TODO: call LikeService
			Comments: 0, // TODO: call CommentService
		}

		feedItem := FeedItem{
			PostID:    post.ID,
			Author:    authorBrief,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			Media:     mediaList,
			Stats:     stats,
			// ViewerReaction: TODO: fetch from LikeService if current user reacted
		}

		feed = append(feed, feedItem)
	}

	return FeedResponse{
		Feed:       feed,
		NextOffset: offset + int64(len(postIDs)),
	}, nil
}
