/* Flyway V1: Twitter Clone Complete Schema
   Features: Users, Refresh Tokens, Tweets, Interactions
   Updates: All timestamps are NOT NULL with defaults
*/

-- 1. USERS
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(100),
    avatar_url TEXT,
    role VARCHAR(20) NOT NULL DEFAULT 'USER',
    provider VARCHAR(50) NOT NULL,
    followers_count INT NOT NULL DEFAULT 0,
    following_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 2. REFRESH TOKENS
CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expiry_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 3. TWEETS
CREATE TABLE tweets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- Content can be null if it's a simple Retweet (no quote)
    content VARCHAR(280),
    media_type VARCHAR(20),
    media_url TEXT,
    -- Self-References
    parent_id BIGINT REFERENCES tweets(id) ON DELETE CASCADE, -- Reply
    retweet_id BIGINT REFERENCES tweets(id) ON DELETE CASCADE, -- Retweet
    -- Counters
    reply_count INT NOT NULL DEFAULT 0,
    retweet_count INT NOT NULL DEFAULT 0,
    like_count INT NOT NULL DEFAULT 0,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for Feed & Lookups
CREATE INDEX idx_tweets_user_created ON tweets(user_id, created_at DESC); -- Critical for Profile Feed
CREATE INDEX idx_tweets_parent ON tweets(parent_id);
CREATE INDEX idx_tweets_retweet ON tweets(retweet_id);

-- 4. LIKES
CREATE TABLE tweet_likes (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tweet_id BIGINT NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, tweet_id) -- Composite PK is faster than a surrogate serial ID
);


-- 5. FOLLOWS
CREATE TABLE follows (
    follower_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id != following_id)
);

-- Indexes for "Who follows X" and "Who is X following"
CREATE INDEX idx_follows_follower ON follows(follower_id);
CREATE INDEX idx_follows_following ON follows(following_id);