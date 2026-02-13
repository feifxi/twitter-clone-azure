package com.fei.twitterbackend.repository;

import com.fei.twitterbackend.model.entity.User;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface UserRepository extends JpaRepository<User, Long> {
    Optional<User> findByEmail(String email);

    @Query("""
        SELECT u FROM User u
        WHERE LOWER(u.username) LIKE LOWER(CONCAT('%', :query, '%'))
           OR LOWER(u.displayName) LIKE LOWER(CONCAT('%', :query, '%'))
        ORDER BY u.followersCount DESC
        """)
    Page<User> searchUsers(@Param("query") String query, Pageable pageable);

    // FOLLOWER COUNTS (The person being followed)
    @Modifying
    @Query("UPDATE User u SET u.followersCount = u.followersCount + 1 WHERE u.id = :userId")
    void incrementFollowersCount(@Param("userId") Long userId);

    @Modifying
    @Query("UPDATE User u SET u.followersCount = u.followersCount - 1 WHERE u.id = :userId")
    void decrementFollowersCount(@Param("userId") Long userId);

    // FOLLOWING COUNTS (The person doing the action)
    @Modifying
    @Query("UPDATE User u SET u.followingCount = u.followingCount + 1 WHERE u.id = :userId")
    void incrementFollowingCount(@Param("userId") Long userId);

    @Modifying
    @Query("UPDATE User u SET u.followingCount = u.followingCount - 1 WHERE u.id = :userId")
    void decrementFollowingCount(@Param("userId") Long userId);

    // Fetch people who follow the target user
    @Query("SELECT u FROM User u JOIN Follow f ON u.id = f.id.followerId WHERE f.id.followingId = :targetUserId")
    Page<User> findFollowersByUserId(@Param("targetUserId") Long targetUserId, Pageable pageable);

    // Fetch people the target user is following
    @Query("SELECT u FROM User u JOIN Follow f ON u.id = f.id.followingId WHERE f.id.followerId = :targetUserId")
    Page<User> findFollowingByUserId(@Param("targetUserId") Long targetUserId, Pageable pageable);

    // Logged-in: Suggested users
    @Query(value = """
        SELECT u FROM User u
        WHERE u.id != :currentUserId
        AND NOT EXISTS (
            SELECT 1 FROM Follow f
            WHERE f.follower.id = :currentUserId AND f.following.id = u.id
        )
        ORDER BY u.followersCount DESC
    """,
            countQuery = """
        SELECT count(u) FROM User u
        WHERE u.id != :currentUserId
        AND NOT EXISTS (
            SELECT 1 FROM Follow f
            WHERE f.follower.id = :currentUserId AND f.following.id = u.id
        )
    """)
    Page<User> findSuggestedUsers(@Param("currentUserId") Long currentUserId, Pageable pageable);

    // Guest: Top users globally
    @Query(value = "SELECT u FROM User u ORDER BY u.followersCount DESC",
            countQuery = "SELECT count(u) FROM User u")
    Page<User> findTopUsersGlobally(Pageable pageable);
}