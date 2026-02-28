package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/chanombude/twitter-go-api/internal/apperr"
	"github.com/chanombude/twitter-go-api/internal/db"
)

type MediaUpload struct {
	Filename    string
	ContentType string
	Reader      interface {
		Read(p []byte) (n int, err error)
	}
}

type CreateTweetInput struct {
	UserID   int64
	Content  *string
	ParentID *int64
	Media    *MediaUpload
}

func (u *Usecase) CreateTweet(ctx context.Context, input CreateTweetInput) (TweetItem, error) {
	trimmedContent := ""
	if input.Content != nil {
		trimmedContent = strings.TrimSpace(*input.Content)
	}

	mediaType := sql.NullString{String: "NONE", Valid: true}
	mediaURL := sql.NullString{Valid: false}
	if input.Media != nil {
		contentType := strings.ToLower(input.Media.ContentType)
		switch {
		case strings.HasPrefix(contentType, "image/"):
			mediaType = sql.NullString{String: "IMAGE", Valid: true}
		case strings.HasPrefix(contentType, "video/"):
			mediaType = sql.NullString{String: "VIDEO", Valid: true}
		default:
			return TweetItem{}, apperr.BadRequest("only images or videos are allowed")
		}

		uploadedURL, err := u.storage.UploadFile(ctx, input.Media.Reader, input.Media.Filename, contentType)
		if err != nil {
			return TweetItem{}, err
		}
		mediaURL = sql.NullString{String: uploadedURL, Valid: true}
	}

	if trimmedContent == "" && !mediaURL.Valid {
		return TweetItem{}, apperr.BadRequest("tweet must include text or media")
	}

	parentID := sql.NullInt64{Valid: false}
	if input.ParentID != nil {
		if _, err := u.store.GetTweet(ctx, db.GetTweetParams{ID: *input.ParentID}); err != nil {
			return TweetItem{}, err
		}
		parentID = sql.NullInt64{Int64: *input.ParentID, Valid: true}
	}

	content := sql.NullString{Valid: false}
	if trimmedContent != "" {
		content = sql.NullString{String: trimmedContent, Valid: true}
	}

	tweet, err := u.store.CreateTweet(ctx, db.CreateTweetParams{
		UserID:    input.UserID,
		Content:   content,
		MediaType: mediaType,
		MediaUrl:  mediaURL,
		ParentID:  parentID,
		RetweetID: sql.NullInt64{Valid: false},
	})
	if err != nil {
		if mediaURL.Valid {
			_ = u.storage.DeleteFile(ctx, mediaURL.String)
		}
		return TweetItem{}, err
	}

	if parentID.Valid {
		if err := u.store.IncrementParentReplyCount(ctx, parentID.Int64); err != nil {
			log.Printf("failed to increment parent reply count: %v", err)
		}

		parentTweet, err := u.store.GetTweet(ctx, db.GetTweetParams{ID: parentID.Int64})
		if err == nil {
			tweetID := tweet.ID
			_ = u.createAndDispatchNotification(ctx, parentTweet.UserID, input.UserID, &tweetID, "REPLY")
		}
	}

	if content.Valid {
		tags := extractHashtags(content.String)
		for _, tag := range tags {
			h, err := u.store.UpsertHashtag(ctx, tag)
			if err != nil {
				log.Printf("failed to upsert hashtag %s: %v", tag, err)
				continue
			}
			if err := u.store.LinkTweetHashtag(ctx, db.LinkTweetHashtagParams{TweetID: tweet.ID, HashtagID: h.ID}); err != nil {
				log.Printf("failed to link hashtag %s: %v", tag, err)
			}
		}
	}

	return u.GetTweet(ctx, tweet.ID, &input.UserID)
}

func (u *Usecase) GetTweet(ctx context.Context, tweetID int64, viewerID *int64) (TweetItem, error) {
	var vID sql.NullInt64
	if viewerID != nil {
		vID = sql.NullInt64{Int64: *viewerID, Valid: true}
	}
	r, err := u.store.GetTweet(ctx, db.GetTweetParams{ID: tweetID, ViewerID: vID})
	if err != nil {
		return TweetItem{}, err
	}
	tweets := []db.Tweet{
		{
			ID:           r.ID,
			UserID:       r.UserID,
			Content:      r.Content,
			MediaType:    r.MediaType,
			MediaUrl:     r.MediaUrl,
			ParentID:     r.ParentID,
			RetweetID:    r.RetweetID,
			ReplyCount:   r.ReplyCount,
			RetweetCount: r.RetweetCount,
			LikeCount:    r.LikeCount,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		},
	}
	items, err := u.populateTweetItems(ctx, tweets, viewerID)
	if err != nil || len(items) == 0 {
		return TweetItem{}, err
	}

	items[0].IsLiked = r.IsLiked
	items[0].IsRetweeted = r.IsRetweeted
	items[0].IsFollowing = r.IsFollowing

	return items[0], nil
}

