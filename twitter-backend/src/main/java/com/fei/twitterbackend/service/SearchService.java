package com.fei.twitterbackend.service;

import com.fei.twitterbackend.mapper.TweetMapper;
import com.fei.twitterbackend.mapper.UserMapper;
import com.fei.twitterbackend.model.dto.common.PageResponse;
import com.fei.twitterbackend.model.dto.tweet.TweetResponse;
import com.fei.twitterbackend.model.dto.user.UserResponse;
import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.repository.TweetRepository;
import com.fei.twitterbackend.repository.UserRepository;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
@RequiredArgsConstructor
@Slf4j
public class SearchService {

    private final TweetRepository tweetRepository;
    private final UserRepository userRepository;
    private final TweetMapper tweetMapper;
    private final UserMapper userMapper;

    @Transactional(readOnly = true)
    public PageResponse<TweetResponse> searchTweets(User currentUser, String rawQuery, int page, int size) {
        log.info("User {} searching for: {}", currentUser != null ? currentUser.getId() : "Guest", rawQuery);

        // Sanitize and format the query
        String sanitizedQuery = prepareTsQuery(rawQuery);
        if (sanitizedQuery.isEmpty()) {
            return PageResponse.from(Page.empty());
        }

        // Execute Native PostgreSQL FTS Query
        PageRequest pageRequest = PageRequest.of(page, size);
        Page<Tweet> tweetPage = tweetRepository.searchTweets(sanitizedQuery, pageRequest);

        return tweetMapper.toResponsePage(tweetPage, currentUser);
    }

    @Transactional(readOnly = true)
    public PageResponse<UserResponse> searchUsers(User currentUser, String rawQuery, int page, int size) {
        log.info("User {} searching for people: {}", currentUser != null ? currentUser.getId() : "Guest", rawQuery);

        if (rawQuery == null || rawQuery.trim().isEmpty()) {
            return PageResponse.from(Page.empty());
        }

        // Clean the query (Just trim it for the LIKE statement)
        String cleanQuery = rawQuery.trim();

        // Execute DB Search
        PageRequest pageRequest = PageRequest.of(page, size);
        Page<User> userPage = userRepository.searchUsers(cleanQuery, pageRequest);

        return userMapper.toResponsePage(userPage, currentUser);
    }

    /**
     * Converts a raw user search like "java spring" into PostgreSQL format "java & spring".
     * Strips dangerous characters to prevent SQL syntax errors.
     */
    private String prepareTsQuery(String query) {
        if (query == null || query.isBlank()) return "";

        // Remove characters that aren't alphanumeric or spaces
        String clean = query.replaceAll("[^a-zA-Z0-9\\s]", "").trim();
        if (clean.isEmpty()) return "";

        // Join words with the '&' operator so PostgreSQL knows to search for ALL words
        return String.join(" & ", clean.split("\\s+"));
    }
}