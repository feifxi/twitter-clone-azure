package server

import (
	"io"
	"net/http"

	"github.com/chanombude/twitter-go-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

func (server *Server) assistant(ctx *gin.Context) {
	var input usecase.AssistantInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		writeError(ctx, err)
		return
	}

	reader, err := server.assistantUC.Chat(ctx.Request.Context(), input)
	if err != nil {
		writeError(ctx, err)
		return
	}

	// Set SSE headers
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Transfer-Encoding", "chunked")

	ctx.Stream(func(w io.Writer) bool {
		buffer := make([]byte, 1024)
		for {
			n, err := reader.Read(buffer)
			if n > 0 {
				_, _ = w.Write(buffer[:n])
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
			}
			if err != nil {
				if err == io.EOF {
					return false // Close stream gracefully
				}
				return false
			}
		}
	})
}
