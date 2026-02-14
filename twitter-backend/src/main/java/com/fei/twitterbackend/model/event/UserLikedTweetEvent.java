package com.fei.twitterbackend.model.event;

import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.User;
import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
@AllArgsConstructor
public class UserLikedTweetEvent {
    private final User actor;  // Who liked it?
    private final Tweet tweet; // Which tweet?
}
