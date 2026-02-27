package usecase

import (
	"context"
	"database/sql"

	"github.com/chanombude/twitter-go-api/internal/db"
)

type AvatarUpload struct {
	Filename    string
	ContentType string
	Reader      interface {
		Read(p []byte) (n int, err error)
	}
}

type UpdateProfileInput struct {
	Bio         *string
	DisplayName *string
	Avatar      *AvatarUpload
}

func (u *Usecase) GetUser(ctx context.Context, targetUserID int64, viewerID *int64) (db.User, *bool, error) {
	user, err := u.store.GetUser(ctx, targetUserID)
	if err != nil {
		return db.User{}, nil, err
	}

	var following *bool
	if viewerID != nil && *viewerID != targetUserID {
		value, err := u.store.IsFollowing(ctx, db.IsFollowingParams{FollowerID: *viewerID, FollowingID: targetUserID})
		if err == nil {
			following = &value
		}
	}

	return user, following, nil
}

func (u *Usecase) UpdateProfile(ctx context.Context, userID int64, input UpdateProfileInput) (db.User, error) {
	existingUser, err := u.store.GetUser(ctx, userID)
	if err != nil {
		return db.User{}, err
	}

	newAvatar := existingUser.AvatarUrl
	uploadedAvatarURL := ""
	if input.Avatar != nil {
		uploadedAvatarURL, err = u.storage.UploadFile(ctx, input.Avatar.Reader, input.Avatar.Filename, input.Avatar.ContentType)
		if err != nil {
			return db.User{}, err
		}
		newAvatar = sql.NullString{String: uploadedAvatarURL, Valid: true}
	}

	bio := existingUser.Bio
	if input.Bio != nil {
		bio = nullStringFromPtr(input.Bio)
	}

	displayName := existingUser.DisplayName
	if input.DisplayName != nil {
		displayName = nullStringFromPtr(input.DisplayName)
	}

	updatedUser, err := u.store.UpdateUserProfile(ctx, db.UpdateUserProfileParams{
		ID:          userID,
		Bio:         bio,
		DisplayName: displayName,
		AvatarUrl:   newAvatar,
	})
	if err != nil {
		if uploadedAvatarURL != "" {
			_ = u.storage.DeleteFile(ctx, uploadedAvatarURL)
		}
		return db.User{}, err
	}

	if uploadedAvatarURL != "" && existingUser.AvatarUrl.Valid {
		_ = u.storage.DeleteFile(ctx, existingUser.AvatarUrl.String)
	}

	return updatedUser, nil
}

func (u *Usecase) FollowUser(ctx context.Context, followerID, targetUserID int64) (bool, error) {
	targetUser, err := u.store.GetUser(ctx, targetUserID)
	if err != nil {
		return false, err
	}

	inserted, err := u.store.FollowUser(ctx, db.FollowUserParams{FollowerID: followerID, FollowingID: targetUserID})
	if err != nil {
		return false, err
	}

	if inserted {
		_ = u.createAndDispatchNotification(ctx, targetUser.ID, followerID, nil, "FOLLOW")
	}
	return inserted, nil
}

func (u *Usecase) UnfollowUser(ctx context.Context, followerID, targetUserID int64) error {
	_, err := u.store.UnfollowUser(ctx, db.UnfollowUserParams{FollowerID: followerID, FollowingID: targetUserID})
	return err
}

func (u *Usecase) ListFollowers(ctx context.Context, targetUserID int64, page, size int32, viewerID *int64) ([]db.User, map[int64]bool, error) {
	users, err := u.store.ListFollowersUsers(ctx, db.ListFollowersUsersParams{
		FollowingID: targetUserID,
		Limit:       size,
		Offset:      page * size,
	})
	if err != nil {
		return nil, nil, err
	}
	followingMap := make(map[int64]bool)
	if viewerID != nil {
		for _, user := range users {
			if *viewerID == user.ID {
				continue
			}
			isFollowing, err := u.store.IsFollowing(ctx, db.IsFollowingParams{FollowerID: *viewerID, FollowingID: user.ID})
			if err == nil {
				followingMap[user.ID] = isFollowing
			}
		}
	}
	return users, followingMap, nil
}

func (u *Usecase) ListFollowing(ctx context.Context, targetUserID int64, page, size int32, viewerID *int64) ([]db.User, map[int64]bool, error) {
	users, err := u.store.ListFollowingUsers(ctx, db.ListFollowingUsersParams{
		FollowerID: targetUserID,
		Limit:      size,
		Offset:     page * size,
	})
	if err != nil {
		return nil, nil, err
	}
	followingMap := make(map[int64]bool)
	if viewerID != nil {
		for _, user := range users {
			if *viewerID == user.ID {
				continue
			}
			isFollowing, err := u.store.IsFollowing(ctx, db.IsFollowingParams{FollowerID: *viewerID, FollowingID: user.ID})
			if err == nil {
				followingMap[user.ID] = isFollowing
			}
		}
	}
	return users, followingMap, nil
}
