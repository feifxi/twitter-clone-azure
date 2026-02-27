package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/chanombude/twitter-go-api/internal/db"
	"github.com/gin-gonic/gin"
)

// A simple in-memory client manager for SSE.
type sseClient struct {
	channel chan notificationResponse
}

var (
	clients = make(map[int64][]*sseClient)
	mu      sync.RWMutex
)

type notificationResponse struct {
	ID          int64     `json:"id"`
	RecipientID int64     `json:"recipient_id"`
	ActorID     int64     `json:"actor_id"`
	TweetID     *int64    `json:"tweet_id,omitempty"`
	Type        string    `json:"type"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

func newNotificationResponse(n db.Notification) notificationResponse {
	var tweetID *int64
	if n.TweetID.Valid {
		tweetID = &n.TweetID.Int64
	}

	return notificationResponse{
		ID:          n.ID,
		RecipientID: n.RecipientID,
		ActorID:     n.ActorID,
		TweetID:     tweetID,
		Type:        n.Type,
		IsRead:      n.IsRead,
		CreatedAt:   n.CreatedAt,
	}
}

func sendNotificationToUser(userID int64, notification notificationResponse) {
	mu.RLock()
	userClients, ok := clients[userID]
	mu.RUnlock()
	if !ok {
		return
	}

	for _, client := range userClients {
		select {
		case client.channel <- notification:
		default:
		}
	}
}

func (server *Server) streamNotifications(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")

	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}

	client := &sseClient{channel: make(chan notificationResponse, 10)}

	mu.Lock()
	clients[userID] = append(clients[userID], client)
	mu.Unlock()

	fmt.Fprintf(ctx.Writer, "event: connected\ndata: {\"status\": \"ok\"}\n\n")
	ctx.Writer.Flush()

	defer func() {
		mu.Lock()
		userClients := clients[userID]
		for i, c := range userClients {
			if c == client {
				clients[userID] = append(userClients[:i], userClients[i+1:]...)
				break
			}
		}
		mu.Unlock()
		close(client.channel)
	}()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Request.Context().Done():
			return
		case <-ticker.C:
			fmt.Fprintf(ctx.Writer, "event: ping\ndata: {}\n\n")
			ctx.Writer.Flush()
		case notification := <-client.channel:
			data, _ := json.Marshal(notification)
			fmt.Fprintf(ctx.Writer, "event: notification\ndata: %s\n\n", data)
			ctx.Writer.Flush()
		}
	}
}

func (server *Server) listNotifications(ctx *gin.Context) {
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}
	page, size, ok := parsePageAndSize(ctx)
	if !ok {
		return
	}
	notifications, err := server.usecase.ListNotifications(ctx, userID, page, size)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := make([]notificationResponse, 0, len(notifications))
	for _, n := range notifications {
		response = append(response, newNotificationResponse(n))
	}
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) getUnreadNotificationCount(ctx *gin.Context) {
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}
	count, err := server.usecase.CountUnreadNotifications(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, count)
}

func (server *Server) markNotificationRead(ctx *gin.Context) {
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}
	if err := server.usecase.MarkAllNotificationsRead(ctx, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
