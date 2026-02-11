package com.fei.twitterbackend.controller;

import com.fei.twitterbackend.model.dto.auth.AuthResponse;
import com.fei.twitterbackend.model.dto.auth.GoogleAuthRequest;
import com.fei.twitterbackend.model.dto.common.ApiResponse;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.service.AuthService;
import jakarta.servlet.http.HttpServletResponse;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.env.Environment;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseCookie;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.server.ResponseStatusException;

@RestController
@RequestMapping("/api/v1/auth")
@RequiredArgsConstructor
public class AuthController {

    @Value("${jwt.refresh-expiration}")
    private long refreshExpirationMs;

    private final Environment environment;
    private final AuthService authService;

    // 1. LOGIN (Google -> Access + Refresh)
    @PostMapping("/google")
    public ResponseEntity<AuthResponse> googleLogin(
            @Valid @RequestBody GoogleAuthRequest authRequest,
            HttpServletResponse response
    ) {
        AuthResponse authToken = authService.loginWithGoogle(authRequest.token());
        // Set HttpOnly Cookie
        setRefreshTokenCookie(response, authToken.refreshToken());
        return ResponseEntity.ok(authToken);
    }

    // 2. REFRESH (Old Refresh -> New Access + New Refresh)
    @PostMapping("/refresh")
    public ResponseEntity<AuthResponse> refreshToken(
            @CookieValue(name = "refreshToken") String refreshToken,
            HttpServletResponse response
    ) {
        if (refreshToken == null || refreshToken.isEmpty()) {
            throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "No refresh token provided");
        }
        AuthResponse authToken = authService.refreshToken(refreshToken);
        // Rotate Cookie
        setRefreshTokenCookie(response, authToken.refreshToken());
        return ResponseEntity.ok(authToken);
    }

    // 3. LOGOUT (Revoke Refresh Token)
    @PostMapping("/logout")
    public ResponseEntity<ApiResponse> logout(@AuthenticationPrincipal User user, HttpServletResponse response) {
        if (user == null) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED).body(new ApiResponse(false, "No user is currently logged in."));
        }
        authService.logout(user);
        // clear refresh token cookie
        clearRefreshTokenCookie(response);
        return ResponseEntity.ok(new ApiResponse(true ,"Logged out"));
    }

    private void setRefreshTokenCookie(HttpServletResponse response, String token) {
        int maxAgeSeconds = (int) (refreshExpirationMs / 1000);
        ResponseCookie cookie = ResponseCookie.from("refreshToken", token)
                .httpOnly(true)
                .secure(isProduction()) // Dynamic check
                .path("/")
                .maxAge(maxAgeSeconds)
                .sameSite("Lax")
                .build();

        response.addHeader("Set-Cookie", cookie.toString());
    }

    private void clearRefreshTokenCookie(HttpServletResponse response) {
        ResponseCookie cookie = ResponseCookie.from("refreshToken", "")
                .httpOnly(true)
                .secure(isProduction()) // Dynamic check
                .path("/")
                .maxAge(0)
                .sameSite("Lax")
                .build();

        response.addHeader("Set-Cookie", cookie.toString());
    }

    private boolean isProduction() {
        String[] activeProfiles = environment.getActiveProfiles();
        // If "dev" profile is active, we are NOT in production
        return activeProfiles.length == 0 || !activeProfiles[0].equals("dev");
    }
}