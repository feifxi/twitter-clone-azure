package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (server *Server) healthz(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (server *Server) readyz(ctx *gin.Context) {
	checkCtx, cancel := context.WithTimeout(ctx.Request.Context(), 2*time.Second)
	defer cancel()

	if err := server.store.Ping(checkCtx); err != nil {
		ctx.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "not_ready",
			"service": "database",
		})
		return
	}

	if server.redis != nil {
		if err := server.redis.Ping(checkCtx).Err(); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "not_ready",
				"service": "redis",
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ready"})
}
