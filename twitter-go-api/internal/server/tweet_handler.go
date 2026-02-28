package server

import (
	"encoding/json"
	"net/http"

	"github.com/chanombude/twitter-go-api/internal/apperr"
	"github.com/chanombude/twitter-go-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type createTweetRequest struct {
	Content  *string `json:"content"`
	ParentID *int64  `json:"parentId"`
}

func (server *Server) createTweet(ctx *gin.Context) {
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}

	if err := ctx.Request.ParseMultipartForm(20 << 20); err != nil {
		writeError(ctx, apperr.BadRequest("failed to parse multipart form"))
		return
	}

	dataBlob := ctx.Request.FormValue("data")
	if dataBlob == "" {
		writeError(ctx, apperr.BadRequest("missing data field in form"))
		return
	}

	var req createTweetRequest
	if err := json.Unmarshal([]byte(dataBlob), &req); err != nil {
		writeError(ctx, apperr.BadRequest("invalid json in data field"))
		return
	}

	input := usecase.CreateTweetInput{UserID: userID, Content: req.Content, ParentID: req.ParentID}
	if file, header, err := ctx.Request.FormFile("media"); err == nil {
		defer file.Close()
		input.Media = &usecase.MediaUpload{
			Filename:    header.Filename,
			ContentType: header.Header.Get("Content-Type"),
			Reader:      file,
		}
	}

	tweet, err := server.usecase.CreateTweet(ctx, input)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, newTweetResponse(tweet))
}

type tweetURIRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getTweet(ctx *gin.Context) {
	var req tweetURIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, err)
		return
	}
	var viewerID *int64
	if id, ok := getCurrentUserID(ctx); ok {
		viewerID = &id
	}
	tweet, err := server.usecase.GetTweet(ctx, req.ID, viewerID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, newTweetResponse(tweet))
}

func (server *Server) deleteTweet(ctx *gin.Context) {
	var req tweetURIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, err)
		return
	}
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}
	if err := server.usecase.DeleteTweet(ctx, userID, req.ID); err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

type getRepliesRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getReplies(ctx *gin.Context) {
	var req getRepliesRequest
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
	tweets, err := server.usecase.ListReplies(ctx, req.ID, page, size, viewerID)
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

func (server *Server) likeTweet(ctx *gin.Context) {
	var req tweetURIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, err)
		return
	}
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}
	if err := server.usecase.LikeTweet(ctx, userID, req.ID); err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (server *Server) unlikeTweet(ctx *gin.Context) {
	var req tweetURIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, err)
		return
	}
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}
	if err := server.usecase.UnlikeTweet(ctx, userID, req.ID); err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (server *Server) retweet(ctx *gin.Context) {
	var req tweetURIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, err)
		return
	}
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}
	tweet, err := server.usecase.Retweet(ctx, userID, req.ID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, newTweetResponse(tweet))
}

func (server *Server) undoRetweet(ctx *gin.Context) {
	var req tweetURIRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, err)
		return
	}
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}
	if err := server.usecase.UndoRetweet(ctx, userID, req.ID); err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}
