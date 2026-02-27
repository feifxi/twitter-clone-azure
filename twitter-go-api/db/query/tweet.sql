-- name: CreateTweet :one
INSERT INTO tweets (
  user_id, content, media_type, media_url, parent_id, retweet_id
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetTweet :one
SELECT * FROM tweets
WHERE id = $1 LIMIT 1;

-- name: DeleteTweetByOwner :one
DELETE FROM tweets
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DecrementParentReplyCount :exec
UPDATE tweets
SET reply_count = GREATEST(0, reply_count - 1)
WHERE id = $1;

-- name: IncrementParentReplyCount :exec
UPDATE tweets
SET reply_count = reply_count + 1
WHERE id = $1;

-- name: CreateRetweet :one
WITH inserted AS (
  INSERT INTO tweets (user_id, retweet_id, media_type)
  VALUES ($1, $2, 'NONE')
  ON CONFLICT DO NOTHING
  RETURNING *
),
updated AS (
  UPDATE tweets
  SET retweet_count = retweet_count + 1
  WHERE id = $2 AND EXISTS (SELECT 1 FROM inserted)
)
SELECT * FROM inserted;

-- name: DeleteRetweetByUser :one
WITH deleted AS (
DELETE FROM tweets
WHERE tweets.user_id = $1 AND tweets.retweet_id = $2
RETURNING *
),
updated AS (
  UPDATE tweets
  SET retweet_count = GREATEST(0, retweet_count - 1)
  WHERE id = $2 AND EXISTS (SELECT 1 FROM deleted)
)
SELECT * FROM deleted;

-- name: GetUserRetweet :one
SELECT * FROM tweets
WHERE user_id = $1 AND retweet_id = $2
LIMIT 1;

-- name: ListForYouFeed :many
SELECT * FROM tweets t
WHERE t.parent_id IS NULL
ORDER BY
  (t.like_count * 2 + t.retweet_count * 3 + t.reply_count + 1) /
  POWER((EXTRACT(EPOCH FROM NOW() - t.created_at) / 3600) + 2, 1.8) DESC,
  t.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListFollowingFeed :many
SELECT t.*
FROM tweets t
JOIN follows f ON t.user_id = f.following_id
WHERE f.follower_id = $1
  AND t.parent_id IS NULL
ORDER BY t.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListUserTweets :many
SELECT * FROM tweets
WHERE user_id = $1
  AND parent_id IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListTweetReplies :many
SELECT * FROM tweets
WHERE parent_id = $1
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;

-- name: SearchTweetsFullText :many
SELECT * FROM tweets
WHERE search_vector @@ to_tsquery('english', $1)
ORDER BY ts_rank(search_vector, to_tsquery('english', $1)) DESC, created_at DESC
LIMIT $2 OFFSET $3;

-- name: SearchTweetsByHashtag :many
SELECT t.*
FROM tweets t
JOIN tweet_hashtags th ON th.tweet_id = t.id
JOIN hashtags h ON h.id = th.hashtag_id
WHERE LOWER(h.text) = LOWER($1)
ORDER BY t.created_at DESC
LIMIT $2 OFFSET $3;
