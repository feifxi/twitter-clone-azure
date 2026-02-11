package com.fei.twitterbackend.service;

import com.fei.twitterbackend.model.dto.tweet.TweetRequest;
import com.fei.twitterbackend.model.dto.tweet.TweetResponse;
import com.fei.twitterbackend.model.dto.user.UserDTO;
import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.model.enums.MediaType;
import com.fei.twitterbackend.repository.LikeRepository;
import com.fei.twitterbackend.repository.TweetRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.transaction.support.TransactionSynchronization;
import org.springframework.transaction.support.TransactionSynchronizationManager;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.http.HttpStatus;

import java.util.Collections;
import java.util.List;
import java.util.Set;

@Service
@RequiredArgsConstructor
public class TweetService {

    private final TweetRepository tweetRepository;
    private final FileStorageService fileStorageService;
    private final LikeRepository likeRepository;

    @Transactional
    public TweetResponse createTweet(User user, TweetRequest request, MultipartFile file) {

        String mediaUrl = null;
        MediaType mediaType = MediaType.NONE;

        // 1. Handle File Upload (If present)
        if (file != null && !file.isEmpty()) {
            String contentType = file.getContentType();

            if (contentType == null) {
                throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "File is missing Content-Type");
            }

            // Strict Checking (The "Allowlist" approach)
            if (contentType.startsWith("image/")) {
                mediaType = MediaType.IMAGE;
            }
            else if (contentType.startsWith("video/")) {
                mediaType = MediaType.VIDEO;
            }
            else {
                // Reject unsupported files (PDFs, EXEs, Text files)
                throw new ResponseStatusException(HttpStatus.UNSUPPORTED_MEDIA_TYPE, "Only images and videos are allowed. Received: " + contentType);
            }

            //  upload if validation passed
            mediaUrl = fileStorageService.uploadFile(file);
        }

        // 2. Handle Parent (Reply Logic)
        Tweet parent = null;
        if (request.parentId() != null) {
            parent = tweetRepository.findById(request.parentId())
                    .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Parent tweet not found"));
            tweetRepository.incrementReplyCount(parent.getId());
        }

        // 3. Save Tweet
        Tweet tweet = Tweet.builder()
                .content(request.content())
                .user(user)
                .parent(parent)
                .mediaUrl(mediaUrl)
                .mediaType(mediaType) // Store the type!
                .replyCount(0)
                .likeCount(0)
                .retweetCount(0)
                .build();

