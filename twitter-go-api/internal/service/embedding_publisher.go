package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/chanombude/twitter-go-api/internal/config"
	"github.com/rs/zerolog/log"
)

type EmbeddingMessage struct {
	TweetID int64  `json:"tweet_id"`
	Content string `json:"content"`
}

type EmbeddingPublisher interface {
	PublishEmbeddingEvent(ctx context.Context, tweetID int64, content string) error
}

type SQSEmbeddingPublisher struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSEmbeddingPublisher(cfg config.Config) (EmbeddingPublisher, error) {
	if !cfg.EnableRAG {
		log.Debug().Msg("RAG is disabled — skipping SQS embedding publisher initialization")
		return &noOpEmbeddingPublisher{}, nil
	}

	if cfg.SQSEmbeddingQueueURL == "" {
		log.Warn().Msg("SQS_EMBEDDING_QUEUE_URL is not set — embedding events will be dropped")
		return &noOpEmbeddingPublisher{}, nil
	}

	region := cfg.S3Region
	if region == "" {
		region = "ap-southeast-1"
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config for SQS: %w", err)
	}

	client := sqs.NewFromConfig(awsCfg)

	return &SQSEmbeddingPublisher{
		client:   client,
		queueURL: cfg.SQSEmbeddingQueueURL,
	}, nil
}

func (p *SQSEmbeddingPublisher) PublishEmbeddingEvent(ctx context.Context, tweetID int64, content string) error {
	msg := EmbeddingMessage{
		TweetID: tweetID,
		Content: content,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal embedding message: %w", err)
	}

	_, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(p.queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("failed to send embedding message to SQS: %w", err)
	}

	log.Info().Int64("tweet_id", tweetID).Msg("Successfully published embedding event to SQS")
	return nil
}

type noOpEmbeddingPublisher struct{}

func (n *noOpEmbeddingPublisher) PublishEmbeddingEvent(ctx context.Context, tweetID int64, content string) error {
	log.Debug().Int64("tweet_id", tweetID).Msg("Dropped embedding event due to missing SQS configuration")
	return nil
}
