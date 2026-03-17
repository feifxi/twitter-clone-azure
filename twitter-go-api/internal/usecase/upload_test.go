package usecase_test

import (
	"context"
	"testing"

	"github.com/chanombude/twitter-go-api/internal/apperr"
	"github.com/chanombude/twitter-go-api/internal/config"
	"github.com/chanombude/twitter-go-api/internal/usecase"
)

type mockStorageService struct {
	generatePresignedURLFn func(ctx context.Context, filename, contentType, folder string) (string, string, error)
}

func (m *mockStorageService) GeneratePresignedURL(ctx context.Context, filename, contentType, folder string) (string, string, error) {
	if m.generatePresignedURLFn != nil {
		return m.generatePresignedURLFn(ctx, filename, contentType, folder)
	}
	return "http://presigned.url", "folder/key.png", nil
}
func (m *mockStorageService) DeleteFile(ctx context.Context, objectKey string) error { return nil }
func (m *mockStorageService) PublicURL(objectKey string) string                      { return "" }

func TestUploadUsecase_GeneratePresignedURL(t *testing.T) {
	cfg := config.Config{MaxBannerBytes: 10 << 20} // 10MB
	uc := usecase.NewUploadUsecase(cfg, &mockStorageService{})

	t.Run("success", func(t *testing.T) {
		var contentLength int64 = 5 << 20
		url, key, err := uc.GeneratePresignedURL(context.Background(), "file.png", "image/png", "banners", &contentLength)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if url == "" || key == "" {
			t.Fatalf("expected valid url and key")
		}
	})

	t.Run("invalid_mime", func(t *testing.T) {
		_, _, err := uc.GeneratePresignedURL(context.Background(), "file.exe", "application/x-msdos-program", "banners", nil)
		if err == nil {
			t.Fatal("expected error for invalid mime type")
		}
		if kind, ok := apperr.KindOf(err); !ok || kind != apperr.KindBadRequest {
			t.Fatalf("expected bad request, got %v", err)
		}
	})

	t.Run("invalid_folder", func(t *testing.T) {
		_, _, err := uc.GeneratePresignedURL(context.Background(), "file.png", "image/png", "malicious", nil)
		if err == nil {
			t.Fatal("expected error for invalid folder")
		}
	})

	t.Run("file_too_large", func(t *testing.T) {
		var oversized int64 = 11 << 20
		_, _, err := uc.GeneratePresignedURL(context.Background(), "big.png", "image/png", "banners", &oversized)
		if err == nil {
			t.Fatal("expected error for oversized file")
		}
	})
}
