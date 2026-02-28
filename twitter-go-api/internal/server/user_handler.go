package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chanombude/twitter-go-api/internal/apperr"
	"github.com/chanombude/twitter-go-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

func (server *Server) getMe(ctx *gin.Context) {
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}

	user, err := server.usecase.GetMe(ctx, userID)
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user, nil))
}

type getUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, apperr.BadRequest("invalid user id"))
		return
	}

	var viewerID *int64
	if id, ok := getCurrentUserID(ctx); ok {
		viewerID = &id
	}

	user, following, err := server.usecase.GetUser(ctx, req.ID, viewerID)
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(user, following))
}

type updateProfileRequest struct {
	Bio         *string `json:"bio"`
	DisplayName *string `json:"display_name"`
}

func (server *Server) updateProfile(ctx *gin.Context) {
	userID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}

	request := updateProfileRequest{}
	input := usecase.UpdateProfileInput{}

	if strings.HasPrefix(ctx.GetHeader("Content-Type"), "multipart/form-data") {
		if err := ctx.Request.ParseMultipartForm(20 << 20); err != nil {
			writeError(ctx, apperr.BadRequest("invalid multipart payload"))
			return
		}

		if dataBlob := ctx.Request.FormValue("data"); dataBlob != "" {
			if err := json.Unmarshal([]byte(dataBlob), &request); err != nil {
				writeError(ctx, apperr.BadRequest("invalid json in data field"))
				return
			}
		}

		file, header, err := ctx.Request.FormFile("avatar")
		if err == nil {
			defer file.Close()
			contentType := header.Header.Get("Content-Type")
			if !strings.HasPrefix(contentType, "image/") {
				writeError(ctx, apperr.BadRequest("avatar must be an image"))
				return
			}
			input.Avatar = &usecase.AvatarUpload{
				Filename:    header.Filename,
				ContentType: contentType,
				Reader:      file,
			}
		}
	} else {
		if err := ctx.ShouldBindJSON(&request); err != nil {
			writeError(ctx, apperr.BadRequest("invalid request payload"))
			return
		}
	}

	input.Bio = request.Bio
	input.DisplayName = request.DisplayName

	updatedUser, err := server.usecase.UpdateProfile(ctx, userID, input)
	if err != nil {
		writeError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, newUserResponse(updatedUser, nil))
}

func (server *Server) followUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, apperr.BadRequest("invalid user id"))
		return
	}

	followerID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}
	if followerID == req.ID {
		writeError(ctx, apperr.BadRequest("cannot follow yourself"))
		return
	}

	_, err := server.usecase.FollowUser(ctx, followerID, req.ID)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

func (server *Server) unfollowUser(ctx *gin.Context) {
	var req getUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, apperr.BadRequest("invalid user id"))
		return
	}

	followerID, ok := mustCurrentUserID(ctx)
	if !ok {
		return
	}

	if err := server.usecase.UnfollowUser(ctx, followerID, req.ID); err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"success": true})
}

type listFollowRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) listFollowers(ctx *gin.Context) {
	var req listFollowRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, apperr.BadRequest("invalid user id"))
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

	users, followingMap, err := server.usecase.ListFollowers(ctx, req.ID, page, size, viewerID)
	if err != nil {
		writeError(ctx, err)
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

func (server *Server) listFollowing(ctx *gin.Context) {
	var req listFollowRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		writeError(ctx, apperr.BadRequest("invalid user id"))
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

	users, followingMap, err := server.usecase.ListFollowing(ctx, req.ID, page, size, viewerID)
	if err != nil {
		writeError(ctx, err)
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
