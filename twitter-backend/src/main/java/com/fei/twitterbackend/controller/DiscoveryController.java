package com.fei.twitterbackend.controller;

import com.fei.twitterbackend.model.dto.common.PageResponse;
import com.fei.twitterbackend.model.dto.hashtag.TrendingHashtagDTO;
import com.fei.twitterbackend.model.dto.user.UserResponse;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.service.DiscoveryService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import java.util.List;

@RestController
@RequestMapping("/api/v1/discovery")
@RequiredArgsConstructor
public class DiscoveryController {

    private final DiscoveryService discoveryService;

    @GetMapping("/trending")
    public ResponseEntity<List<TrendingHashtagDTO>> getTrendingHashtags(
            @RequestParam(defaultValue = "10") int limit) {

        // Capping the limit to prevent malicious requests asking for 100,000 tags
        int safeLimit = Math.min(limit, 50);

        return ResponseEntity.ok(discoveryService.getTrendingHashtags(safeLimit));
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