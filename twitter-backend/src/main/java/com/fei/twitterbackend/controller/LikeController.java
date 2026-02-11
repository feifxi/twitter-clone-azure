package com.fei.twitterbackend.controller;

import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.service.LikeService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/tweets")
@RequiredArgsConstructor
public class LikeController {

    private final LikeService likeService;

    // POST /api/v1/tweets/{id}/like
    @PostMapping("/{id}/like")
    public ResponseEntity<Void> likeTweet(
            @PathVariable Long id,
            @AuthenticationPrincipal User user
    ) {
        likeService.likeTweet(user, id);
        return ResponseEntity.ok().build();
    }

    // DELETE /api/v1/tweets/{id}/like
    @DeleteMapping("/{id}/like")
    public ResponseEntity<Void> unlikeTweet(
            @PathVariable Long id,
            @AuthenticationPrincipal User user
    ) {
        likeService.unlikeTweet(user, id);
        return ResponseEntity.ok().build();
    }
}