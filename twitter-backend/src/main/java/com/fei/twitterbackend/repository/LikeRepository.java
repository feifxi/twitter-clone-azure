package com.fei.twitterbackend.repository;

import com.fei.twitterbackend.model.entity.TweetLike;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Set;

@Repository
public interface LikeRepository extends JpaRepository<TweetLike, Long> {
    boolean existsByUserIdAndTweetId(Long userId, Long tweetId);
    void deleteByUserIdAndTweetId(Long userId, Long tweetId);

    // OPTIMIZED: Find all Tweet IDs that the user liked from a specific list
    @Query("SELECT tl.tweet.id FROM TweetLike tl WHERE tl.user.id = :userId AND tl.tweet.id IN :tweetIds")
    Set<Long> findLikedTweetIdsByUserId(@Param("userId") Long userId, @Param("tweetIds") List<Long> tweetIds);
}