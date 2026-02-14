package com.fei.twitterbackend.model.event;

import com.fei.twitterbackend.model.entity.User;
import lombok.AllArgsConstructor;
import lombok.Getter;

@Getter
@AllArgsConstructor
public class UserFollowedEvent {
    private final User actor;   // Who followed?
    private final User target;  // Who was followed?
}