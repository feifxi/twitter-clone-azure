package server

import (
	"github.com/chanombude/twitter-go-api/internal/apperr"
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
		writeError(ctx, apperr.Unauthorized("authentication required"))
		return 0, false
	}
	return userID, true
}

func parsePageAndSize(ctx *gin.Context) (int32, int32, bool) {
	type paginationQuery struct {
		Page *int32 `form:"page" binding:"omitempty,min=0"`
		Size *int32 `form:"size" binding:"omitempty,min=1"`
	}

	var req paginationQuery
	if err := ctx.ShouldBindQuery(&req); err != nil {
		writeError(ctx, err)
		return 0, 0, false
	}

	page := defaultPage
	if req.Page != nil {
		page = *req.Page
	}

	size := defaultSize
	if req.Size != nil {
		size = min(*req.Size, maxSize)
	}

	return page, size, true
}
