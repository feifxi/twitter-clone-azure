package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type presignRequest struct {
	Filename      string `json:"filename" binding:"required"`
	ContentType   string `json:"contentType" binding:"required"`
	Folder        string `json:"folder" binding:"required"`
	ContentLength *int64 `json:"contentLength" binding:"omitempty,min=1"`
}

func (server *Server) presignUpload(ctx *gin.Context) {
	if _, ok := mustCurrentUserID(ctx); !ok {
		return
	}

	var req presignRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		writeError(ctx, err)
		return
	}

	presignedURL, objectKey, err := server.uploadUC.GeneratePresignedURL(ctx, req.Filename, req.ContentType, req.Folder, req.ContentLength)
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"presignedUrl": presignedURL,
		"objectKey":    objectKey,
	})
}