        Tweet savedTweet = tweetRepository.save(tweet);
        return mapToDTO(savedTweet, false);
    }

    // GET FEED (Global)
    @Transactional(readOnly = true)
    public Page<TweetResponse> getGlobalFeed(User currentUser, int page, int size) {
        Pageable pageable = createPageable(page, size, Sort.Direction.DESC);

        // A. Fetch Tweets
        Page<Tweet> tweetsPage = tweetRepository.findAllByParentIdIsNull(pageable);

        // B. Map to DTOs with "likedByMe" logic
        return mapPageToResponse(tweetsPage, currentUser);
    }

    // GET SINGLE TWEET (Detail View)
    @Transactional(readOnly = true)
    public TweetResponse getTweetById(User currentUser, Long tweetId) {
        Tweet tweet = tweetRepository.findById(tweetId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Tweet not found"));

        // Check if liked (Single check is fine here)
        boolean likedByMe = currentUser != null &&
                likeRepository.existsByUserIdAndTweetId(currentUser.getId(), tweetId);

        return mapToDTO(tweet, likedByMe);
    }

    // GET USER TWEETS (Profile)
    @Transactional(readOnly = true)
    public Page<TweetResponse> getUserTweets(User currentUser, Long userId, int page, int size) {
        Pageable pageable = createPageable(page, size, Sort.Direction.DESC);
        Page<Tweet> tweetsPage = tweetRepository.findAllByUserIdAndParentIdIsNull(userId, pageable);
        return mapPageToResponse(tweetsPage, currentUser);
    }

    // GET REPLIES (Flat)
    @Transactional(readOnly = true)
    public Page<TweetResponse> getReplies(User currentUser, Long tweetId, int page, int size) {
        // Validation
        if (!tweetRepository.existsById(tweetId)) {
            throw new ResponseStatusException(HttpStatus.NOT_FOUND, "Parent tweet not found");
        }

        Pageable pageable = createPageable(page, size, Sort.Direction.ASC);
        Page<Tweet> repliesPage = tweetRepository.findAllByParentId(tweetId, pageable);

        return mapPageToResponse(repliesPage, currentUser);
    }

    // UPDATE TWEET
    @Transactional
    public TweetResponse updateTweet(User user, Long tweetId, TweetRequest request) {
        Tweet tweet = tweetRepository.findById(tweetId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Tweet not found"));

        // 1. Authorization: Only the owner can edit
        if (!tweet.getUser().getId().equals(user.getId())) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, "You can only edit your own tweets");
        }

        // 2. Update Content
        tweet.setContent(request.content());
        // Note: @UpdateTimestamp in Entity handles the 'updatedAt' field automatically

        Tweet updatedTweet = tweetRepository.save(tweet);

        // 3. Return DTO (Need to check if *I* liked my own tweet to set the boolean correctly)
        boolean likedByMe = likeRepository.existsByUserIdAndTweetId(user.getId(), tweetId);
        return mapToDTO(updatedTweet, likedByMe);
    }

    // DELETE TWEET
    @Transactional
    public void deleteTweet(User user, Long tweetId) {
        Tweet tweet = tweetRepository.findById(tweetId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Tweet not found"));

        // 1. Authorization
        if (!tweet.getUser().getId().equals(user.getId())) {
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, "You can only delete your own tweets");
        }

        // 2. Data Integrity: Decrement parent reply count
        if (tweet.getParent() != null) {
            tweetRepository.decrementReplyCount(tweet.getParent().getId());
        }

        // 3. Harvest: Get ALL media URLs (parent + all children/grandchildren)
        // This is fast because it returns List<String>, not Entities.
        List<String> allMediaToDelete = tweetRepository.findAllMediaUrlsInThread(tweetId);

        // 4. Delete from DB
        // CascadeType.ALL handles the rows, but we needed the URLs first.
        tweetRepository.delete(tweet);

        // 5. Schedule Azure Cleanup (Best Practice)
        // We register a hook to delete files ONLY if the DB transaction commits successfully.
        // This prevents "Broken Images" (images deleted but tweet reappears due to rollback).
        if (!allMediaToDelete.isEmpty()) {
            TransactionSynchronizationManager.registerSynchronization(new TransactionSynchronization() {
                @Override
                public void afterCommit() {
                    // 1 HTTP request (or extremely few) instead of N requests
                    fileStorageService.deleteFiles(allMediaToDelete);
                }
            });
        }
    }

    private Pageable createPageable(int page, int size, Sort.Direction direction) {
        return PageRequest.of(page, size, Sort.by(direction, "createdAt"));
    }

    // HELPER: Batch Mapping to Avoid N+1 Problem
    private Page<TweetResponse> mapPageToResponse(Page<Tweet> tweetsPage, User currentUser) {
        // 1. Extract all Tweet IDs from the page
        List<Long> tweetIds = tweetsPage.getContent().stream()
                .map(Tweet::getId)
                .toList();

        // 2. Optimization: If list is empty, return empty
        if (tweetIds.isEmpty()) {
            return tweetsPage.map(t -> mapToDTO(t, false));
        }

        // 3. Fetch ONLY the likes for these specific tweets by this user
        Set<Long> likedTweetIds;
        if (currentUser != null) {
            likedTweetIds = likeRepository.findLikedTweetIdsByUserId(currentUser.getId(), tweetIds);
        } else {
            likedTweetIds = Collections.emptySet();
        }

        // 4. Map the page, checking against the Set
        return tweetsPage.map(tweet -> {
            boolean isLiked = likedTweetIds.contains(tweet.getId());
            return mapToDTO(tweet, isLiked);
        });
    }

    private TweetResponse mapToDTO(Tweet tweet, boolean likedByMe) {
        return new TweetResponse(
                tweet.getId(),
                tweet.getContent(),
                tweet.getMediaType().name(),
                tweet.getMediaUrl(),
                UserDTO.fromEntity(tweet.getUser()),
                tweet.getReplyCount(),
                tweet.getLikeCount(),
                tweet.getRetweetCount(),
                likedByMe,
                tweet.getCreatedAt()
        );
    }
}