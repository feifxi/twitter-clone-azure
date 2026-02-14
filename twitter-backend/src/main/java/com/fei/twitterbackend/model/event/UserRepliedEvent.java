package com.fei.twitterbackend.model.event;

import com.fei.twitterbackend.model.entity.Tweet;
import com.fei.twitterbackend.model.entity.User;
import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
@AllArgsConstructor
public class UserRepliedEvent {
    private final User actor;
    private final Tweet parentTweet; // The tweet being replied TO
    private final Tweet replyTweet;  // The new reply itself
}