package com.fei.twitterbackend;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.scheduling.annotation.EnableAsync;

@SpringBootApplication
@EnableAsync
public class TwitterBackendApplication {

    public static void main(String[] args) {
        SpringApplication.run(TwitterBackendApplication.class, args);
    }

}
