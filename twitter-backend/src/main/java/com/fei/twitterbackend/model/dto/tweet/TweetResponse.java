package com.fei.twitterbackend.model.dto.tweet;

import com.fei.twitterbackend.model.dto.user.UserDTO;
import java.time.LocalDateTime;

public record TweetResponse(
        Long id,
        String content,         // Will be NULL for a Retweet
        String mediaType,
        String mediaUrl,
        UserDTO user,
        int replyCount,
        int likeCount,
        int retweetCount,
        boolean likedByMe,
        boolean retweetedByMe,
        TweetResponse originalTweet, // Null if not a retweet
        LocalDateTime createdAt
) {
    public boolean isRetweet() {
        return originalTweet != null && (content == null || content.isBlank());
    }
}