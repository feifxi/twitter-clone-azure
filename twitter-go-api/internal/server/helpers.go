package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/chanombude/twitter-go-api/internal/middleware"
	"github.com/chanombude/twitter-go-api/internal/token"
	"github.com/gin-gonic/gin"
)

const (
	defaultPage = int32(0)
	defaultSize = int32(20)
	maxSize     = int32(50)
)

func getCurrentUserID(ctx *gin.Context) (int64, bool) {
	payload, ok := ctx.Get(middleware.AuthorizationPayloadKey)
	if !ok {
		return 0, false
	}
	authPayload, ok := payload.(*token.Payload)
	if !ok {
		return 0, false
	}
	return authPayload.UserID, true
}

func mustCurrentUserID(ctx *gin.Context) (int64, bool) {
	userID, ok := getCurrentUserID(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return 0, false
	}
	return userID, true
}

func parsePageAndSize(ctx *gin.Context) (int32, int32, bool) {
	page := defaultPage
	size := defaultSize

	if rawPage := strings.TrimSpace(ctx.Query("page")); rawPage != "" {
		value, err := strconv.ParseInt(rawPage, 10, 32)
		if err != nil || value < 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "page must be >= 0"})
			return 0, 0, false
		}
		page = int32(value)
	}

	if rawSize := strings.TrimSpace(ctx.Query("size")); rawSize != "" {
		value, err := strconv.ParseInt(rawSize, 10, 32)
		if err != nil || value <= 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "size must be > 0"})
			return 0, 0, false
		}
		size = int32(value)
	}

	if size > maxSize {
		size = maxSize
	}

	return page, size, true
}

func toOffset(page, size int32) int32 {
	return page * size
}
