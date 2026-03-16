CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE IF NOT EXISTS tweet_embeddings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tweet_id BIGINT UNIQUE NOT NULL REFERENCES tweets(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    embedding vector(768) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS tweet_embeddings_embedding_idx 
ON tweet_embeddings USING hnsw (embedding vector_cosine_ops);