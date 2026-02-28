package server

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/chanombude/twitter-go-api/internal/apperr"
	"github.com/gin-gonic/gin"
)

type apiErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeError(ctx *gin.Context, err error) {
	if err == nil {
		return
	}

	status := http.StatusInternalServerError
	code := "INTERNAL_ERROR"
	message := "internal server error"

	if errors.Is(err, sql.ErrNoRows) {
		status = http.StatusNotFound
		code = "NOT_FOUND"
		message = "resource not found"
	} else if kind, ok := apperr.KindOf(err); ok {
		switch kind {
		case apperr.KindBadRequest:
			status = http.StatusBadRequest
			code = "BAD_REQUEST"
			message = defaultMessage(apperr.MessageOf(err), "invalid request")
		case apperr.KindUnauthorized:
			status = http.StatusUnauthorized
			code = "UNAUTHORIZED"
			message = defaultMessage(apperr.MessageOf(err), "authentication required")
		case apperr.KindForbidden:
			status = http.StatusForbidden
			code = "FORBIDDEN"
			message = defaultMessage(apperr.MessageOf(err), "forbidden")
		case apperr.KindNotFound:
			status = http.StatusNotFound
			code = "NOT_FOUND"
			message = defaultMessage(apperr.MessageOf(err), "resource not found")
		case apperr.KindConflict:
			status = http.StatusConflict
			code = "CONFLICT"
			message = defaultMessage(apperr.MessageOf(err), "conflict")
		case apperr.KindInternal:
			status = http.StatusInternalServerError
			code = "INTERNAL_ERROR"
			message = "internal server error"
		}
	}

	if status >= 500 {
		log.Printf("[ERROR] %s %s: %v", ctx.Request.Method, ctx.FullPath(), err)
	}

	ctx.JSON(status, apiErrorResponse{Code: code, Message: message})
}

func defaultMessage(in, fallback string) string {
	if in == "" {
		return fallback
	}
	return in
}
