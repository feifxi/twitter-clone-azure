package com.fei.twitterbackend.repository;

import com.fei.twitterbackend.model.entity.Tweet;
import jakarta.transaction.Transactional;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;

public interface TweetRepository extends JpaRepository<Tweet, Long> {
    // Global Feed (Only root tweets, no replies)
    Page<Tweet> findAllByParentIdIsNull(Pageable pageable);

    // Main Profile Feed
    Page<Tweet> findAllByUserIdAndParentIdIsNull(Long userId, Pageable pageable);

    // Tweet Replies (Flat strategy)
    Page<Tweet> findAllByParentId(Long parentId, Pageable pageable);

    // Fast Counter Update, fire a raw SQL update. It's atomic and fast.
    @Modifying
    @Transactional
    @Query("UPDATE Tweet t SET t.replyCount = t.replyCount + 1 WHERE t.id = :tweetId")
    void incrementReplyCount(Long tweetId);

    @Modifying
    @Transactional
    @Query("UPDATE Tweet t SET t.replyCount = t.replyCount - 1 WHERE t.id = :tweetId")
    void decrementReplyCount(Long tweetId);

    @Modifying
    @Transactional
    @Query("UPDATE Tweet t SET t.likeCount = t.likeCount + 1 WHERE t.id = :tweetId")
    void incrementLikeCount(Long tweetId);

    @Modifying
    @Transactional
    @Query("UPDATE Tweet t SET t.likeCount = t.likeCount - 1 WHERE t.id = :tweetId")
    void decrementLikeCount(Long tweetId);

    // Efficiently finds the tweet and ALL its descendants' media URLs
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