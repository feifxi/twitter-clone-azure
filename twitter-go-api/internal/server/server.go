package server

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chanombude/twitter-go-api/internal/config"
	db "github.com/chanombude/twitter-go-api/internal/db"
	"github.com/chanombude/twitter-go-api/internal/middleware"
	"github.com/chanombude/twitter-go-api/internal/service"
	"github.com/chanombude/twitter-go-api/internal/token"
	"github.com/chanombude/twitter-go-api/internal/usecase"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config     config.Config
	store      db.Querier
	tokenMaker token.Maker
	storage    service.StorageService
	usecase    *usecase.Usecase
	router     *gin.Engine
}

func NewServer(config config.Config, store db.Querier) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	storageService, err := service.NewAzureStorageService(config)
	if err != nil {
		return nil, fmt.Errorf("cannot create storage service: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		storage:    storageService,
	}
	server.usecase = usecase.New(config, store, tokenMaker, storageService, server.publishNotification)

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	allowOrigins := []string{"http://localhost:3000"}
	if strings.TrimSpace(server.config.FrontendURL) != "" {
		allowOrigins = strings.Split(server.config.FrontendURL, ",")
		for i := range allowOrigins {
			allowOrigins[i] = strings.TrimSpace(allowOrigins[i])
		}
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Basic anti-spam rate limiter.
	router.Use(middleware.RateLimiter(5, 10))

	router.POST("/api/v1/auth/google", server.loginGoogle)
	router.POST("/api/v1/auth/refresh", server.refreshToken)
	router.POST("/api/v1/auth/logout", server.logout)

	publicRoutes := router.Group("/api/v1")
	publicRoutes.Use(middleware.OptionalAuthMiddleware(server.tokenMaker))
	publicRoutes.GET("/users/:id", server.getUser)
	publicRoutes.GET("/users/:id/followers", server.listFollowers)
	publicRoutes.GET("/users/:id/following", server.listFollowing)
	publicRoutes.GET("/tweets/:id", server.getTweet)
	publicRoutes.GET("/tweets/:id/replies", server.getReplies)
	publicRoutes.GET("/feeds/global", server.getGlobalFeed)
	publicRoutes.GET("/feeds/user/:userId", server.getUserFeed)
	publicRoutes.GET("/search/users", server.searchUsers)
	publicRoutes.GET("/search/tweets", server.searchTweets)
	publicRoutes.GET("/search/hashtags", server.searchHashtags)
	publicRoutes.GET("/discovery/trending", server.getTrendingHashtags)
	publicRoutes.GET("/discovery/users", server.getSuggestedUsers)

	authRoutes := router.Group("/api/v1")
	authRoutes.Use(middleware.AuthMiddleware(server.tokenMaker))
	authRoutes.GET("/auth/me", server.getMe)
	authRoutes.GET("/users/me", server.getMe)
	authRoutes.PUT("/users/profile", server.updateProfile)
	authRoutes.POST("/users/:id/follow", server.followUser)
	authRoutes.DELETE("/users/:id/follow", server.unfollowUser)
	authRoutes.POST("/tweets", server.createTweet)
	authRoutes.DELETE("/tweets/:id", server.deleteTweet)
	authRoutes.POST("/tweets/:id/like", server.likeTweet)
	authRoutes.DELETE("/tweets/:id/like", server.unlikeTweet)
	authRoutes.POST("/tweets/:id/retweet", server.retweet)
	authRoutes.DELETE("/tweets/:id/retweet", server.undoRetweet)
	authRoutes.GET("/feeds/following", server.getFollowingFeed)
	authRoutes.GET("/notifications", server.listNotifications)
	authRoutes.GET("/notifications/stream", server.streamNotifications)
	authRoutes.GET("/notifications/unread-count", server.getUnreadNotificationCount)
	authRoutes.POST("/notifications/mark-read", server.markNotificationRead)

	server.router = router
}

func (server *Server) HTTPServer(address string) *http.Server {
	return &http.Server{
		Addr:    address,
		Handler: server.router,
	}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) publishNotification(notification db.Notification) {
	sendNotificationToUser(notification.RecipientID, newNotificationResponse(notification))
}
