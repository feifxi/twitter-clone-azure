package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type searchQueryRequest struct {
	Query string `form:"q" binding:"required"`
}

func (server *Server) searchUsers(ctx *gin.Context) {
	var req searchQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		writeError(ctx, err)
		return
	}
	offset, size, ok := parseOffsetAndSize(ctx)
	if !ok {
		return
	}
	page := offset / size
	var viewerID *int64
	if id, ok := getCurrentUserID(ctx); ok {
		viewerID = &id
	}
	users, err := server.searchUC.SearchUsers(ctx, req.Query, page, size+1, viewerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	response := newUserResponseList(users)
	ctx.JSON(http.StatusOK, buildPageResponse(response, size, offset))
}

func (server *Server) searchTweets(ctx *gin.Context) {
	var req searchQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		writeError(ctx, err)
		return
	}
	offset, size, ok := parseOffsetAndSize(ctx)
	if !ok {
		return
	}
	page := offset / size
	var viewerID *int64
	if id, ok := getCurrentUserID(ctx); ok {
		viewerID = &id
	}
	tweets, err := server.searchUC.SearchTweets(ctx, req.Query, page, size+1, viewerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	response := newTweetResponseList(tweets)
	ctx.JSON(http.StatusOK, buildPageResponse(response, size, offset))
}

func (server *Server) searchHashtags(ctx *gin.Context) {
	query := strings.TrimSpace(ctx.Query("q"))
	limit := int32(5)
	if rawLimit := strings.TrimSpace(ctx.Query("limit")); rawLimit != "" {
		if parsed, err := strconv.ParseInt(rawLimit, 10, 32); err == nil && parsed > 0 {
			limit = int32(parsed)
		}
	}
	if limit > maxSize {
		limit = maxSize
	}
	hashtags, err := server.searchUC.SearchHashtags(ctx, query, limit)
	if err != nil {
		writeError(ctx, err)
		return
	}
	response := newHashtagResponseList(hashtags)
	ctx.JSON(http.StatusOK, response)
}
