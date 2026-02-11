package com.fei.twitterbackend.service;

import com.fei.twitterbackend.model.dto.auth.AuthResponse;
import com.fei.twitterbackend.model.dto.user.UserDTO;
import com.fei.twitterbackend.model.entity.RefreshToken;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.model.enums.Role;
import com.fei.twitterbackend.repository.UserRepository;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdToken;
import com.google.api.client.googleapis.auth.oauth2.GoogleIdTokenVerifier;
import com.google.api.client.http.javanet.NetHttpTransport;
import com.google.api.client.json.gson.GsonFactory;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.web.server.ResponseStatusException;

import java.util.Collections;

@Service
@RequiredArgsConstructor
public class AuthService {

    @Value("${spring.security.oauth2.client.registration.google.client-id}")
    private String googleClientId;

    private final UserRepository userRepository;
    private final JwtService jwtService; // (Code provided in previous response)
    private final RefreshTokenService refreshTokenService;

    public AuthResponse loginWithGoogle(String googleIdToken) {
        // Verify Google Token
        GoogleIdToken.Payload payload = verifyGoogleToken(googleIdToken);
        if (payload == null) throw new ResponseStatusException(HttpStatus.BAD_REQUEST,"Invalid Google Token.");

        String email = payload.getEmail();

        // Create or Get User (One-Time Import)
        User user = userRepository.findByEmail(email).orElseGet(() -> {
            User newUser = new User();
            newUser.setEmail(email);
            newUser.setUsername(email.split("@")[0]);
            newUser.setDisplayName((String) payload.get("name"));
            newUser.setAvatarUrl((String) payload.get("picture"));
            newUser.setRole(Role.USER);
            newUser.setProvider("GOOGLE");
            return userRepository.save(newUser);
        });

        // Generate Tokens
        String accessToken = jwtService.generateToken(user);
        RefreshToken refreshToken = refreshTokenService.createRefreshToken(user);

        return new AuthResponse(
                accessToken,
                refreshToken.getToken(),
                UserDTO.fromEntity(user)
        );
    }

    public AuthResponse refreshToken(String requestRefreshToken) {
        return refreshTokenService.findByToken(requestRefreshToken)
                .map(refreshTokenService::verifyExpiration)
                .map(RefreshToken::getUser)
                .map(user -> {
                    String newAccessToken = jwtService.generateToken(user);
                    RefreshToken newRefreshToken = refreshTokenService.createRefreshToken(user);

                    return new AuthResponse(
                            newAccessToken,
                            newRefreshToken.getToken(),
                            UserDTO.fromEntity(user)
                    );
                })
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.UNAUTHORIZED,"Refresh token not found."));
    }

    public void logout(User user) {
        // delete the refresh token from DB
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
            return null;
        }
    }
}