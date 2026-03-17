package usecase_test

import (
	"context"
	"testing"

	"github.com/chanombude/twitter-go-api/internal/db"
	"github.com/chanombude/twitter-go-api/internal/usecase"
)

type mockStorage struct {
	publicURLFn  func(key string) string
	deleteFileFn func(ctx context.Context, key string) error
}

func (m *mockStorage) PublicURL(key string) string {
	return m.publicURLFn(key)
}
func (m *mockStorage) DeleteFile(ctx context.Context, key string) error {
	return m.deleteFileFn(ctx, key)
}
func (m *mockStorage) GeneratePresignedURL(ctx context.Context, filename, contentType, folder string) (string, string, error) {
	return "", "", nil
}

func TestUpdateProfile_ClearingFields(t *testing.T) {
	ctx := context.Background()

	oldAvatar := "https://cdn.com/old_avatar.png"
	oldBanner := "https://cdn.com/old_banner.png"

	store := &MockStore{
		GetUserFn: func(ctx context.Context, arg db.GetUserParams) (db.GetUserRow, error) {
			return db.GetUserRow{
				User: db.User{
					ID:        arg.ID,
					AvatarUrl: &oldAvatar,
					BannerUrl: &oldBanner,
				},
			}, nil
		},
		UpdateUserProfileFn: func(ctx context.Context, arg db.UpdateUserProfileParams) (db.User, error) {
			if arg.AvatarUrl != nil {
				t.Errorf("expected AvatarUrl to be nil, got %v", *arg.AvatarUrl)
			}
			if arg.BannerUrl != nil {
				t.Errorf("expected BannerUrl to be nil, got %v", *arg.BannerUrl)
			}
			return db.User{ID: arg.ID}, nil
		},
	}

	storage := &mockStorage{
		deleteFileFn: func(ctx context.Context, key string) error {
			return nil
		},
	}

	uc := usecase.NewUserUsecase(store, storage, nil)

	empty := ""
	input := usecase.UpdateProfileInput{
		AvatarKey: &empty,
		BannerKey: &empty,
	}

	_, err := uc.UpdateProfile(ctx, 1, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFollowUser(t *testing.T) {
	t.Run("cannot_follow_yourself", func(t *testing.T) {
		uc := usecase.NewUserUsecase(nil, nil, nil)
		_, err := uc.FollowUser(context.Background(), 1, 1)
		if err == nil {
			t.Fatal("expected error for self-follow")
		}
	})

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		store := &MockStore{
			GetUserFn: func(ctx context.Context, arg db.GetUserParams) (db.GetUserRow, error) {
				return db.GetUserRow{User: db.User{ID: 2}}, nil
			},
			FollowUserFn: func(ctx context.Context, arg db.FollowUserParams) (bool, error) {
				return true, nil
			},
		}
		uc := usecase.NewUserUsecase(store, nil, nil)
		inserted, err := uc.FollowUser(ctx, 1, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !inserted {
			t.Fatal("expected inserted to be true")
		}
	})
}

func TestUnfollowUser(t *testing.T) {
	ctx := context.Background()
	store := &MockStore{
		GetUserFn: func(ctx context.Context, arg db.GetUserParams) (db.GetUserRow, error) {
			return db.GetUserRow{User: db.User{ID: 2}}, nil
		},
		UnfollowUserFn: func(ctx context.Context, arg db.UnfollowUserParams) (bool, error) {
			return true, nil
		},
	}
	uc := usecase.NewUserUsecase(store, nil, nil)
	err := uc.UnfollowUser(ctx, 1, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListFollowersUsers(t *testing.T) {
	ctx := context.Background()
	store := &MockStore{
		GetUserFn: func(ctx context.Context, arg db.GetUserParams) (db.GetUserRow, error) {
			return db.GetUserRow{User: db.User{ID: 1}}, nil
		},
		ListFollowersUsersFn: func(ctx context.Context, arg db.ListFollowersUsersParams) ([]db.ListFollowersUsersRow, error) {
			return []db.ListFollowersUsersRow{
				{User: db.User{ID: 2, Username: "user2"}, IsFollowing: true},
			}, nil
		},
	}
	uc := usecase.NewUserUsecase(store, nil, nil)
	followers, err := uc.ListFollowers(ctx, 1, 0, 10, ptr[int64](1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(followers) != 1 || followers[0].ID != 2 {
		t.Fatal("expected 1 follower with ID 2")
	}
}
