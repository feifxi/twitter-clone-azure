package com.fei.twitterbackend.controller;

import com.fei.twitterbackend.model.dto.common.PageResponse;
import com.fei.twitterbackend.model.dto.tweet.TweetResponse;
import com.fei.twitterbackend.model.dto.user.UserResponse;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.service.SearchService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/search")
@RequiredArgsConstructor
public class SearchController {

    private final SearchService searchService;

    @GetMapping("/tweets")
    public ResponseEntity<PageResponse<TweetResponse>> searchTweets(
            @AuthenticationPrincipal User user,
            @RequestParam(name = "q") String query,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {

        PageResponse<TweetResponse> results = searchService.searchTweets(user, query, page, size);
        return ResponseEntity.ok(results);
    }

    @GetMapping("/users")
    public ResponseEntity<PageResponse<UserResponse>> searchUsers(
            @AuthenticationPrincipal User user,
            @RequestParam(name = "q") String query,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size) {

        PageResponse<UserResponse> results = searchService.searchUsers(user, query, page, size);
        return ResponseEntity.ok(results);
    }
}