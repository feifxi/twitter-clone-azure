package com.fei.twitterbackend.model.dto.auth;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fei.twitterbackend.model.dto.user.UserResponse;

public record AuthResponse(
        String accessToken,
        @JsonIgnore
        String refreshToken,
        UserResponse user
) {}