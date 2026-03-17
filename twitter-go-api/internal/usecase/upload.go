package usecase

import (
	"context"
	"strings"

	"github.com/chanombude/twitter-go-api/internal/apperr"
)

var (
	presignAllowedMIMEs = map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
		"video/mp4":  true,
		"video/webm": true,
	}

	presignAllowedFolders = map[string]bool{
		"tweets":  true,
		"avatars": true,
		"banners": true,
	}
)

func (u *UploadUsecase) GeneratePresignedURL(ctx context.Context, filename, contentType, folder string, contentLength *int64) (string, string, error) {
	ctype := strings.ToLower(strings.TrimSpace(contentType))
	if !presignAllowedMIMEs[ctype] {
		return "", "", apperr.BadRequest("unsupported content type")
	}

	fldr := strings.ToLower(strings.TrimSpace(folder))
	if !presignAllowedFolders[fldr] {
		return "", "", apperr.BadRequest("unsupported folder")
	}

	if contentLength != nil {
		maxSize := u.maxUploadSize(fldr)
		if maxSize > 0 && *contentLength > maxSize {
			return "", "", apperr.BadRequest("file too large")
		}
	}

	presignedURL, objectKey, err := u.storage.GeneratePresignedURL(ctx, filename, ctype, fldr)
	if err != nil {
		return "", "", apperr.Internal("failed to generate upload URL", err)
	}

	return presignedURL, objectKey, nil
}

func (u *UploadUsecase) maxUploadSize(folder string) int64 {
	switch folder {
	case "avatars":
		return u.config.MaxAvatarBytes
	case "banners":
		return u.config.MaxBannerBytes
	case "tweets":
		return u.config.MaxMediaBytes
	default:
		return 0
	}
}
