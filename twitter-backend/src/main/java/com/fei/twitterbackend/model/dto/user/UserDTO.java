package com.fei.twitterbackend.model.dto.user;

import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.model.enums.Role;

public record UserDTO(
        Long id,
        String username,
        String email,
        String displayName,
        String avatarUrl,
        Role role,
        int followersCount,
        int followingCount,
        boolean followedByMe
) {

    // 1. FEED MAPPER: Use when viewing a list of tweets or a profile
    public static UserDTO fromEntity(User user, boolean followedByMe) {
        return new UserDTO(
                user.getId(),
                user.getUsername(),
                user.getEmail(),
                user.getDisplayName(),
                user.getAvatarUrl(),
                user.getRole(),
                user.getFollowersCount(),
                user.getFollowingCount(),
                followedByMe
        );
    }

    // 2. AUTH MAPPER: Use for Login/Register (You can't follow yourself)
    public static UserDTO fromEntity(User user) {
        return fromEntity(user, false);
    }
}