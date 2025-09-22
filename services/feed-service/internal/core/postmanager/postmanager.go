package postmanager

import (
	"feedservice/internal/infra/store"
	"feedservice/internal/model"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type PostManager struct {
	MediaStore        *store.MediaStore
	PostStore         *store.PostStore
	PostMediaStore    *store.PostMediaStore
	newposteventqueue *chan model.NewPostEvent
}

func NewPostManager(mediaStore *store.MediaStore, newposteventqueue_ *chan model.NewPostEvent) *PostManager {
	return &PostManager{
		MediaStore:        mediaStore,
		newposteventqueue: newposteventqueue_,
	}
}

func (p *PostManager) CreatePost(userID string, content string, mediaIDs []string) (string, error) {
	// 1. Validate mediaIDs belong to this user
	if err := p.MediaStore.ValidateUserMedia(userID, mediaIDs); err != nil {
		return "", err
	}

	// 2. Insert post record into posts table
	postID := uuid.New().String()
	err := p.PostStore.InsertPost(postID, userID, content)
	if err != nil {
		return "", err
	}

	// 3. Link media to post in post_media table
	if len(mediaIDs) > 0 {
		if err := p.PostMediaStore.LinkMediaToPost(postID, mediaIDs); err != nil {
			return "", err
		}

		// 4. Update media status to 'uploaded'
		if err := p.MediaStore.UpdateMediaStatus(mediaIDs, "uploaded"); err != nil {
			return "", err
		}
	}

	if p.newposteventqueue != nil {
		*p.newposteventqueue <- model.NewPostEvent{
			PostID: postID,
			UserID: userID,
		}
	}

	// 5. Return new post ID
	return postID, nil
}

func (p *PostManager) CreateMedia(userID string, mediatype string, mediafilename string) (model.Media, error) {
	// Step 1: Generate new media_id
	mediaID := uuid.New().String()

	// Step 2: Build S3 object key (random name under user folder)
	objectKey := fmt.Sprintf("%s/%s_%s", userID, mediaID, mediafilename)

	// Step 3: Generate pre-signed upload URL (valid 5 min)
	uploadURL := p.MediaStore.S3client.GeneratePreSignedURL(objectKey, 5*time.Minute)

	// Step 4: Insert into DB with status=pending
	media, err := p.MediaStore.CreateMediaRecord(
		mediaID,
		userID,
		mediatype,
		"pending",
		objectKey, // use this key to gen GET url S3 media in later.
	)
	if err != nil {
		return model.Media{}, fmt.Errorf("failed to insert media record: %w", err)
	}

	// Step 5: Return Media with presigned URL
	media.Url.String = uploadURL
	media.Url.Valid = true

	return media, nil
}
