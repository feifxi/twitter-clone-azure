package usecase_test

import (
	"context"
	"testing"

	"github.com/chanombude/twitter-go-api/internal/db"
	"github.com/chanombude/twitter-go-api/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestDiscoveryUsecase_GetTrendingHashtags(t *testing.T) {
	ctx := context.Background()

	t.Run("success_24h", func(t *testing.T) {
		store := &MockStore{
			GetTrendingHashtagsLast24hFn: func(ctx context.Context, limit int32) ([]db.Hashtag, error) {
				return []db.Hashtag{{Text: "trending"}}, nil
			},
		}
		uc := usecase.NewDiscoveryUsecase(store)
		res, err := uc.GetTrendingHashtags(ctx, 5)
		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, "trending", res[0].Text)
	})

	t.Run("fallback_alltime", func(t *testing.T) {
		store := &MockStore{
			GetTrendingHashtagsLast24hFn: func(ctx context.Context, limit int32) ([]db.Hashtag, error) {
				return []db.Hashtag{}, nil
			},
			GetTopHashtagsAllTimeFn: func(ctx context.Context, limit int32) ([]db.Hashtag, error) {
				return []db.Hashtag{{Text: "alltime"}}, nil
			},
		}
		uc := usecase.NewDiscoveryUsecase(store)
		res, err := uc.GetTrendingHashtags(ctx, 5)
		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, "alltime", res[0].Text)
	})
}

func TestDiscoveryUsecase_GetSuggestedUsers(t *testing.T) {
	ctx := context.Background()
	viewerID := int64(1)

	t.Run("personalized", func(t *testing.T) {
		store := &MockStore{
			ListSuggestedUsersFn: func(ctx context.Context, arg db.ListSuggestedUsersParams) ([]db.ListSuggestedUsersRow, error) {
				return []db.ListSuggestedUsersRow{
					{User: db.User{ID: 100, Username: "suggested"}},
				}, nil
			},
		}
		uc := usecase.NewDiscoveryUsecase(store)
		res, err := uc.GetSuggestedUsers(ctx, 0, 10, &viewerID)
		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, "suggested", res[0].Username)
	})

	t.Run("guest_fallback", func(t *testing.T) {
		store := &MockStore{
			ListTopUsersFn: func(ctx context.Context, arg db.ListTopUsersParams) ([]db.User, error) {
				return []db.User{{ID: 200, Username: "top"}}, nil
			},
		}
		uc := usecase.NewDiscoveryUsecase(store)
		res, err := uc.GetSuggestedUsers(ctx, 0, 10, nil)
		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, "top", res[0].Username)
	})
}
