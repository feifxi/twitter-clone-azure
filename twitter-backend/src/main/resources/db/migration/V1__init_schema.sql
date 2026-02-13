/* Flyway V1: Twitter Clone Complete Schema
   Features: Users, Refresh Tokens, Tweets, Interactions
*/

-- 1. USERS
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(100),
    bio VARCHAR(160),
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

CREATE INDEX idx_tweets_user_created ON tweets(user_id, created_at DESC);
CREATE INDEX idx_tweets_parent ON tweets(parent_id);
CREATE INDEX idx_tweets_retweet ON tweets(retweet_id);
CREATE INDEX idx_tweets_created_at ON tweets(created_at DESC);

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

CREATE INDEX idx_follows_follower ON follows(follower_id);
CREATE INDEX idx_follows_following ON follows(following_id);

-- 6. HASHTAGS
CREATE TABLE hashtags (
    id BIGSERIAL PRIMARY KEY,
    text VARCHAR(100) NOT NULL UNIQUE, -- Stored as lowercase for consistency
    usage_count INT NOT NULL DEFAULT 1,
    last_used_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tweet_hashtags (
    tweet_id BIGINT NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    hashtag_id BIGINT NOT NULL REFERENCES hashtags(id) ON DELETE CASCADE,
    PRIMARY KEY (tweet_id, hashtag_id)
);

CREATE INDEX idx_hashtags_text ON hashtags(text);
CREATE INDEX idx_tweet_hashtags_tag ON tweet_hashtags(hashtag_id);

-- 7. FULL TEXT SEARCH
ALTER TABLE tweets ADD COLUMN search_vector tsvector 
    GENERATED ALWAYS AS (to_tsvector('english', coalesce(content, ''))) STORED;

-- GIN Index for lightning-fast text search
CREATE INDEX idx_tweets_search ON tweets USING GIN(search_vector);

-- 8. USER SEARCH OPTIMIZATION 
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_users_username_trgm ON users USING GIN (username gin_trgm_ops);
CREATE INDEX idx_users_displayname_trgm ON users USING GIN (display_name gin_trgm_ops);