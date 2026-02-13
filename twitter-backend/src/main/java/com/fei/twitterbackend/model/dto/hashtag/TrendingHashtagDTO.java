package com.fei.twitterbackend.model.dto.hashtag;

public record TrendingHashtagDTO(
        String hashtag,
        int recentCount
) {}
