package com.fei.twitterbackend.controller;

import com.fei.twitterbackend.model.dto.common.PageResponse;
import com.fei.twitterbackend.model.dto.hashtag.TrendingHashtagDTO;
import com.fei.twitterbackend.model.dto.tweet.TweetResponse;
import com.fei.twitterbackend.model.dto.user.UserResponse;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.service.DiscoveryService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/v1/discovery")
@RequiredArgsConstructor
public class DiscoveryController {

    private final DiscoveryService discoveryService;

    // Usage: GET /api/v1/discovery/tags?query=java
    @GetMapping("/tags")
    public ResponseEntity<List<TrendingHashtagDTO>> searchHashtags(
            @RequestParam("q") String query,
            @RequestParam(defaultValue = "5") int limit
    ) {
        return ResponseEntity.ok(discoveryService.searchHashtags(query, limit));
    }

    @GetMapping("/trending")
    public ResponseEntity<List<TrendingHashtagDTO>> getTrendingHashtags(
            @RequestParam(defaultValue = "10") int limit) {

        // Capping the limit to prevent malicious requests asking for 100,000 tags
        int safeLimit = Math.min(limit, 50);

        return ResponseEntity.ok(discoveryService.getTrendingHashtags(safeLimit));
    }

    @GetMapping("/hashtags/{hashtag}/tweets")
    public ResponseEntity<PageResponse<TweetResponse>> getTweetsForHashtag(
            @AuthenticationPrincipal User user,
            @PathVariable String hashtag,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {

        PageResponse<TweetResponse> results = discoveryService.getTweetsByHashtag(user, hashtag, page, size);
        return ResponseEntity.ok(results);
    }

    // TIP for Frontend:
    // 1. For the Sidebar Widget: call /api/v1/discovery/users?page=0&size=3
    // 2. For the "Show More" Page: call /api/v1/discovery/users?page=0&size=20
    @GetMapping("/users")
    public ResponseEntity<PageResponse<UserResponse>> getSuggestedUsers(
            @AuthenticationPrincipal User user,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size
    ) {
        return ResponseEntity.ok(discoveryService.getSuggestedUsers(user, page, size));
    }
}