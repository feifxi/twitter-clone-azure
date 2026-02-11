package com.fei.twitterbackend.model.dto.tweet;

import com.fei.twitterbackend.model.dto.user.UserDTO;
import com.fei.twitterbackend.model.entity.Tweet;

import java.time.LocalDateTime;

public record TweetResponse(
        Long id,
        String content,
        String mediaType,
        String mediaUrl,
        UserDTO user,
        int replyCount,
        int likeCount,
        int retweetCount,
        boolean likedByMe,
        LocalDateTime createdAt
) {}