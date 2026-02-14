package com.fei.twitterbackend.model.event;

import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.User;
import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
@AllArgsConstructor
public class UserRetweetedEvent {
    private final User actor;
    private final Tweet targetTweet; // The tweet being retweeted
}