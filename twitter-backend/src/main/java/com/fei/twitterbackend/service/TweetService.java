package com.fei.twitterbackend.service;

import com.fei.twitterbackend.model.dto.tweet.TweetRequest;
import com.fei.twitterbackend.model.dto.tweet.TweetResponse;
import com.fei.twitterbackend.model.dto.user.UserDTO;
import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.model.enums.MediaType;
import com.fei.twitterbackend.repository.FollowRepository;
import com.fei.twitterbackend.repository.LikeRepository;
import com.fei.twitterbackend.repository.TweetRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.transaction.support.TransactionSynchronization;
import org.springframework.transaction.support.TransactionSynchronizationManager;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.server.ResponseStatusException;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;

@Service
@RequiredArgsConstructor
@Slf4j
public class TweetService {

    private final TweetRepository tweetRepository;
    private final LikeRepository likeRepository;
    private final FollowRepository followRepository;
    private final FileStorageService fileStorageService;

    // ========================================================================
    // 1. CREATE / UPDATE / DELETE (WRITE OPERATIONS)
    // ========================================================================

    @Transactional
    public TweetResponse createTweet(User user, TweetRequest request, MultipartFile file) {
        log.info("User {} is creating a new tweet", user.getId());

        String mediaUrl = null;
        MediaType mediaType = MediaType.NONE;

        // 1. Handle File Upload
        if (file != null && !file.isEmpty()) {
            String contentType = file.getContentType();
            log.debug("Processing file upload. ContentType: {}", contentType);

            if (contentType == null) {
                throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "File is missing Content-Type");
            }

            if (contentType.startsWith("image/")) {
                mediaType = MediaType.IMAGE;
            } else if (contentType.startsWith("video/")) {
                mediaType = MediaType.VIDEO;
            } else {
                log.warn("Unsupported media upload attempt by User {}: {}", user.getId(), contentType);
                throw new ResponseStatusException(HttpStatus.UNSUPPORTED_MEDIA_TYPE, "Only images and videos are allowed");
            }

            mediaUrl = fileStorageService.uploadFile(file);
            log.debug("File uploaded successfully to: {}", mediaUrl);
        }

        // 2. Handle Parent (Reply Logic)
        Tweet parent = null;
        if (request.parentId() != null) {
            log.debug("Fetching parent tweet ID: {}", request.parentId());
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
                .mediaType(mediaType)
                .replyCount(0)
                .likeCount(0)
                .retweetCount(0)
                .build();

        Tweet savedTweet = tweetRepository.save(tweet);
        log.info("Tweet created successfully with ID: {}", savedTweet.getId());

        return mapToDTO(savedTweet, false, false, false);
    }

    @Transactional
    public TweetResponse updateTweet(User user, Long tweetId, TweetRequest request) {
        log.info("User {} requesting update for tweet {}", user.getId(), tweetId);

        Tweet tweet = tweetRepository.findById(tweetId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Tweet not found"));

        // Authorization
        if (!tweet.getUser().getId().equals(user.getId())) {
            log.warn("Unauthorized update attempt. User {} tried to edit Tweet {} owned by User {}",
                    user.getId(), tweetId, tweet.getUser().getId());
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, "You can only edit your own tweets");
        }

        tweet.setContent(request.content());
        Tweet updatedTweet = tweetRepository.save(tweet);
        log.info("Tweet {} updated successfully", tweetId);

        boolean likedByMe = likeRepository.existsByUserIdAndTweetId(user.getId(), tweetId);
        boolean retweetedByMe = tweetRepository.existsByUserIdAndRetweetId(user.getId(), tweetId);

        return mapToDTO(updatedTweet, likedByMe, retweetedByMe, false);
    }

