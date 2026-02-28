package usecase

import (
	"context"
	"database/sql"
	"regexp"
	"strings"

	"github.com/chanombude/twitter-go-api/internal/config"
	"github.com/chanombude/twitter-go-api/internal/db"
	"github.com/chanombude/twitter-go-api/internal/service"
	"github.com/chanombude/twitter-go-api/internal/token"
	"github.com/redis/go-redis/v9"
)

type TweetItem struct {
	db.Tweet
	Author         UserItem
	IsLiked        bool
	IsRetweeted    bool
	IsFollowing    bool
	ParentUsername *string
	OriginalTweet  *TweetItem
}

type UserItem struct {
	db.User
	IsFollowing bool
}

type Usecase struct {
	config              config.Config
	store               db.Querier
	tokenMaker          token.Maker
	storage             service.StorageService
	redis               *redis.Client
	publishNotification func(db.Notification)
}

func New(cfg config.Config, store db.Querier, tokenMaker token.Maker, storage service.StorageService, redisClient *redis.Client, publishNotification func(db.Notification)) *Usecase {
	return &Usecase{
		config:              cfg,
		store:               store,
		tokenMaker:          tokenMaker,
		storage:             storage,
		redis:               redisClient,
		publishNotification: publishNotification,
	}
}

func (u *Usecase) Store() db.Querier {
	return u.store
}

func nullStringFromPtr(v *string) sql.NullString {
	if v == nil {
		return sql.NullString{Valid: false}
	}
	trimmed := strings.TrimSpace(*v)
	if trimmed == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: trimmed, Valid: true}
}

var hashtagRegex = regexp.MustCompile(`(?i)(?:^|\s)#([a-z0-9_]+)`)

func extractHashtags(content string) []string {
	matches := hashtagRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil
	}

	seen := make(map[string]struct{})
	result := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		tag := strings.ToLower(strings.TrimSpace(m[1]))
		if tag == "" {
			continue
		}
		if _, exists := seen[tag]; exists {
			continue
		}
		seen[tag] = struct{}{}
		result = append(result, tag)
	}

	return result
}

func buildTSQuery(raw string) string {
	clean := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == ' ' {
			return r
		}
		return ' '
	}, raw)
	parts := strings.Fields(clean)
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " & ")
}

func (u *Usecase) createAndDispatchNotification(ctx context.Context, recipientID, actorID int64, tweetID *int64, typ string) error {
	if recipientID == actorID {
		return nil
	}

	arg := db.CreateNotificationParams{
		RecipientID: recipientID,
		ActorID:     actorID,
		Type:        typ,
		TweetID:     sql.NullInt64{Valid: false},
	}
	if tweetID != nil {
		arg.TweetID = sql.NullInt64{Int64: *tweetID, Valid: true}
	}

	notification, err := u.store.CreateNotification(ctx, arg)
	if err != nil {
		return err
	}

	if u.publishNotification != nil {
		u.publishNotification(notification)
	}
	return nil
}
