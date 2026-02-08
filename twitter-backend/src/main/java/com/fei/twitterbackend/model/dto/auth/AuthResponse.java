package com.fei.twitterbackend.model.dto.auth;

import com.fasterxml.jackson.annotation.JsonIgnore;
import com.fei.twitterbackend.model.dto.user.UserDTO;

public record AuthResponse(
        String accessToken,
        @JsonIgnore
        String refreshToken,
        UserDTO user
) {}