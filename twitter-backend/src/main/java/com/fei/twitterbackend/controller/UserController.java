package com.fei.twitterbackend.controller;

import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.service.UserService;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/users")
@RequiredArgsConstructor
public class UserController {

    private final UserService userService;

    // Follow / Unfollow Toggle Endpoint
    @PostMapping("/{id}/follow")
    public ResponseEntity<Void> followUser(
            @AuthenticationPrincipal User user,
            @PathVariable Long id
    ) {
        userService.followUser(user, id);
        return ResponseEntity.ok().build();
    }
}