package com.fei.twitterbackend.mapper;

import com.fei.twitterbackend.model.dto.common.PageResponse;
import com.fei.twitterbackend.model.dto.tweet.TweetResponse;
import com.fei.twitterbackend.model.dto.user.UserResponse;
import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.User;
import com.fei.twitterbackend.repository.FollowRepository;
import com.fei.twitterbackend.repository.LikeRepository;
import com.fei.twitterbackend.repository.TweetRepository;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.stereotype.Component;

import java.util.Collections;
import java.util.List;
import java.util.Set;

@Component
@RequiredArgsConstructor
public class TweetMapper {

    private final LikeRepository likeRepository;
    private final TweetRepository tweetRepository;
    private final FollowRepository followRepository;

    // Maps a single Tweet to a TweetResponse.
    public TweetResponse toResponse(Tweet tweet, boolean likedByMe, boolean retweetedByMe, boolean isFollowingAuthor) {
        TweetResponse originalTweetDTO = null;
        if (tweet.getRetweet() != null) {
            // For nested tweets, we usually don't need deep interaction states
            originalTweetDTO = toResponse(tweet.getRetweet(), false, false, false);
        }

        return new TweetResponse(
                tweet.getId(),
                tweet.getContent(),
                tweet.getMediaType() != null ? tweet.getMediaType().name() : null,
                tweet.getMediaUrl(),
                UserResponse.fromEntity(tweet.getUser(), isFollowingAuthor),
                tweet.getReplyCount(),
                tweet.getLikeCount(),
                tweet.getRetweetCount(),
                likedByMe,
                retweetedByMe,
                originalTweetDTO,
                tweet.getCreatedAt()
        );
    }

    // Maps a Page of Tweets to a Page of TweetResponses with batch-fetched interaction states.
    public PageResponse<TweetResponse> toResponsePage(Page<Tweet> tweetsPage, User currentUser) {
        List<Tweet> tweets = tweetsPage.getContent();

        // 1. Handle Empty Case
        if (tweets.isEmpty()) {
            Page<TweetResponse> emptyPage = tweetsPage.map(t -> toResponse(t, false, false, false));
            return PageResponse.from(emptyPage); // Wrap and return
        }

        // 2. Extract IDs
        List<Long> tweetIds = tweets.stream().map(Tweet::getId).toList();
        List<Long> authorIds = tweets.stream().map(t -> t.getUser().getId()).distinct().toList();

        // 3. Batch Fetch Data
        Set<Long> likedTweetIds;
        Set<Long> retweetedTweetIds;
        Set<Long> followedAuthorIds;

        if (currentUser != null) {
            likedTweetIds = likeRepository.findLikedTweetIdsByUserId(currentUser.getId(), tweetIds);
            retweetedTweetIds = tweetRepository.findRetweetedTweetIdsByUserId(currentUser.getId(), tweetIds);
            followedAuthorIds = followRepository.findFollowedUserIds(currentUser.getId(), authorIds);
        } else {
            likedTweetIds = Collections.emptySet();
            retweetedTweetIds = Collections.emptySet();
            followedAuthorIds = Collections.emptySet();
        }

        // 4. Map the Page
        Page<TweetResponse> mappedPage = tweetsPage.map(tweet -> {
            boolean isLiked = likedTweetIds.contains(tweet.getId());
            boolean isRetweeted = retweetedTweetIds.contains(tweet.getId());
            boolean isFollowingAuthor = followedAuthorIds.contains(tweet.getUser().getId());

            return toResponse(tweet, isLiked, isRetweeted, isFollowingAuthor);
        });

        return PageResponse.from(mappedPage);
    }
}