package com.fei.twitterbackend.model.dto.auth;

import jakarta.validation.constraints.NotBlank;

public record GoogleAuthRequest(
        @NotBlank
        String token
) {}
