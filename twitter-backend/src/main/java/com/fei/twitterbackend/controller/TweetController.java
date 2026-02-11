package com.fei.twitterbackend.controller;

import com.fei.twitterbackend.model.dto.common.PageResponse;
import com.fei.twitterbackend.model.dto.tweet.TweetRequest;
import com.fei.twitterbackend.model.dto.tweet.TweetResponse;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.service.TweetService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/api/v1/tweets")
@RequiredArgsConstructor
public class TweetController {

    private final TweetService tweetService;

    @PostMapping(consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    public ResponseEntity<TweetResponse> createTweet(
            @AuthenticationPrincipal User user,
            @RequestPart("data") @Valid TweetRequest request,
            @RequestPart(value = "file", required = false) MultipartFile file
    ) {
        return ResponseEntity.ok(tweetService.createTweet(user, request, file));
    }

    // Global Feed
    @GetMapping
    public ResponseEntity<PageResponse<TweetResponse>> getGlobalFeed(
            @AuthenticationPrincipal User user, // Can be null if public? Usually not in this app.
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size
    ) {
        Page<TweetResponse> tweetPage = tweetService.getGlobalFeed(user, page, size);
        return ResponseEntity.ok(PageResponse.from(tweetPage));
    }

    // Get Single Tweet
    @GetMapping("/{id}")
    public ResponseEntity<TweetResponse> getTweet(
            @AuthenticationPrincipal User user,
            @PathVariable Long id
    ) {
        return ResponseEntity.ok(tweetService.getTweetById(user, id));
    }

    // Get Replies
    @GetMapping("/{id}/replies")
    public ResponseEntity<PageResponse<TweetResponse>> getReplies(
            @AuthenticationPrincipal User user,
            @PathVariable Long id,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size
    ) {
        Page<TweetResponse> tweetPage = tweetService.getReplies(user, id, page, size);
        return ResponseEntity.ok(PageResponse.from(tweetPage));
    }

    // Get User Profile Tweets
    @GetMapping("/user/{userId}")
    public ResponseEntity<PageResponse<TweetResponse>> getUserTweets(
            @AuthenticationPrincipal User user,
            @PathVariable Long userId,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size
    ) {
        Page<TweetResponse> tweetPage = tweetService.getUserTweets(user, userId, page, size);
        return ResponseEntity.ok(PageResponse.from(tweetPage));
    }

    // UPDATE
    @PutMapping("/{id}")
    public ResponseEntity<TweetResponse> updateTweet(
            @AuthenticationPrincipal User user,
            @PathVariable Long id,
            @RequestBody @Valid TweetRequest request
    ) {
        // We ignore parentId in update logic, only content is used
        return ResponseEntity.ok(tweetService.updateTweet(user, id, request));
    }

    // DELETE
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteTweet(
            @AuthenticationPrincipal User user,
            @PathVariable Long id
    ) {
        tweetService.deleteTweet(user, id);
        return ResponseEntity.noContent().build(); // Returns 204 No Content
    }

    @PostMapping("/{id}/retweet")
    public ResponseEntity<Void> retweet(
            @AuthenticationPrincipal User user,
            @PathVariable Long id
    ) {
        tweetService.retweet(user, id);
        return ResponseEntity.ok().build();
    }
}