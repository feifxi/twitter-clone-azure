package com.fei.twitterbackend.service;

import com.fei.twitterbackend.mapper.TweetMapper;
import com.fei.twitterbackend.mapper.UserMapper;
import com.fei.twitterbackend.model.dto.common.PageResponse;
import com.fei.twitterbackend.model.dto.hashtag.TrendingHashtagDTO;
import com.fei.twitterbackend.model.dto.tweet.TweetResponse;
import com.fei.twitterbackend.model.dto.user.UserResponse;
import com.fei.twitterbackend.model.entity.Hashtag;
import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.model.projection.TrendingHashtagProjection;
import com.fei.twitterbackend.repository.HashtagRepository;
import com.fei.twitterbackend.repository.TweetRepository;
import com.fei.twitterbackend.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Collections;
import java.util.List;
import java.util.stream.Collectors;

@Service
@RequiredArgsConstructor
@Slf4j
public class DiscoveryService {

    private final HashtagRepository hashtagRepository;
    private final TweetRepository tweetRepository;
    private final UserRepository userRepository;
    private final TweetMapper tweetMapper;
    private final UserMapper userMapper;

    /**
     * Autocomplete for the "Compose Tweet" box.
     * Query: "jav" -> Returns top 5 tags starting with "jav"
     */
    @Transactional(readOnly = true)
    public List<TrendingHashtagDTO> searchHashtags(String query, int limit) {
        if (query == null || query.isBlank()) {
            return Collections.emptyList();
        }

        // Remove the '#' if the frontend sent it (e.g., "#jav" -> "jav")
        String cleanPrefix = query.replace("#", "").trim();

        Pageable pageable = PageRequest.of(0, limit);
        List<Hashtag> hashtags = hashtagRepository.searchHashtagsByPrefix(cleanPrefix, pageable);

        return hashtags.stream()
                .map(h -> new TrendingHashtagDTO(h.getText(), h.getUsageCount()))
                .collect(Collectors.toList());
    }

    /**
     * Calculates the top trending hashtags over the last 24 hours.
     * Caches the result for performance.
     */
    @Transactional(readOnly = true)
    public List<TrendingHashtagDTO> getTrendingHashtags(int limit) {
        log.info("Calculating trending hashtags limit: {}", limit);

        // 1. Try to get tags from the last 24 hours
        List<TrendingHashtagProjection> projections = hashtagRepository.findTrendingHashtags(limit);

        // 2. FALLBACK: If nothing happened today, get the all-time most popular
        if (projections.isEmpty()) {
            projections = hashtagRepository.findAllTimeTopHashtags(limit);
        }

        return projections.stream()
                .map(proj -> new TrendingHashtagDTO(proj.getText(), proj.getCount()))
                .collect(Collectors.toList());
    }

    @Transactional(readOnly = true)
    public PageResponse<TweetResponse> getTweetsByHashtag(User currentUser, String hashtag, int page, int size) {
        // Clean the input automatically removes '#', '!', spaces, and malicious symbols.
        String cleanHashtag = hashtag.replaceAll("[^a-zA-Z0-9_]", "");

        if (cleanHashtag.isEmpty()) {
            return PageResponse.from(Page.empty());
        }

        PageRequest pageRequest = PageRequest.of(page, size);
        Page<Tweet> tweetPage = tweetRepository.findTweetsByHashtag(cleanHashtag, pageRequest);

        return tweetMapper.toResponsePage(tweetPage, currentUser);
    }

    @Transactional(readOnly = true)
    public PageResponse<UserResponse> getSuggestedUsers(User currentUser, int page, int size) {
        Pageable pageable = PageRequest.of(page, size);
        Page<User> usersPage;

        if (currentUser == null) {
            // Guest? Show global top users
            usersPage = userRepository.findTopUsersGlobally(pageable);
        } else {
            // Logged in? Show personalized suggestions (excluding followed)
            usersPage = userRepository.findSuggestedUsers(currentUser.getId(), pageable);
        }

        // MAPPER OPTIMIZATION:
        // Since the SQL query filtered out people we already follow,
        // This saves us from doing ANY batch-fetching!
        return PageResponse.from(
                usersPage.map(user -> UserResponse.fromEntity(user, false))
        );
    }
}