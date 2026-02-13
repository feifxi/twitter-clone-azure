package com.fei.twitterbackend.controller;

import com.fei.twitterbackend.model.dto.common.PageResponse;
import com.fei.twitterbackend.model.dto.tweet.TweetRequest;
import com.fei.twitterbackend.model.dto.tweet.TweetResponse;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.service.LikeService;
import com.fei.twitterbackend.service.RetweetService;
import com.fei.twitterbackend.service.TweetService;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
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
    private final LikeService likeService;
    private final RetweetService retweetService;

    @PostMapping(consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    public ResponseEntity<TweetResponse> createTweet(
            @AuthenticationPrincipal User user,
            @RequestPart("data") @Valid TweetRequest request,
            @RequestPart(value = "media", required = false) MultipartFile media
    ) {
        return ResponseEntity.status(HttpStatus.CREATED)
                .body(tweetService.createTweet(user, request, media));
    }

    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteTweet(
            @AuthenticationPrincipal User user,
            @PathVariable Long id
    ) {
        tweetService.deleteTweet(user, id);
        return ResponseEntity.noContent().build();
    }

    @GetMapping("/{id}")
    public ResponseEntity<TweetResponse> getTweet(
            @AuthenticationPrincipal User user,
            @PathVariable Long id
    ) {
        return ResponseEntity.ok(tweetService.getTweetById(user, id));
    }

    @GetMapping("/{id}/replies")
    public ResponseEntity<PageResponse<TweetResponse>> getReplies(
            @AuthenticationPrincipal User user,
            @PathVariable Long id,
            @RequestParam(defaultValue = "0") int page,
            @RequestParam(defaultValue = "20") int size
    ) {
        PageResponse<TweetResponse> tweetPage = tweetService.getReplies(user, id, page, size);
        return ResponseEntity.ok(tweetPage);
    }

    @PostMapping("/{id}/like")
    public ResponseEntity<Void> likeTweet(
            @PathVariable Long id,
            @AuthenticationPrincipal User user
    ) {
        likeService.likeTweet(user, id);
        return ResponseEntity.ok().build();
    }

    @DeleteMapping("/{id}/like")
    public ResponseEntity<Void> unlikeTweet(
            @PathVariable Long id,
            @AuthenticationPrincipal User user
    ) {
        likeService.unlikeTweet(user, id);
        return ResponseEntity.ok().build();
    }

    @PostMapping("/{id}/retweet")
    public ResponseEntity<Void> retweet(
            @AuthenticationPrincipal User user,
            @PathVariable Long id
    ) {
        retweetService.retweet(user, id);
        return ResponseEntity.ok().build();
    }

    @DeleteMapping("/{id}/retweet")
    public ResponseEntity<Void> unretweet(
            @AuthenticationPrincipal User user,
            @PathVariable Long id
    ) {
        retweetService.unretweet(user, id);
        return ResponseEntity.ok().build();
    }
}