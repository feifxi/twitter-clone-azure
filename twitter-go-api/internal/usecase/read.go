package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/chanombude/twitter-go-api/internal/db"
)

func (u *Usecase) GetGlobalFeed(ctx context.Context, page, size int32, viewerID *int64) ([]TweetItem, error) {
	var vID sql.NullInt64
	if viewerID != nil {
		vID = sql.NullInt64{Int64: *viewerID, Valid: true}
	}
	cacheKey := fmt.Sprintf("cache:global_feed:page:%d:size:%d", page, size)
	var rows []db.ListForYouFeedRow
	fromCache := false

	if viewerID == nil && u.redis != nil {
		if cached, err := u.redis.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			if json.Unmarshal([]byte(cached), &rows) == nil {
				fromCache = true
			}
		}
	}

	if !fromCache {
		var err error
		rows, err = u.store.ListForYouFeed(ctx, db.ListForYouFeedParams{
			Limit:    size,
			Offset:   page * size,
			ViewerID: vID,
		})
		if err != nil {
			return nil, err
		}

		if viewerID == nil && u.redis != nil && len(rows) > 0 {
			if data, err := json.Marshal(rows); err == nil {
				u.redis.Set(ctx, cacheKey, data, 30*time.Second)
			}
		}
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

	// Attach query-specific dynamic fields (IsLiked, IsRetweeted, IsFollowing)
	// Because dataloader doesn't know about the `rows` specific joined booleans.
	for i, r := range rows {
		items[i].IsLiked = r.IsLiked
		items[i].IsRetweeted = r.IsRetweeted
		items[i].IsFollowing = r.IsFollowing
	}

	return items, nil
}

