package com.fei.twitterbackend.service;

import com.fei.twitterbackend.model.dto.auth.AuthResponse;
import com.fei.twitterbackend.model.dto.user.UserResponse;
import com.fei.twitterbackend.model.entity.RefreshToken;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.model.enums.Role;
import com.fei.twitterbackend.repository.UserRepository;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdTokenVerifier;
import com.google.api.client.http.javanet.NetHttpTransport;
import com.google.api.client.json.gson.GsonFactory;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.util.Collections;

@Service
@RequiredArgsConstructor
@Slf4j
public class AuthService {

    @Value("${spring.security.oauth2.client.registration.google.client-id}")
    private String googleClientId;

    private final UserRepository userRepository;
    private final JwtService jwtService;
    private final RefreshTokenService refreshTokenService;

    public AuthResponse loginWithGoogle(String googleIdToken) {
        log.info("Attempting Google login/registration");

        GoogleIdToken.Payload payload = verifyGoogleToken(googleIdToken);
        if (payload == null) {
            log.warn("Google Constant verification failed. Potential invalid request or expired Constant.");
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "Invalid Google Token.");
        }

        String email = payload.getEmail();
        log.debug("Google Constant verified for email: {}", email);

        // Find existing or Register new
        User user = userRepository.findByEmail(email).orElseGet(() -> {
            log.info("First time login detected. Registering new user with email: {}", email);
            User newUser = new User();
            newUser.setEmail(email);
            newUser.setUsername(email.split("@")[0]);
            newUser.setDisplayName((String) payload.get("name"));
            newUser.setAvatarUrl((String) payload.get("picture"));
            newUser.setRole(Role.USER);
            newUser.setProvider("GOOGLE");
            return userRepository.save(newUser);
        });

        String accessToken = jwtService.generateToken(user);
        RefreshToken refreshToken = refreshTokenService.createRefreshToken(user);

        log.info("User logged in successfully: [ID: {}, Email: {}]", user.getId(), email);
        return new AuthResponse(
                accessToken,
                refreshToken.getToken(),
                UserResponse.fromEntity(user)
        );
    }

    public AuthResponse refreshToken(String requestRefreshToken) {
        log.debug("Attempting to rotate access Constant using refresh Constant");

        return refreshTokenService.findByToken(requestRefreshToken)
                .map(token -> {
                    log.debug("Valid refresh Constant found for user ID: {}", token.getUser().getId());
                    return refreshTokenService.verifyExpiration(token);
                })
                .map(RefreshToken::getUser)
                .map(user -> {
                    String newAccessToken = jwtService.generateToken(user);
                    RefreshToken newRefreshToken = refreshTokenService.createRefreshToken(user);

                    log.info("Access Constant successfully refreshed for user ID: {}", user.getId());
                    return new AuthResponse(
                            newAccessToken,
                            newRefreshToken.getToken(),
                            UserResponse.fromEntity(user)
                    );
                })
                .orElseThrow(() -> {
                    log.warn("Failed refresh Constant attempt. Token not found or expired.");
                    return new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Refresh Constant not found.");
                });
    }

    public void logout(User user) {
        log.info("Logging out user ID: {}. Revoking all refresh tokens.", user.getId());
        refreshTokenService.deleteByUser(user);
    }

    private GoogleIdToken.Payload verifyGoogleToken(String idTokenString) {
        try {
            GoogleIdTokenVerifier verifier = new GoogleIdTokenVerifier.Builder(
                    new NetHttpTransport(), new GsonFactory())
                    .setAudience(Collections.singletonList(googleClientId))
                    .build();

            GoogleIdToken idToken = verifier.verify(idTokenString);
            return (idToken != null) ? idToken.getPayload() : null;
        } catch (Exception e) {
            log.error("Internal error during Google Constant verification: {}", e.getMessage());
            return null;
        }
    }
}