package usecase

import (
	"github.com/chanombude/twitter-go-api/internal/config"
	"github.com/chanombude/twitter-go-api/internal/db"
	"github.com/chanombude/twitter-go-api/internal/service"
	"github.com/chanombude/twitter-go-api/internal/token"
)

func NewAuthUsecase(cfg config.Config, store db.Store, tokenMaker token.Maker) *AuthUsecase {
	return &AuthUsecase{
		config:     cfg,
		store:      store,
		tokenMaker: tokenMaker,
	}
}

func NewUserUsecase(store db.Store, storage service.StorageService, publishNotification func(db.Notification)) *UserUsecase {
	return &UserUsecase{
		store:               store,
		storage:             storage,
		publishNotification: publishNotification,
	}
}

func NewTweetUsecase(store db.Store, storage service.StorageService, publishNotification func(db.Notification)) *TweetUsecase {
	return &TweetUsecase{
		store:               store,
		storage:             storage,
		publishNotification: publishNotification,
	}
}

func NewFeedUsecase(store db.Store) *FeedUsecase {
	return &FeedUsecase{store: store}
}

func NewSearchUsecase(store db.Store) *SearchUsecase {
	return &SearchUsecase{store: store}
}

func NewDiscoveryUsecase(store db.Store) *DiscoveryUsecase {
	return &DiscoveryUsecase{store: store}
}

func NewNotificationUsecase(store db.Store) *NotificationUsecase {
	return &NotificationUsecase{store: store}
}

func NewMessageUsecase(store db.Store) *MessageUsecase {
	return &MessageUsecase{store: store}
}