func (u *Usecase) DeleteTweet(ctx context.Context, userID, tweetID int64) error {
	tweet, err := u.store.GetTweet(ctx, db.GetTweetParams{ID: tweetID})
	if err != nil {
		return err
	}
	if tweet.UserID != userID {
		return apperr.Forbidden("you can only delete your own tweets")
	}

	if tweet.RetweetID.Valid {
		_, err = u.store.DeleteRetweetByUser(ctx, db.DeleteRetweetByUserParams{
			UserID:    userID,
			RetweetID: sql.NullInt64{Int64: tweet.RetweetID.Int64, Valid: true},
		})
	} else {
		_, err = u.store.DeleteTweetByOwner(ctx, db.DeleteTweetByOwnerParams{ID: tweetID, UserID: userID})
	}
	if err != nil {
		return err
	}

	if tweet.ParentID.Valid {
		if err := u.store.DecrementParentReplyCount(ctx, tweet.ParentID.Int64); err != nil {
			log.Printf("failed to decrement parent reply count: %v", err)
		}
	}

	if tweet.MediaUrl.Valid && tweet.MediaUrl.String != "" {
		_ = u.storage.DeleteFile(ctx, tweet.MediaUrl.String)
	}

	return nil
}

func (u *Usecase) ListReplies(ctx context.Context, tweetID int64, page, size int32, viewerID *int64) ([]TweetItem, error) {
	var vID sql.NullInt64
	if viewerID != nil {
		vID = sql.NullInt64{Int64: *viewerID, Valid: true}
	}
	rows, err := u.store.ListTweetReplies(ctx, db.ListTweetRepliesParams{
		ParentID: sql.NullInt64{Int64: tweetID, Valid: true},
		Limit:    size,
		Offset:   page * size,
		ViewerID: vID,
	})
	if err != nil {
		return nil, err
	}
	// Map raw rows to db.Tweet slice so it can be passed to dataloader
	tweets := make([]db.Tweet, 0, len(rows))
	for _, r := range rows {
		tweets = append(tweets, db.Tweet{
			ID:           r.ID,
			UserID:       r.UserID,
			Content:      r.Content,
			MediaType:    r.MediaType,
			MediaUrl:     r.MediaUrl,
			ParentID:     r.ParentID,
			RetweetID:    r.RetweetID,
			ReplyCount:   r.ReplyCount,
			RetweetCount: r.RetweetCount,
			LikeCount:    r.LikeCount,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
	}

	items, err := u.populateTweetItems(ctx, tweets, viewerID)
	if err != nil {
		return nil, err
	}

	for i, r := range rows {
		items[i].IsLiked = r.IsLiked
		items[i].IsRetweeted = r.IsRetweeted
		items[i].IsFollowing = r.IsFollowing
	}

	return items, nil
}

func (u *Usecase) LikeTweet(ctx context.Context, userID, tweetID int64) error {
	tweet, err := u.store.GetTweet(ctx, db.GetTweetParams{ID: tweetID})
	if err != nil {
		return err
	}

	liked, err := u.store.LikeTweet(ctx, db.LikeTweetParams{UserID: userID, TweetID: tweetID})
	if err != nil {
		return err
	}

	if liked {
		id := tweet.ID
		_ = u.createAndDispatchNotification(ctx, tweet.UserID, userID, &id, "LIKE")
	}

	return nil
}

func (u *Usecase) UnlikeTweet(ctx context.Context, userID, tweetID int64) error {
	_, err := u.store.UnlikeTweet(ctx, db.UnlikeTweetParams{UserID: userID, TweetID: tweetID})
	return err
}

func (u *Usecase) Retweet(ctx context.Context, userID, tweetID int64) (TweetItem, error) {
	targetTweet, err := u.store.GetTweet(ctx, db.GetTweetParams{ID: tweetID})
	if err != nil {
		return TweetItem{}, err
	}

	originalTweet := targetTweet
	if targetTweet.RetweetID.Valid {
		originalTweet, err = u.store.GetTweet(ctx, db.GetTweetParams{ID: targetTweet.RetweetID.Int64})
		if err != nil {
			return TweetItem{}, err
		}
	}

	created, err := u.store.CreateRetweet(ctx, db.CreateRetweetParams{
		UserID:    userID,
		RetweetID: sql.NullInt64{Int64: originalTweet.ID, Valid: true},
	})
	if err != nil {
		return TweetItem{}, err
	}

	id := originalTweet.ID
	_ = u.createAndDispatchNotification(ctx, originalTweet.UserID, userID, &id, "RETWEET")

	return u.GetTweet(ctx, created.ID, &userID)
}

func (u *Usecase) UndoRetweet(ctx context.Context, userID, tweetID int64) error {
	targetTweet, err := u.store.GetTweet(ctx, db.GetTweetParams{ID: tweetID})
	if err != nil {
		return err
	}

	originalID := targetTweet.ID
	if targetTweet.RetweetID.Valid {
		originalID = targetTweet.RetweetID.Int64
	}

	_, err = u.store.DeleteRetweetByUser(ctx, db.DeleteRetweetByUserParams{
		UserID:    userID,
		RetweetID: sql.NullInt64{Int64: originalID, Valid: true},
	})
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return err
}
