package com.fei.twitterbackend.service;

import com.fei.twitterbackend.mapper.TweetMapper;
import com.fei.twitterbackend.model.dto.common.PageResponse;
import com.fei.twitterbackend.model.dto.tweet.TweetResponse;
import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.User;
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
import org.springframework.web.server.ResponseStatusException;

@Service
@RequiredArgsConstructor
@Slf4j
public class FeedService {

    private final TweetRepository tweetRepository;
    private final TweetMapper tweetMapper;

    // TAB 1: FOR YOU (The Global/Discovery)
    @Transactional(readOnly = true)
    public PageResponse<TweetResponse> getForYouFeed(User currentUser, int page, int size) {
        log.debug("Loading 'For You' feed for user: {}", currentUser != null ? currentUser.getId() : "Guest");
        Pageable pageable = createPageable(page, size, Sort.Direction.DESC);

        Page<Tweet> tweets = tweetRepository.findAllByParentIdIsNull(pageable);
        return tweetMapper.toResponsePage(tweets, currentUser);
    }

    // TAB 2: FOLLOWING (Only people you follow)
    @Transactional(readOnly = true)
    public PageResponse<TweetResponse> getFollowingTimeline(User currentUser, int page, int size) {
        // If not logged in, they can't have a following feed
        if (currentUser == null) {
            throw new ResponseStatusException(HttpStatus.UNAUTHORIZED, "Login to see following feed");
        }
        log.debug("Loading 'Following' timeline for user: {}", currentUser.getId());
        Pageable pageable = createPageable(page, size, Sort.Direction.DESC);

        Page<Tweet> tweets = tweetRepository.findFollowingTimeline(currentUser.getId(), pageable);
        return tweetMapper.toResponsePage(tweets, currentUser);
    }

    @Transactional(readOnly = true)
    public PageResponse<TweetResponse> getUserTweets(User currentUser, Long userId, int page, int size) {
        log.debug("Fetching profile feed for user {}. Page: {}", userId, page);
        Pageable pageable = createPageable(page, size, Sort.Direction.DESC);

        Page<Tweet> tweetsPage = tweetRepository.findAllByUserIdAndParentIdIsNull(userId, pageable);
        return tweetMapper.toResponsePage(tweetsPage, currentUser);
    }

    // HELPER METHODS
    private Pageable createPageable(int page, int size, Sort.Direction direction) {
        return PageRequest.of(page, size, Sort.by(direction, "createdAt"));
    }
}
