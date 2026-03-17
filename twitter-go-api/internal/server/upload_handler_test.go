package server

import (
	"bytes"
	"context"
	"net/http"
	"testing"

	"github.com/chanombude/twitter-go-api/internal/apperr"
)

type mockUploadUC struct {
	generatePresignedURLFn func(ctx context.Context, filename, contentType, folder string, contentLength *int64) (string, string, error)
}

func (m *mockUploadUC) GeneratePresignedURL(ctx context.Context, filename, contentType, folder string, contentLength *int64) (string, string, error) {
	return m.generatePresignedURLFn(ctx, filename, contentType, folder, contentLength)
}

func TestPresignUpload_Success(t *testing.T) {
	t.Parallel()

	mock := &mockUploadUC{
		generatePresignedURLFn: func(_ context.Context, filename, contentType, folder string, contentLength *int64) (string, string, error) {
			return "http://presigned.url", "banners/uuid_file.png", nil
		},
	}

	reqBody := `{"filename":"banner.png","contentType":"image/png","folder":"banners"}`
	ctx, rec := newHandlerTestContext(http.MethodPost, "/api/v1/uploads/presign", bytes.NewBufferString(reqBody), "application/json")
	setAuthorizedUser(ctx, 1)

	s := &Server{
		uploadUC: mock,
	}
	s.presignUpload(ctx)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestPresignUpload_UsecaseError(t *testing.T) {
	t.Parallel()

	mock := &mockUploadUC{
		generatePresignedURLFn: func(_ context.Context, filename, contentType, folder string, contentLength *int64) (string, string, error) {
			return "", "", apperr.BadRequest("unsupported folder")
		},
	}

	reqBody := `{"filename":"file.png","contentType":"image/png","folder":"malicious"}`
	ctx, rec := newHandlerTestContext(http.MethodPost, "/api/v1/uploads/presign", bytes.NewBufferString(reqBody), "application/json")
	setAuthorizedUser(ctx, 1)

	s := &Server{
		uploadUC: mock,
	}
	s.presignUpload(ctx)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad request from usecase, got %d, body=%s", rec.Code, rec.Body.String())
	}
}
