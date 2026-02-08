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
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- 2. REFRESH TOKENS
CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expiry_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- 3. TWEETS
CREATE TABLE tweets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content VARCHAR(280),
    media_url TEXT,
    -- Structure
    parent_id BIGINT REFERENCES tweets(id) ON DELETE CASCADE, -- Reply
    retweet_id BIGINT REFERENCES tweets(id) ON DELETE CASCADE, -- Retweet
    -- Stats
    reply_count INT DEFAULT 0,
    retweet_count INT DEFAULT 0,
    like_count INT DEFAULT 0,

    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Indexing for speed
CREATE INDEX idx_tweets_user ON tweets(user_id);
CREATE INDEX idx_tweets_parent ON tweets(parent_id);

-- 4. LIKES
CREATE TABLE likes (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tweet_id BIGINT NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, tweet_id)
);

-- 5. FOLLOWS
CREATE TABLE follows (
    follower_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    following_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id != following_id)
);