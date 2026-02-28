package server

import (
	"time"

	"github.com/chanombude/twitter-go-api/internal/usecase"
)

type tweetResponse struct {
	ID           int64          `json:"id"`
	Content      *string        `json:"content"`
	MediaType    *string        `json:"media_type"`
	MediaUrl     *string        `json:"media_url"`
	User         userResponse   `json:"user"`
	ReplyCount   int32          `json:"reply_count"`
	LikeCount    int32          `json:"like_count"`
	RetweetCount int32          `json:"retweet_count"`
	IsLiked      bool           `json:"is_liked"`
	IsRetweeted  bool           `json:"is_retweeted"`
	Original     *tweetResponse `json:"original,omitempty"`
	ParentID     *int64         `json:"parent_id"`
	ParentHandle *string        `json:"parent_handle"`
	CreatedAt    time.Time      `json:"created_at"`
}

// Converts a Tweet DB model into a structured API response.
func newTweetResponse(tweet usecase.TweetItem, originalTweet *tweetResponse, parentHandle *string) tweetResponse {
	var content, mediaType, mediaurl *string
	if tweet.Content.Valid {
		content = &tweet.Content.String
	}
	if tweet.MediaType.Valid {
		mediaType = &tweet.MediaType.String
	}
	if tweet.MediaUrl.Valid {
		mediaurl = &tweet.MediaUrl.String
	}

	var parentId *int64
	if tweet.ParentID.Valid {
		parentId = &tweet.ParentID.Int64
	}

	return tweetResponse{
		ID:           tweet.ID,
		Content:      content,
		MediaType:    mediaType,
		MediaUrl:     mediaurl,
		User:         newUserResponse(tweet.Author),
		ReplyCount:   tweet.ReplyCount,
		LikeCount:    tweet.LikeCount,
		RetweetCount: tweet.RetweetCount,
		IsLiked:      tweet.IsLiked,
		IsRetweeted:  tweet.IsRetweeted,
		Original:     originalTweet,
		ParentID:     parentId,
		ParentHandle: parentHandle,
		CreatedAt:    tweet.CreatedAt,
	}
}