func (u *Usecase) GetFollowingFeed(ctx context.Context, userID int64, page, size int32) ([]TweetItem, error) {
	vID := sql.NullInt64{Int64: userID, Valid: true}
	rows, err := u.store.ListFollowingFeed(ctx, db.ListFollowingFeedParams{
		FollowerID: userID,
		Limit:      size,
		Offset:     page * size,
		ViewerID:   vID,
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

	items, err := u.populateTweetItems(ctx, tweets, &userID)
	if err != nil {
		return nil, err
	}

	// Attach query-specific dynamic fields (IsLiked, IsRetweeted, IsFollowing)
	// Because dataloader doesn't know about the `rows` specific joined booleans.
	for i, r := range rows {
		items[i].IsLiked = r.IsLiked
		items[i].IsRetweeted = r.IsRetweeted
		items[i].IsFollowing = r.IsFollowing
	}

	return items, nil
}

func (u *Usecase) GetUserFeed(ctx context.Context, userID int64, page, size int32, viewerID *int64) ([]TweetItem, error) {
	var vID sql.NullInt64
	if viewerID != nil {
		vID = sql.NullInt64{Int64: *viewerID, Valid: true}
	}
	rows, err := u.store.ListUserTweets(ctx, db.ListUserTweetsParams{
		UserID:   userID,
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

	// Attach query-specific dynamic fields (IsLiked, IsRetweeted, IsFollowing)
	// Because dataloader doesn't know about the `rows` specific joined booleans.
	for i, r := range rows {
		items[i].IsLiked = r.IsLiked
		items[i].IsRetweeted = r.IsRetweeted
		items[i].IsFollowing = r.IsFollowing
	}

	return items, nil
}

func (u *Usecase) SearchUsers(ctx context.Context, query string, page, size int32, viewerID *int64) ([]UserItem, error) {
	var vID sql.NullInt64
	if viewerID != nil {
		vID = sql.NullInt64{Int64: *viewerID, Valid: true}
	}
	rows, err := u.store.SearchUsers(ctx, db.SearchUsersParams{
		Column1:  sql.NullString{String: strings.TrimSpace(query), Valid: true},
		Limit:    size,
		Offset:   page * size,
		ViewerID: vID,
	})
	if err != nil {
		return nil, err
	}
	items := make([]UserItem, 0, len(rows))
	for _, r := range rows {
		items = append(items, UserItem{
			User: db.User{
				ID:             r.ID,
				Username:       r.Username,
				Email:          r.Email,
				DisplayName:    r.DisplayName,
				Bio:            r.Bio,
				AvatarUrl:      r.AvatarUrl,
				Role:           r.Role,
				Provider:       r.Provider,
				FollowersCount: r.FollowersCount,
				FollowingCount: r.FollowingCount,
				CreatedAt:      r.CreatedAt,
				UpdatedAt:      r.UpdatedAt,
			},
			IsFollowing: r.IsFollowing,
		})
	}
	return items, nil
}

func (u *Usecase) SearchTweets(ctx context.Context, query string, page, size int32, viewerID *int64) ([]TweetItem, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return []TweetItem{}, nil
	}

	var vID sql.NullInt64
	if viewerID != nil {
		vID = sql.NullInt64{Int64: *viewerID, Valid: true}
	}

	if strings.HasPrefix(trimmed, "#") {
		hashtag := strings.TrimSpace(strings.ToLower(strings.TrimLeft(trimmed, "#")))
		if hashtag == "" {
			return []TweetItem{}, nil
		}
		rows, err := u.store.SearchTweetsByHashtag(ctx, db.SearchTweetsByHashtagParams{
			Lower:    hashtag,
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

		// Attach query-specific dynamic fields (IsLiked, IsRetweeted, IsFollowing)
		// Because dataloader doesn't know about the `rows` specific joined booleans.
		for i, r := range rows {
			items[i].IsLiked = r.IsLiked
			items[i].IsRetweeted = r.IsRetweeted
			items[i].IsFollowing = r.IsFollowing
		}

		return items, nil
	}

	tsQuery := buildTSQuery(trimmed)
	if tsQuery == "" {
		return []TweetItem{}, nil
	}
	rows, err := u.store.SearchTweetsFullText(ctx, db.SearchTweetsFullTextParams{
		ToTsquery: tsQuery,
		Limit:     size,
		Offset:    page * size,
		ViewerID:  vID,
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

	// Attach query-specific dynamic fields (IsLiked, IsRetweeted, IsFollowing)
	// Because dataloader doesn't know about the `rows` specific joined booleans.
	for i, r := range rows {
		items[i].IsLiked = r.IsLiked
		items[i].IsRetweeted = r.IsRetweeted
		items[i].IsFollowing = r.IsFollowing
	}

	return items, nil
}

func (u *Usecase) SearchHashtags(ctx context.Context, query string, limit int32) ([]db.Hashtag, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return []db.Hashtag{}, nil
	}
	return u.store.SearchHashtagsByPrefix(ctx, db.SearchHashtagsByPrefixParams{
		Column1: sql.NullString{String: strings.TrimPrefix(trimmed, "#"), Valid: true},
		Limit:   limit,
	})
}

func (u *Usecase) GetTrendingHashtags(ctx context.Context, limit int32) ([]db.Hashtag, error) {
	cacheKey := fmt.Sprintf("cache:trending_hashtags:limit:%d", limit)
	if u.redis != nil {
		if cached, err := u.redis.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
			var hashtags []db.Hashtag
			if json.Unmarshal([]byte(cached), &hashtags) == nil {
				return hashtags, nil
			}
		}
	}

	hashtags, err := u.store.GetTrendingHashtagsLast24h(ctx, limit)
	if err != nil {
		return nil, err
	}
	if len(hashtags) == 0 {
		if hashtags, err = u.store.GetTopHashtagsAllTime(ctx, limit); err != nil {
			return nil, err
		}
	}

	if u.redis != nil && len(hashtags) > 0 {
		if data, err := json.Marshal(hashtags); err == nil {
			u.redis.Set(ctx, cacheKey, data, 5*time.Minute)
		}
	}

	return hashtags, nil
}

func (u *Usecase) GetSuggestedUsers(ctx context.Context, page, size int32, viewerID *int64) ([]UserItem, error) {
	if viewerID != nil {
		rows, err := u.store.ListSuggestedUsers(ctx, db.ListSuggestedUsersParams{
			FollowerID: *viewerID,
			Limit:      size,
			Offset:     page * size,
			ViewerID:   sql.NullInt64{Int64: *viewerID, Valid: true},
		})
		if err != nil {
			return nil, err
		}
		items := make([]UserItem, 0, len(rows))
		for _, r := range rows {
			items = append(items, UserItem{
				User: db.User{
					ID:             r.ID,
					Username:       r.Username,
					Email:          r.Email,
					DisplayName:    r.DisplayName,
					Bio:            r.Bio,
					AvatarUrl:      r.AvatarUrl,
					Role:           r.Role,
					Provider:       r.Provider,
					FollowersCount: r.FollowersCount,
					FollowingCount: r.FollowingCount,
					CreatedAt:      r.CreatedAt,
					UpdatedAt:      r.UpdatedAt,
				},
				IsFollowing: r.IsFollowing,
			})
		}
		return items, nil
	} else {
		users, err := u.store.ListTopUsers(ctx, db.ListTopUsersParams{Limit: size, Offset: page * size})
		if err != nil {
			return nil, err
		}
		items := make([]UserItem, 0, len(users))
		for _, r := range users {
			items = append(items, UserItem{
				User:        r,
				IsFollowing: false,
			})
		}
		return items, nil
	}
}

func (u *Usecase) ListNotifications(ctx context.Context, userID int64, page, size int32) ([]db.Notification, error) {
	return u.store.ListNotifications(ctx, db.ListNotificationsParams{RecipientID: userID, Limit: size, Offset: page * size})
}

func (u *Usecase) CountUnreadNotifications(ctx context.Context, userID int64) (int64, error) {
	return u.store.GetUnreadNotificationCount(ctx, userID)
}

func (u *Usecase) MarkAllNotificationsRead(ctx context.Context, userID int64) error {
	return u.store.MarkAllNotificationsRead(ctx, userID)
}
