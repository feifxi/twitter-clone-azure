package usecase

import (
	"context"
	"database/sql"
	"strings"

	"github.com/chanombude/twitter-go-api/internal/db"
)

func (u *Usecase) GetGlobalFeed(ctx context.Context, page, size int32, viewerID *int64) ([]TweetItem, error) {
	var vID sql.NullInt64
	if viewerID != nil {
		vID = sql.NullInt64{Int64: *viewerID, Valid: true}
	}
	rows, err := u.store.ListForYouFeed(ctx, db.ListForYouFeedParams{
		Limit:    size,
		Offset:   page * size,
		ViewerID: vID,
	})
	if err != nil {
		return nil, err
	}

	items := make([]TweetItem, 0, len(rows))
	for _, r := range rows {
		userRow, _ := u.store.GetUser(ctx, db.GetUserParams{ID: r.UserID, ViewerID: vID})
		items = append(items, TweetItem{
			Tweet: db.Tweet{
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
			Author: UserItem{
				User: db.User{
					ID:             userRow.ID,
					Username:       userRow.Username,
					Email:          userRow.Email,
					DisplayName:    userRow.DisplayName,
					Bio:            userRow.Bio,
					AvatarUrl:      userRow.AvatarUrl,
					Role:           userRow.Role,
					FollowersCount: userRow.FollowersCount,
					FollowingCount: userRow.FollowingCount,
					CreatedAt:      userRow.CreatedAt,
					UpdatedAt:      userRow.UpdatedAt,
				},
				IsFollowing: userRow.IsFollowing,
			},
			IsLiked:     r.IsLiked,
			IsRetweeted: r.IsRetweeted,
			IsFollowing: r.IsFollowing,
		})
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
	items := make([]TweetItem, 0, len(rows))
	for _, r := range rows {
		userRow, _ := u.store.GetUser(ctx, db.GetUserParams{ID: r.UserID, ViewerID: vID})
		items = append(items, TweetItem{
			Tweet: db.Tweet{
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
			Author: UserItem{
				User: db.User{
					ID:             userRow.ID,
					Username:       userRow.Username,
					Email:          userRow.Email,
					DisplayName:    userRow.DisplayName,
					Bio:            userRow.Bio,
					AvatarUrl:      userRow.AvatarUrl,
					Role:           userRow.Role,
					FollowersCount: userRow.FollowersCount,
					FollowingCount: userRow.FollowingCount,
					CreatedAt:      userRow.CreatedAt,
					UpdatedAt:      userRow.UpdatedAt,
				},
				IsFollowing: userRow.IsFollowing,
			},
			IsLiked:     r.IsLiked,
			IsRetweeted: r.IsRetweeted,
			IsFollowing: r.IsFollowing,
		})
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
	items := make([]TweetItem, 0, len(rows))
	for _, r := range rows {
		userRow, _ := u.store.GetUser(ctx, db.GetUserParams{ID: r.UserID, ViewerID: vID})
		items = append(items, TweetItem{
			Tweet: db.Tweet{
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
			Author: UserItem{
				User: db.User{
					ID:             userRow.ID,
					Username:       userRow.Username,
					Email:          userRow.Email,
					DisplayName:    userRow.DisplayName,
					Bio:            userRow.Bio,
					AvatarUrl:      userRow.AvatarUrl,
					Role:           userRow.Role,
					FollowersCount: userRow.FollowersCount,
					FollowingCount: userRow.FollowingCount,
					CreatedAt:      userRow.CreatedAt,
					UpdatedAt:      userRow.UpdatedAt,
				},
				IsFollowing: userRow.IsFollowing,
			},
			IsLiked:     r.IsLiked,
			IsRetweeted: r.IsRetweeted,
			IsFollowing: r.IsFollowing,
		})
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
		items := make([]TweetItem, 0, len(rows))
		for _, r := range rows {
			userRow, _ := u.store.GetUser(ctx, db.GetUserParams{ID: r.UserID, ViewerID: vID})
			items = append(items, TweetItem{
				Tweet: db.Tweet{
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
				Author: UserItem{
					User: db.User{
						ID:             userRow.ID,
						Username:       userRow.Username,
						Email:          userRow.Email,
						DisplayName:    userRow.DisplayName,
						Bio:            userRow.Bio,
						AvatarUrl:      userRow.AvatarUrl,
						Role:           userRow.Role,
						Provider:       userRow.Provider,
						FollowersCount: userRow.FollowersCount,
						FollowingCount: userRow.FollowingCount,
						CreatedAt:      userRow.CreatedAt,
						UpdatedAt:      userRow.UpdatedAt,
					},
					IsFollowing: userRow.IsFollowing,
				},
				IsLiked:     r.IsLiked,
				IsRetweeted: r.IsRetweeted,
				IsFollowing: r.IsFollowing,
			})
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
	items := make([]TweetItem, 0, len(rows))
	for _, r := range rows {
		userRow, _ := u.store.GetUser(ctx, db.GetUserParams{ID: r.UserID, ViewerID: vID})
		items = append(items, TweetItem{
			Tweet: db.Tweet{
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
			Author: UserItem{
				User: db.User{
					ID:             userRow.ID,
					Username:       userRow.Username,
					Email:          userRow.Email,
					DisplayName:    userRow.DisplayName,
					Bio:            userRow.Bio,
					AvatarUrl:      userRow.AvatarUrl,
					Role:           userRow.Role,
					Provider:       userRow.Provider,
					FollowersCount: userRow.FollowersCount,
					FollowingCount: userRow.FollowingCount,
					CreatedAt:      userRow.CreatedAt,
					UpdatedAt:      userRow.UpdatedAt,
				},
				IsFollowing: userRow.IsFollowing,
			},
			IsLiked:     r.IsLiked,
			IsRetweeted: r.IsRetweeted,
			IsFollowing: r.IsFollowing,
		})
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
	hashtags, err := u.store.GetTrendingHashtagsLast24h(ctx, limit)
	if err != nil {
		return nil, err
	}
	if len(hashtags) == 0 {
		return u.store.GetTopHashtagsAllTime(ctx, limit)
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
