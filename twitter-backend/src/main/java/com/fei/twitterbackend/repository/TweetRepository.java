package com.fei.twitterbackend.repository;

import com.fei.twitterbackend.model.entity.Tweet;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;
import java.util.Optional;
import java.util.Set;

public interface TweetRepository extends JpaRepository<Tweet, Long> {

    // ========================================================================
    // 1. CORE FEEDS (READS)
    // ========================================================================

    // Global Feed (Only root tweets, no replies)
    Page<Tweet> findAllByParentIdIsNull(Pageable pageable);

    // Main Profile Feed (User's tweets + retweets)
    Page<Tweet> findAllByUserIdAndParentIdIsNull(Long userId, Pageable pageable);

    // Reply Thread (Flat strategy)
    Page<Tweet> findAllByParentId(Long parentId, Pageable pageable);


    // ========================================================================
    // 2. ATOMIC COUNTERS (WRITES)
    // Uses direct SQL updates for performance (Avoids loading entity -> modifying -> saving)
    // ========================================================================

    // --- REPLY COUNTERS ---
    @Modifying
    @Transactional
    @Query("UPDATE Tweet t SET t.replyCount = t.replyCount + 1 WHERE t.id = :tweetId")
    void incrementReplyCount(@Param("tweetId") Long tweetId);

    @Modifying
    @Transactional
    @Query("UPDATE Tweet t SET t.replyCount = t.replyCount - 1 WHERE t.id = :tweetId")
    void decrementReplyCount(@Param("tweetId") Long tweetId);

    // --- LIKE COUNTERS ---
    @Modifying
    @Transactional
    @Query("UPDATE Tweet t SET t.likeCount = t.likeCount + 1 WHERE t.id = :tweetId")
    void incrementLikeCount(@Param("tweetId") Long tweetId);

    @Modifying
    @Transactional
    @Query("UPDATE Tweet t SET t.likeCount = t.likeCount - 1 WHERE t.id = :tweetId")
    void decrementLikeCount(@Param("tweetId") Long tweetId);

    // --- RETWEET COUNTERS ---
    @Modifying
    @Transactional
    @Query("UPDATE Tweet t SET t.retweetCount = t.retweetCount + 1 WHERE t.id = :tweetId")
    void incrementRetweetCount(@Param("tweetId") Long tweetId);

    @Modifying
    @Transactional
    @Query("UPDATE Tweet t SET t.retweetCount = t.retweetCount - 1 WHERE t.id = :tweetId")
    void decrementRetweetCount(@Param("tweetId") Long tweetId);


    // ========================================================================
    // 3. RETWEET LOGIC
    // ========================================================================

    // Check if user has already retweeted a specific tweet (Fast)
    boolean existsByUserIdAndRetweetId(Long userId, Long retweetId);

    // Find specific retweet entity (Used for "Un-retweet" to delete the row)
    Optional<Tweet> findByUserIdAndRetweetId(Long userId, Long retweetId);

    // Batch Fetch: Finds which of the given tweetIds were retweeted by the user
    // Returns: A Set of IDs that should have the "Green Retweet Button" active
    @Query("SELECT t.retweet.id FROM Tweet t WHERE t.user.id = :userId AND t.retweet.id IN :tweetIds")
    Set<Long> findRetweetedTweetIdsByUserId(@Param("userId") Long userId, @Param("tweetIds") List<Long> tweetIds);


    // ========================================================================
    // 4. UTILITIES & COMPLEX QUERIES
    // ========================================================================

    // Recursive Cleanup: Finds the tweet and ALL descendants' media URLs
    // Uses CTE (Common Table Expression) for tree traversal
    @Query(value = """
        WITH RECURSIVE tweet_tree AS (
            SELECT id, media_url, parent_id
            FROM tweets
            WHERE id = :tweetId
            UNION ALL
            SELECT t.id, t.media_url, t.parent_id
            FROM tweets t
            INNER JOIN tweet_tree tt ON t.parent_id = tt.id
        )
        SELECT media_url FROM tweet_tree WHERE media_url IS NOT NULL
    """, nativeQuery = true)
    List<String> findAllMediaUrlsInThread(@Param("tweetId") Long tweetId);
}