package customer

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
	GetProfile(ctx context.Context, id string) (*Customer, error)
	UpdateProfile(ctx context.Context, id string, req *UpdateCustomerReq) (*Customer, error)
	UploadAvatar(ctx context.Context, customerID string, file io.Reader, size int64, contentType string) (*Customer, error)
}

type service struct {
	repo     *Repository
	bus      events.EventBus
	uploader upload.Uploader
}

func NewService(repo *Repository, bus events.EventBus, uploader upload.Uploader) Service {
	return &service{repo: repo, bus: bus, uploader: uploader}
}

func (s *service) GetProfile(ctx context.Context, id string) (*Customer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) UpdateProfile(ctx context.Context, id string, req *UpdateCustomerReq) (*Customer, error) {
	if err := s.repo.UpdateProfile(ctx, id, req); err != nil {
		return nil, fmt.Errorf("update profile: %w", err)
	}
	return s.repo.GetByID(ctx, id)
}

func (s *service) UploadAvatar(ctx context.Context, customerID string, file io.Reader, size int64, contentType string) (*Customer, error) {
	if err := upload.ValidateSize(size, maxAvatarSize); err != nil {
		return nil, err
	}

	objectName := fmt.Sprintf("customers/%s/avatar%s", customerID, upload.ExtFromContentType(contentType))

	url, err := s.uploader.Upload(ctx, "avatars", objectName, file, size, contentType)
	if err != nil {
		return nil, fmt.Errorf("upload avatar: %w", err)
	}

	url = fmt.Sprintf("%s?v=%d", url, time.Now().Unix())

	if err := s.repo.UpdateAvatarURL(ctx, customerID, url); err != nil {
		return nil, fmt.Errorf("persist avatar url: %w", err)
	}

	return s.repo.GetByID(ctx, customerID)
}
