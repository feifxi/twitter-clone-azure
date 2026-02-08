package com.fei.twitterbackend.model.dto.common;

public record ApiResponse(
        Boolean success,
        String message
) {}