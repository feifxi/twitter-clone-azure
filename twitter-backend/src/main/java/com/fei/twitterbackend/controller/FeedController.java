package com.fei.twitterbackend.controller;

import com.fei.twitterbackend.model.dto.common.PageResponse;
import com.fei.twitterbackend.model.dto.tweet.TweetResponse;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.service.FeedService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/feeds")
@RequiredArgsConstructor
public class FeedController {

    private final FeedService feedService;

    @GetMapping
    public ResponseEntity<PageResponse<TweetResponse>> getFeed(
            @AuthenticationPrincipal User user,
            @RequestParam(defaultValue = "foryou") String type,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size
    ) {
        if ("following".equalsIgnoreCase(type)) {
            return ResponseEntity.ok(feedService.getFollowingTimeline(user, page, size));
        }
        return ResponseEntity.ok(feedService.getForYouFeed(user, page, size));
    }

    @GetMapping("/user/{userId}")
    public ResponseEntity<PageResponse<TweetResponse>> getUserProfileFeed(
            @AuthenticationPrincipal User user,
            @PathVariable Long userId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size
    ) {
        PageResponse<TweetResponse> tweetPage = feedService.getUserTweets(user, userId, page, size);
        return ResponseEntity.ok(tweetPage);
    }
}