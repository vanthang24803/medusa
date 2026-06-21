package identity

import (
	"context"
	"fmt"
	"io"
	"time"

	"ecommerce/packages/events"
	"ecommerce/packages/upload"
)

const maxAvatarSize = 5 << 20 // 5 MB

type Service interface {
	Ping(ctx context.Context) error
	GetProfile(ctx context.Context, userID string) (*User, error)
	UploadAvatar(ctx context.Context, userID string, file io.Reader, size int64, contentType string) (*User, error)
}

type service struct {
	repo     *Repository
	bus      events.EventBus
	uploader upload.Uploader
}

func NewService(repo *Repository, bus events.EventBus, uploader upload.Uploader) Service {
	return &service{repo: repo, bus: bus, uploader: uploader}
}

func (s *service) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *service) GetProfile(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetByID(ctx, userID)
}

func (s *service) UploadAvatar(ctx context.Context, userID string, file io.Reader, size int64, contentType string) (*User, error) {
	if err := upload.ValidateSize(size, maxAvatarSize); err != nil {
		return nil, err
	}

	objectName := fmt.Sprintf("identity/%s/avatar%s", userID, upload.ExtFromContentType(contentType))

	url, err := s.uploader.Upload(ctx, "avatars", objectName, file, size, contentType)
	if err != nil {
		return nil, fmt.Errorf("upload avatar: %w", err)
	}

	url = fmt.Sprintf("%s?v=%d", url, time.Now().Unix())

	if err := s.repo.UpdateAvatarURL(ctx, userID, url); err != nil {
		return nil, fmt.Errorf("persist avatar url: %w", err)
	}

	return s.repo.GetByID(ctx, userID)
}
