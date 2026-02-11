package com.fei.twitterbackend.service;

import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.repository.FollowRepository;
import com.fei.twitterbackend.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

@Service
@RequiredArgsConstructor
@Slf4j
public class UserService {

    private final UserRepository userRepository;
    private final FollowRepository followRepository;

    @Transactional
    public void followUser(User currentUser, Long targetUserId) {
        // 1. Validation
        if (currentUser.getId().equals(targetUserId)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "You cannot follow yourself");
        }

        // Check if target user exists
        if (!userRepository.existsById(targetUserId)) {
            throw new ResponseStatusException(HttpStatus.NOT_FOUND, "User not found");
        }

        // 2. Toggle Logic
        boolean isAlreadyFollowing = followRepository.isFollowing(currentUser.getId(), targetUserId);

        if (isAlreadyFollowing) {
            // UNFOLLOW FLOW
            log.info("User {} is unfollowing User {}", currentUser.getId(), targetUserId);

            followRepository.unfollowUser(currentUser.getId(), targetUserId);

            // Update Counters
            userRepository.decrementFollowingCount(currentUser.getId());
            userRepository.decrementFollowersCount(targetUserId);

        } else {
            // FOLLOW FLOW
            log.info("User {} is following User {}", currentUser.getId(), targetUserId);

            followRepository.followUser(currentUser.getId(), targetUserId);

            // Update Counters
            userRepository.incrementFollowingCount(currentUser.getId());
            userRepository.incrementFollowersCount(targetUserId);
        }
    }
}