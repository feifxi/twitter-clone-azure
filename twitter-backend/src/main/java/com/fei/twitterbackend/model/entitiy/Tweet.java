package com.fei.twitterbackend.model.entitiy;

import jakarta.persistence.*;
import lombok.Data;
import org.hibernate.annotations.CreationTimestamp;

import java.time.LocalDateTime;

@Entity
@Table(name = "tweets")
@Data
public class Tweet {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "user_id", nullable = false)
    private User user;

    @Column(length = 280)
    private String content;

    private String mediaUrl;

    // --- REPLY LOGIC (Option A: Flat) ---
    // We only need to know the immediate parent.
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "parent_id")
    private Tweet parentTweet;

    // We don't map the "List<Tweet> replies" here usually,
    // because fetching a tweet shouldn't automatically fetch 10,000 replies.
    // You fetch replies via a separate Repository query.

    // --- RETWEET LOGIC ---
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "retweet_id")
    private Tweet originalTweet;

    // --- COUNTERS ---
    private int replyCount = 0;
    private int retweetCount = 0;
    private int likeCount = 0;

    @CreationTimestamp
    private LocalDateTime createdAt;

    // Helper to check type
    public boolean isRetweet() {
        return originalTweet != null;
    }

    public boolean isReply() {
        return parentTweet != null;
    }
}
