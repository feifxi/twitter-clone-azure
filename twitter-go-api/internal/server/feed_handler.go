package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) getGlobalFeed(ctx *gin.Context) {
	page, size, ok := parsePageAndSize(ctx)
	if !ok {
		return
	}
	var viewerID *int64
	if id, ok := getCurrentUserID(ctx); ok {
		viewerID = &id
	}
	tweets, err := server.usecase.GetGlobalFeed(ctx, page, size, viewerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	response := make([]tweetResponse, 0, len(tweets))
	for _, t := range tweets {
		response = append(response, newTweetResponse(t))
	}
	ctx.JSON(http.StatusOK, response)
}

func (server *Server) getFollowingFeed(ctx *gin.Context) {
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}
	page, size, ok := parsePageAndSize(ctx)
	if !ok {
		return
	}
	tweets, err := server.usecase.GetFollowingFeed(ctx, userID, page, size)
	if err != nil {
		writeError(ctx, err)
		return
	}
	response := make([]tweetResponse, 0, len(tweets))
	for _, t := range tweets {
		response = append(response, newTweetResponse(t))
	}
	ctx.JSON(http.StatusOK, response)
}

type userFeedRequest struct {
	UserID int64 `uri:"userId" binding:"required,min=1"`
}

func (server *Server) getUserFeed(ctx *gin.Context) {
	var req userFeedRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, err)
		return
	}
	page, size, ok := parsePageAndSize(ctx)
	if !ok {
		return
	}
	var viewerID *int64
	if id, ok := getCurrentUserID(ctx); ok {
		viewerID = &id
	}
	tweets, err := server.usecase.GetUserFeed(ctx, req.UserID, page, size, viewerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	response := make([]tweetResponse, 0, len(tweets))
	for _, t := range tweets {
		response = append(response, newTweetResponse(t))
	}
	ctx.JSON(http.StatusOK, response)
}
