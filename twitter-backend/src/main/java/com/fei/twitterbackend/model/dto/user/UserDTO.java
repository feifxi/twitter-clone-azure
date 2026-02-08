package com.fei.twitterbackend.model.dto.user;

import com.fei.twitterbackend.model.entitiy.User;
import com.fei.twitterbackend.model.enums.Role;

public record UserDTO(
        Long id,
        String username,
        String email,
        String displayName,
        String avatarUrl,
        Role role
) {
    // Static Mapper Method: Converts Entity -> DTO
    public static UserDTO fromEntity(User user) {
        return new UserDTO(
                user.getId(),
                user.getUsername(),
                user.getEmail(),
                user.getDisplayName(),
                user.getAvatarUrl(),
                user.getRole()
        );
    }
}