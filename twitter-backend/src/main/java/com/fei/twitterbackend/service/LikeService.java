package com.fei.twitterbackend.service;

import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.TweetLike;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.repository.LikeRepository;
import com.fei.twitterbackend.repository.TweetRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

@Service
@RequiredArgsConstructor
@Slf4j
public class LikeService {

    private final LikeRepository likeRepository;
    private final TweetRepository tweetRepository;

    @Transactional
    public void likeTweet(User user, Long tweetId) {
        log.info("User {} is liking tweet {}", user.getId(), tweetId);

        // 1. Check if already liked
        if (likeRepository.existsByUserIdAndTweetId(user.getId(), tweetId)) {
            log.warn("User {} attempted to like tweet {} but already liked it", user.getId(), tweetId);
            throw new ResponseStatusException(HttpStatus.CONFLICT, "You already liked this tweet");
        }

        // 2. Load Tweet
        Tweet tweet = tweetRepository.findById(tweetId)
                .orElseThrow(() -> {
                    log.warn("Like attempt failed: Tweet {} not found", tweetId);
                    return new ResponseStatusException(HttpStatus.NOT_FOUND, "Tweet not found");
                });

        // 3. Save Like
        TweetLike like = TweetLike.builder()
                .id(new TweetLike.TweetLikeId(user.getId(), tweetId))
                .user(user)
                .tweet(tweet)
                .build();

        likeRepository.save(like);
        log.debug("TweetLike entity saved for User: {}, Tweet: {}", user.getId(), tweetId);

        // 4. Increment Counter
        tweetRepository.incrementLikeCount(tweetId);
        log.info("Like count incremented for tweet {}", tweetId);
    }

    @Transactional
    public void unlikeTweet(User user, Long tweetId) {
        log.info("User {} is unliking tweet {}", user.getId(), tweetId);

        // 1. Check if Like Exists
        if (!likeRepository.existsByUserIdAndTweetId(user.getId(), tweetId)) {
            log.warn("User {} attempted to unlike tweet {} but no record found", user.getId(), tweetId);
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "You have not liked this tweet");
        }

        // 2. Delete Like
        likeRepository.deleteByUserIdAndTweetId(user.getId(), tweetId);
        log.debug("TweetLike record deleted for User: {}, Tweet: {}", user.getId(), tweetId);

        // 3. Decrement Counter
        tweetRepository.decrementLikeCount(tweetId);
        log.info("Like count decremented for tweet {}", tweetId);
    }
}