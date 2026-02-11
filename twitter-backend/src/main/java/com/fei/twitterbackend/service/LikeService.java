package com.fei.twitterbackend.service;

import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.TweetLike;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.repository.LikeRepository;
import com.fei.twitterbackend.repository.TweetRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

@Service
@RequiredArgsConstructor
public class LikeService {

    private final LikeRepository likeRepository;
    private final TweetRepository tweetRepository;

    @Transactional
    public void likeTweet(User user, Long tweetId) {
        // 1. Verify Tweet Exists
        Tweet tweet = tweetRepository.findById(tweetId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Tweet not found"));

        // 2. Prevent Duplicate Like
        if (likeRepository.existsByUserIdAndTweetId(user.getId(), tweetId)) {
            throw new ResponseStatusException(HttpStatus.CONFLICT, "You already liked this tweet");
        }

        // 3. Save Like
        TweetLike like = TweetLike.builder()
                .user(user)
                .tweet(tweet)
                .build();
        likeRepository.save(like);

        // 4. Increment Counter
        tweetRepository.incrementLikeCount(tweetId);
    }

    @Transactional
    public void unlikeTweet(User user, Long tweetId) {
        // 1. Check if Like Exists
        if (!likeRepository.existsByUserIdAndTweetId(user.getId(), tweetId)) {
            throw new ResponseStatusException(HttpStatus.BAD_REQUEST, "You have not liked this tweet");
        }

        // 2. Delete Like
        likeRepository.deleteByUserIdAndTweetId(user.getId(), tweetId);

        // 3. Decrement Counter
        tweetRepository.decrementLikeCount(tweetId);
    }
}
