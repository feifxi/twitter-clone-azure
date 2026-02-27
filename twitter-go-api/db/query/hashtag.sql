-- name: UpsertHashtag :one
INSERT INTO hashtags (text, usage_count, last_used_at)
VALUES ($1, 1, CURRENT_TIMESTAMP)
ON CONFLICT (text) DO UPDATE
SET usage_count = hashtags.usage_count + 1,
    last_used_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: LinkTweetHashtag :exec
INSERT INTO tweet_hashtags (tweet_id, hashtag_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: SearchHashtagsByPrefix :many
SELECT * FROM hashtags
WHERE LOWER(text) LIKE LOWER($1 || '%')
ORDER BY usage_count DESC
LIMIT $2;

-- name: GetTrendingHashtagsLast24h :many
SELECT h.*
FROM hashtags h
JOIN tweet_hashtags th ON th.hashtag_id = h.id
JOIN tweets t ON t.id = th.tweet_id
WHERE t.created_at >= NOW() - INTERVAL '24 hours'
GROUP BY h.id
ORDER BY COUNT(th.tweet_id) DESC, MAX(t.created_at) DESC
LIMIT $1;

-- name: GetTopHashtagsAllTime :many
SELECT * FROM hashtags
ORDER BY usage_count DESC, last_used_at DESC
LIMIT $1;
