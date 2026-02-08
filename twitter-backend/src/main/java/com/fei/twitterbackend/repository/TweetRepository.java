package com.fei.twitterbackend.repository;

import com.fei.twitterbackend.model.entitiy.Tweet;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.util.List;

public interface TweetRepository extends JpaRepository<Tweet, Long> {

    // 1. Get Main Feed (Users I follow)
    @Query("SELECT t FROM Tweet t WHERE t.user.id IN :followingIds ORDER BY t.createdAt DESC")
    Page<Tweet> findFeed(@Param("followingIds") List<Long> followingIds, Pageable pageable);

    // 2. Get Replies for a specific Tweet (Flat list)
    Page<Tweet> findByParentTweetId(Long parentId, Pageable pageable);
}