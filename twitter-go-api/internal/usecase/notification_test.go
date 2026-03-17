package usecase_test

import (
	"context"
	"testing"

	"github.com/chanombude/twitter-go-api/internal/db"
	"github.com/chanombude/twitter-go-api/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestNotificationUsecase_ListNotifications(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)

	t.Run("success", func(t *testing.T) {
		tweetID := int64(10)
		store := &MockStore{
			ListNotificationsFn: func(ctx context.Context, arg db.ListNotificationsParams) ([]db.Notification, error) {
				return []db.Notification{
					{ID: 100, RecipientID: userID, ActorID: 2, Type: "like", TweetID: &tweetID},
				}, nil
			},
			GetUsersByIDsFn: func(ctx context.Context, arg db.GetUsersByIDsParams) ([]db.GetUsersByIDsRow, error) {
				return []db.GetUsersByIDsRow{
					{User: db.User{ID: 2, Username: "actor"}},
				}, nil
			},
			GetTweetsByIDsFn: func(ctx context.Context, arg db.GetTweetsByIDsParams) ([]db.GetTweetsByIDsRow, error) {
				content := "tweet content"
				return []db.GetTweetsByIDsRow{
					{Tweet: db.Tweet{ID: tweetID, Content: &content}},
				}, nil
			},
		}

		uc := usecase.NewNotificationUsecase(store)
		res, err := uc.ListNotifications(ctx, userID, 0, 10)

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, "actor", res[0].Actor.Username)
		require.Equal(t, "tweet content", *res[0].TweetContent)
	})
}

func TestNotificationUsecase_Counters(t *testing.T) {
	ctx := context.Background()
	userID := int64(1)

	t.Run("get_count", func(t *testing.T) {
		store := &MockStore{
			GetUnreadNotificationCountFn: func(ctx context.Context, id int64) (int64, error) {
				return 5, nil
			},
		}
		uc := usecase.NewNotificationUsecase(store)
		count, err := uc.CountUnreadNotifications(ctx, userID)
		require.NoError(t, err)
		require.Equal(t, int64(5), count)
	})

	t.Run("mark_all_read", func(t *testing.T) {
		called := false
		store := &MockStore{
			MarkAllNotificationsReadFn: func(ctx context.Context, id int64) error {
				called = true
				require.Equal(t, userID, id)
				return nil
			},
		}
		uc := usecase.NewNotificationUsecase(store)
		err := uc.MarkAllNotificationsRead(ctx, userID)
		require.NoError(t, err)
		require.True(t, called)
	})
}
