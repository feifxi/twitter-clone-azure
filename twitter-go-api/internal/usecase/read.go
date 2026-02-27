package usecase

import (
	"context"
	"database/sql"
	"strings"

	"github.com/chanombude/twitter-go-api/internal/db"
)

func (u *Usecase) GetGlobalFeed(ctx context.Context, page, size int32) ([]db.Tweet, error) {
	return u.store.ListForYouFeed(ctx, db.ListForYouFeedParams{Limit: size, Offset: page * size})
}

func (u *Usecase) GetFollowingFeed(ctx context.Context, userID int64, page, size int32) ([]db.Tweet, error) {
	return u.store.ListFollowingFeed(ctx, db.ListFollowingFeedParams{FollowerID: userID, Limit: size, Offset: page * size})
}

func (u *Usecase) GetUserFeed(ctx context.Context, userID int64, page, size int32) ([]db.Tweet, error) {
	return u.store.ListUserTweets(ctx, db.ListUserTweetsParams{UserID: userID, Limit: size, Offset: page * size})
}

func (u *Usecase) SearchUsers(ctx context.Context, query string, page, size int32, viewerID *int64) ([]db.User, map[int64]bool, error) {
	users, err := u.store.SearchUsers(ctx, db.SearchUsersParams{
		Column1: sql.NullString{String: strings.TrimSpace(query), Valid: true},
		Limit:   size,
		Offset:  page * size,
	})
	if err != nil {
		return nil, nil, err
	}

	followingMap := make(map[int64]bool)
	if viewerID != nil {
		for _, user := range users {
			if user.ID == *viewerID {
				continue
			}
			f, err := u.store.IsFollowing(ctx, db.IsFollowingParams{FollowerID: *viewerID, FollowingID: user.ID})
			if err == nil {
				followingMap[user.ID] = f
			}
		}
	}

	return users, followingMap, nil
}

func (u *Usecase) SearchTweets(ctx context.Context, query string, page, size int32) ([]db.Tweet, error) {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return []db.Tweet{}, nil
	}

	if strings.HasPrefix(trimmed, "#") {
		hashtag := strings.TrimSpace(strings.ToLower(strings.TrimLeft(trimmed, "#")))
		if hashtag == "" {
			return []db.Tweet{}, nil
		}
		return u.store.SearchTweetsByHashtag(ctx, db.SearchTweetsByHashtagParams{Lower: hashtag, Limit: size, Offset: page * size})
	}

	tsQuery := buildTSQuery(trimmed)
	if tsQuery == "" {
		return []db.Tweet{}, nil
	}
	return u.store.SearchTweetsFullText(ctx, db.SearchTweetsFullTextParams{ToTsquery: tsQuery, Limit: size, Offset: page * size})
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

func (u *Usecase) GetSuggestedUsers(ctx context.Context, page, size int32, viewerID *int64) ([]db.User, map[int64]bool, error) {
	var (
		users []db.User
		err   error
	)
	if viewerID != nil {
		users, err = u.store.ListSuggestedUsers(ctx, db.ListSuggestedUsersParams{FollowerID: *viewerID, Limit: size, Offset: page * size})
	} else {
		users, err = u.store.ListTopUsers(ctx, db.ListTopUsersParams{Limit: size, Offset: page * size})
	}
	if err != nil {
		return nil, nil, err
	}

	followingMap := make(map[int64]bool)
	if viewerID != nil {
		for _, user := range users {
			if user.ID == *viewerID {
				continue
			}
			f, err := u.store.IsFollowing(ctx, db.IsFollowingParams{FollowerID: *viewerID, FollowingID: user.ID})
			if err == nil {
				followingMap[user.ID] = f
			}
		}
	}

	return users, followingMap, nil
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
