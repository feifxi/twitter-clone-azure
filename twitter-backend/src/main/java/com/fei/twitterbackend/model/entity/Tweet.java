package com.fei.twitterbackend.model.entity;

import com.fei.twitterbackend.model.enums.MediaType;
import jakarta.persistence.*;
import lombok.*;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;

@Getter
@Setter
@NoArgsConstructor
@AllArgsConstructor
@Builder
@Entity
@Table(name = "tweets")
public class Tweet {
    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "content", length = 280, nullable = false)
    private String content;

    @Enumerated(EnumType.STRING)
    @Column(name = "media_type")
    private MediaType mediaType; // IMAGE, VIDEO, NONE

    @Column(name = "media_url")
    private String mediaUrl; // For images/videos (Azure Blob URL)

    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "user_id", nullable = false)
    private User user;

    // SELF-REFERENCE: A tweet can be a reply to another tweet
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "parent_id")
    private Tweet parent;

    // A tweet can have many replies (One-to-Many)
    @OneToMany(mappedBy = "parent", cascade = CascadeType.ALL, orphanRemoval = true)
    private List<Tweet> replies = new ArrayList<>();

    // Simple Counters (Optimization: Don't count(*) every time)
    @Builder.Default
    @Column(name = "like_count")
    private int likeCount = 0;

    @Builder.Default
    @Column(name = "retweet_count")
    private int retweetCount = 0;

    @Builder.Default
    @Column(name = "reply_count")
    private int replyCount = 0;

    @CreationTimestamp
    @Column(name = "created_at", nullable = false)
    private LocalDateTime createdAt;

    @UpdateTimestamp
    @Column(name = "updated_at", nullable = false)
    private LocalDateTime updatedAt;
}