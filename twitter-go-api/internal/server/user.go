package server

import (
	"github.com/chanombude/twitter-go-api/internal/usecase"
)

type userResponse struct {
	ID             int64   `json:"id"`
	Username       string  `json:"username"`
	Email          string  `json:"email"`
	DisplayName    *string `json:"display_name"`
	Bio            *string `json:"bio"`
	AvatarUrl      *string `json:"avatar_url"`
	IsFollowing    bool    `json:"is_following"`
	FollowersCount int32   `json:"followers_count"`
	FollowingCount int32   `json:"following_count"`
}

func newUserResponse(user usecase.UserItem) userResponse {
	var displayName, bio, avatarUrl *string
	if user.DisplayName.Valid {
		displayName = &user.DisplayName.String
	}
	if user.Bio.Valid {
		bio = &user.Bio.String
	}
	if user.AvatarUrl.Valid {
		avatarUrl = &user.AvatarUrl.String
	}

	return userResponse{
		ID:             user.ID,
		Username:       user.Username,
		Email:          user.Email,
		DisplayName:    displayName,
		Bio:            bio,
		AvatarUrl:      avatarUrl,
		IsFollowing:    user.IsFollowing,
		FollowersCount: user.FollowersCount,
		FollowingCount: user.FollowingCount,
	}
}
