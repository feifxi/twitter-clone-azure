-- name: ListRelatedTweetsByEmbedding :many
SELECT tweet_id, content
FROM tweet_embeddings
ORDER BY embedding <-> $1::vector ASC
LIMIT $2;
