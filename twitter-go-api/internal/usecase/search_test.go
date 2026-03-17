package usecase_test

import (
	"context"
	"testing"

	"github.com/chanombude/twitter-go-api/internal/db"
	"github.com/chanombude/twitter-go-api/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestSearchUsecase_SearchUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		store := &MockStore{
			SearchUsersFn: func(ctx context.Context, arg db.SearchUsersParams) ([]db.SearchUsersRow, error) {
				return []db.SearchUsersRow{
					{User: db.User{ID: 1, Username: "founduser"}},
				}, nil
			},
		}

		uc := usecase.NewSearchUsecase(store)
		res, err := uc.SearchUsers(ctx, "query", 0, 10, nil)

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, "founduser", res[0].Username)
	})

	t.Run("empty_query", func(t *testing.T) {
		uc := usecase.NewSearchUsecase(nil)
		res, err := uc.SearchUsers(ctx, "  ", 0, 10, nil)
		require.NoError(t, err)
		require.Empty(t, res)
	})
}

func TestSearchUsecase_SearchTweets(t *testing.T) {
	ctx := context.Background()

	t.Run("hashtag_search", func(t *testing.T) {
		store := &MockStore{
			SearchTweetsByHashtagFn: func(ctx context.Context, arg db.SearchTweetsByHashtagParams) ([]db.SearchTweetsByHashtagRow, error) {
				return []db.SearchTweetsByHashtagRow{
					{Tweet: db.Tweet{ID: 10, UserID: 100}},
				}, nil
			},
			GetUsersByIDsFn: func(ctx context.Context, arg db.GetUsersByIDsParams) ([]db.GetUsersByIDsRow, error) {
				return []db.GetUsersByIDsRow{
					{User: db.User{ID: 100, Username: "author"}},
				}, nil
			},
		}

		uc := usecase.NewSearchUsecase(store)
		res, err := uc.SearchTweets(ctx, "#Golang", 0, 10, nil)

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, int64(10), res[0].ID)
	})

	t.Run("fulltext_search", func(t *testing.T) {
		store := &MockStore{
			SearchTweetsFullTextFn: func(ctx context.Context, arg db.SearchTweetsFullTextParams) ([]db.SearchTweetsFullTextRow, error) {
				return []db.SearchTweetsFullTextRow{
					{Tweet: db.Tweet{ID: 20, UserID: 200}},
				}, nil
			},
			GetUsersByIDsFn: func(ctx context.Context, arg db.GetUsersByIDsParams) ([]db.GetUsersByIDsRow, error) {
				return []db.GetUsersByIDsRow{
					{User: db.User{ID: 200, Username: "author2"}},
				}, nil
			},
		}

		uc := usecase.NewSearchUsecase(store)
		res, err := uc.SearchTweets(ctx, "hello world", 0, 10, nil)

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, int64(20), res[0].ID)
	})
}

func TestSearchUsecase_SearchHashtags(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		store := &MockStore{
			SearchHashtagsByPrefixFn: func(ctx context.Context, arg db.SearchHashtagsByPrefixParams) ([]db.Hashtag, error) {
				return []db.Hashtag{{Text: "golang"}}, nil
			},
		}

		uc := usecase.NewSearchUsecase(store)
		res, err := uc.SearchHashtags(ctx, "go", 5)

		require.NoError(t, err)
		require.Len(t, res, 1)
		require.Equal(t, "golang", res[0].Text)
	})
}