    @Transactional
    public void deleteTweet(User user, Long tweetId) {
        log.info("User {} requesting deletion of tweet {}", user.getId(), tweetId);

        Tweet tweet = tweetRepository.findById(tweetId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Tweet not found"));

        // Authorization
        if (!tweet.getUser().getId().equals(user.getId())) {
            log.warn("Unauthorized delete attempt by User {} on Tweet {}", user.getId(), tweetId);
            throw new ResponseStatusException(HttpStatus.FORBIDDEN, "You can only delete your own tweets");
        }

        // Data Integrity
        if (tweet.getParent() != null) {
            tweetRepository.decrementReplyCount(tweet.getParent().getId());
        }

        // Harvest Media
        List<String> allMediaToDelete = tweetRepository.findAllMediaUrlsInThread(tweetId);
        log.debug("Found {} media files associated with tweet thread {} to cleanup", allMediaToDelete.size(), tweetId);

        // Delete DB
        tweetRepository.delete(tweet);
        log.info("Tweet {} deleted from database", tweetId);

        // Schedule Storage Cleanup
        if (!allMediaToDelete.isEmpty()) {
            TransactionSynchronizationManager.registerSynchronization(new TransactionSynchronization() {
                @Override
                public void afterCommit() {
                    log.info("Transaction committed. Starting batch delete for {} files.", allMediaToDelete.size());
                    fileStorageService.deleteFiles(allMediaToDelete);
                }
            });
        }
    }

    // ========================================================================
    // 2. READ OPERATIONS (FEEDS & SINGLE)
    // ========================================================================

    @Transactional(readOnly = true)
    public Page<TweetResponse> getGlobalFeed(User currentUser, int page, int size) {
        log.debug("Fetching global feed. Page: {}, Size: {}, User: {}", page, size, currentUser != null ? currentUser.getId() : "Guest");
        Pageable pageable = createPageable(page, size, Sort.Direction.DESC);
        Page<Tweet> tweetsPage = tweetRepository.findAllByParentIdIsNull(pageable);
        return mapPageToResponse(tweetsPage, currentUser);
    }

    @Transactional(readOnly = true)
    public TweetResponse getTweetById(User currentUser, Long tweetId) {
        log.debug("Fetching single tweet details: {}", tweetId);

        Tweet tweet = tweetRepository.findById(tweetId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Tweet not found"));

        boolean likedByMe = false;
        boolean retweetedByMe = false;
        boolean isFollowingAuthor = false;

        if (currentUser != null) {
            Long currentUserId = currentUser.getId();
            Long authorId = tweet.getUser().getId();

            likedByMe = likeRepository.existsByUserIdAndTweetId(currentUser.getId(), tweetId);
            retweetedByMe = tweetRepository.existsByUserIdAndRetweetId(currentUser.getId(), tweetId);

            // Only check follow if looking at someone else's tweet
            if (!currentUserId.equals(authorId)) {
                isFollowingAuthor = followRepository.isFollowing(currentUserId, authorId);
            }
        }

        return mapToDTO(tweet, likedByMe, retweetedByMe, isFollowingAuthor);
    }

    @Transactional(readOnly = true)
    public Page<TweetResponse> getUserTweets(User currentUser, Long userId, int page, int size) {
        log.debug("Fetching profile feed for user {}. Page: {}", userId, page);
        Pageable pageable = createPageable(page, size, Sort.Direction.DESC);
        Page<Tweet> tweetsPage = tweetRepository.findAllByUserIdAndParentIdIsNull(userId, pageable);
        return mapPageToResponse(tweetsPage, currentUser);
    }

    @Transactional(readOnly = true)
    public Page<TweetResponse> getReplies(User currentUser, Long tweetId, int page, int size) {
        log.debug("Fetching replies for tweet {}. Page: {}", tweetId, page);

        if (!tweetRepository.existsById(tweetId)) {
            throw new ResponseStatusException(HttpStatus.NOT_FOUND, "Parent tweet not found");
        }

        Pageable pageable = createPageable(page, size, Sort.Direction.ASC);
        Page<Tweet> repliesPage = tweetRepository.findAllByParentId(tweetId, pageable);
        return mapPageToResponse(repliesPage, currentUser);
    }

    // ========================================================================
    // 3. INTERACTION ACTIONS (RETWEET)
    // ========================================================================

    @Transactional
    public void retweet(User user, Long tweetId) {
        log.info("User {} processing retweet for tweet {}", user.getId(), tweetId);

        Tweet targetTweet = tweetRepository.findById(tweetId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Tweet not found"));

        // Flatten Logic
        if (targetTweet.getRetweet() != null) {
            log.debug("Target is already a retweet. Redirecting to original tweet ID: {}", targetTweet.getRetweet().getId());
            targetTweet = targetTweet.getRetweet();
        }

        // Toggle Logic
        Optional<Tweet> existingRetweet = tweetRepository.findByUserIdAndRetweetId(user.getId(), targetTweet.getId());

        if (existingRetweet.isPresent()) {
            log.info("Retweet exists. Removing (Undoing) retweet.");
            tweetRepository.delete(existingRetweet.get());
            tweetRepository.decrementRetweetCount(targetTweet.getId());
        } else {
            log.info("Retweet does not exist. Creating new retweet.");
            Tweet retweet = Tweet.builder()
                    .user(user)
                    .retweet(targetTweet)
                    .content(null)
                    .mediaType(MediaType.NONE)
                    .build();

            tweetRepository.save(retweet);
            tweetRepository.incrementRetweetCount(targetTweet.getId());
        }
    }

    // ========================================================================
    // 4. HELPER METHODS
    // ========================================================================

    private Pageable createPageable(int page, int size, Sort.Direction direction) {
        return PageRequest.of(page, size, Sort.by(direction, "createdAt"));
    }

    private Page<TweetResponse> mapPageToResponse(Page<Tweet> tweetsPage, User currentUser) {
        // 1. Get all Tweet IDs
        List<Long> tweetIds = tweetsPage.getContent().stream().map(Tweet::getId).toList();

        // 2. Get all Author IDs (to check follows)
        List<Long> authorIds = tweetsPage.getContent().stream()
                .map(t -> t.getUser().getId())
                .distinct()
                .toList();

        if (tweetIds.isEmpty()) return tweetsPage.map(t -> mapToDTO(t, false, false, false));

        // 3. Batch Fetch Data
        Set<Long> likedTweetIds;
        Set<Long> retweetedTweetIds;
        Set<Long> followedAuthorIds; // <--- NEW SET

        if (currentUser != null) {
            likedTweetIds = likeRepository.findLikedTweetIdsByUserId(currentUser.getId(), tweetIds);
            retweetedTweetIds = tweetRepository.findRetweetedTweetIdsByUserId(currentUser.getId(), tweetIds);
            // Fetch Follows
            followedAuthorIds = followRepository.findFollowedUserIds(currentUser.getId(), authorIds);
        } else {
            likedTweetIds = Collections.emptySet();
            retweetedTweetIds = Collections.emptySet();
            followedAuthorIds = Collections.emptySet();
        }

        // 4. Map
        return tweetsPage.map(tweet -> {
            boolean isLiked = likedTweetIds.contains(tweet.getId());
            boolean isRetweeted = retweetedTweetIds.contains(tweet.getId());
            boolean isFollowingAuthor = followedAuthorIds.contains(tweet.getUser().getId()); // <--- CHECK

            return mapToDTO(tweet, isLiked, isRetweeted, isFollowingAuthor);
        });
    }

    private TweetResponse mapToDTO(Tweet tweet, boolean likedByMe, boolean retweetedByMe, boolean isFollowingAuthor) {

        TweetResponse originalTweetDTO = null;
        if (tweet.getRetweet() != null) {
            // Recursion: For the nested tweet, we pass false for everything to be safe/fast
            originalTweetDTO = mapToDTO(tweet.getRetweet(), false, false, false);
        }

        return new TweetResponse(
                tweet.getId(),
                tweet.getContent(),
                tweet.getMediaType() != null ? tweet.getMediaType().name() : null,
                tweet.getMediaUrl(),
                UserDTO.fromEntity(tweet.getUser(), isFollowingAuthor),
                tweet.getReplyCount(),
                tweet.getLikeCount(),
                tweet.getRetweetCount(),
                likedByMe,
                retweetedByMe,
                originalTweetDTO,
                tweet.getCreatedAt()
        );
    }
}