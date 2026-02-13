package com.fei.twitterbackend.service;

import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.model.enums.MediaType;
import com.fei.twitterbackend.repository.TweetRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.server.ResponseStatusException;

import java.util.Optional;

@Service
@RequiredArgsConstructor
@Slf4j
public class RetweetService {

    private final TweetRepository tweetRepository;

    @Transactional
    public void retweet(User user, Long tweetId) {
        log.info("User {} is retweeting tweet {}", user.getId(), tweetId);

        Tweet targetTweet = tweetRepository.findById(tweetId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Tweet not found"));

        // Flatten Logic: Always retweet the original post
        if (targetTweet.getRetweet() != null) {
            targetTweet = targetTweet.getRetweet();
        }

        // Idempotency: If already retweeted, just return
        if (tweetRepository.existsByUserIdAndRetweetId(user.getId(), targetTweet.getId())) {
            log.debug("User {} already retweeted tweet {}. Skipping.", user.getId(), targetTweet.getId());
            return;
        }

        Tweet retweet = Tweet.builder()
                .user(user)
                .retweet(targetTweet)
                .content(null)
                .mediaType(MediaType.NONE)
                .build();

        tweetRepository.save(retweet);
        tweetRepository.incrementRetweetCount(targetTweet.getId());
        log.info("Retweet created for User {} on Tweet {}", user.getId(), targetTweet.getId());
    }

    @Transactional
    public void unretweet(User user, Long tweetId) {
        log.info("User {} is undoing retweet for tweet {}", user.getId(), tweetId);

        Tweet targetTweet = tweetRepository.findById(tweetId)
                .orElseThrow(() -> new ResponseStatusException(HttpStatus.NOT_FOUND, "Tweet not found"));

        // Flatten Logic: Find the original if target is a retweet
        if (targetTweet.getRetweet() != null) {
            targetTweet = targetTweet.getRetweet();
        }

        // Idempotency: If the retweet doesn't exist, just return
        Optional<Tweet> existingRetweet = tweetRepository.findByUserIdAndRetweetId(user.getId(), targetTweet.getId());

        if (existingRetweet.isEmpty()) {
            log.debug("No retweet record found for User {} and Tweet {}. Skipping.", user.getId(), targetTweet.getId());
            return;
        }

        tweetRepository.delete(existingRetweet.get());
        tweetRepository.decrementRetweetCount(targetTweet.getId());
        log.info("Retweet removed for User {} on Tweet {}", user.getId(), targetTweet.getId());
    }
}