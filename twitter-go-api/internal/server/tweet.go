package server

import (
	"time"

	"github.com/chanombude/twitter-go-api/internal/usecase"
)

type tweetResponse struct {
	ID             int64          `json:"id"`
	Content        *string        `json:"content"`
	MediaType      *string        `json:"mediaType"`
	MediaUrl       *string        `json:"mediaUrl"`
	User           userResponse   `json:"user"`
	ReplyCount     int32          `json:"replyCount"`
	LikeCount      int32          `json:"likeCount"`
	RetweetCount   int32          `json:"retweetCount"`
	IsLiked        bool           `json:"isLiked"`
	IsRetweeted    bool           `json:"isRetweeted"`
	Original       *tweetResponse `json:"original,omitempty"`
	ParentID       *int64         `json:"parentId"`
	ParentUsername *string        `json:"parentUsername"`
	CreatedAt      time.Time      `json:"createdAt"`
}

// Converts a Tweet DB model into a structured API response.
func newTweetResponse(tweet usecase.TweetItem) tweetResponse {
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

	var originalResp *tweetResponse
	if tweet.OriginalTweet != nil {
		orig := newTweetResponse(*tweet.OriginalTweet)
		originalResp = &orig
	}

	return tweetResponse{
		ID:             tweet.ID,
		Content:        content,
		MediaType:      mediaType,
		MediaUrl:       mediaurl,
		User:           newUserResponse(tweet.Author),
		ReplyCount:     tweet.ReplyCount,
		LikeCount:      tweet.LikeCount,
		RetweetCount:   tweet.RetweetCount,
		IsLiked:        tweet.IsLiked,
		IsRetweeted:    tweet.IsRetweeted,
		Original:       originalResp,
		ParentID:       parentId,
		ParentUsername: tweet.ParentUsername,
		CreatedAt:      tweet.CreatedAt,
	}
}
