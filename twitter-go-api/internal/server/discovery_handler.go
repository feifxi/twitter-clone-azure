package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (server *Server) getTrendingHashtags(ctx *gin.Context) {
	limit := int32(10)
	if rawLimit := strings.TrimSpace(ctx.Query("limit")); rawLimit != "" {
		if parsed, err := strconv.ParseInt(rawLimit, 10, 32); err == nil && parsed > 0 {
			limit = int32(parsed)
		}
	}
	if limit > maxSize {
		limit = maxSize
	}
	hashtags, err := server.usecase.GetTrendingHashtags(ctx, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, hashtags)
}

func (server *Server) getSuggestedUsers(ctx *gin.Context) {
	page, size, ok := parsePageAndSize(ctx)
	if !ok {
		return
	}
	var viewerID *int64
	if id, ok := getCurrentUserID(ctx); ok {
		viewerID = &id
	}
	users, followingMap, err := server.usecase.GetSuggestedUsers(ctx, page, size, viewerID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := make([]userResponse, 0, len(users))
	for _, user := range users {
		var following *bool
		if v, ok := followingMap[user.ID]; ok {
			f := v
			following = &f
		}
		response = append(response, newUserResponse(user, following))
	}
	ctx.JSON(http.StatusOK, response)
}
