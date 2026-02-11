package com.fei.twitterbackend.model.dto.tweet;

public record TweetRequest(
        String content,
        Long parentId
) {}