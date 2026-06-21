package brand

import (
	"context"
	"time"

	"ecommerce/packages/actor"
	"ecommerce/packages/events"
	"ecommerce/packages/types"
)

type Service interface {
	List(ctx context.Context, q ListQuery) ([]Brand, error)
	GetByID(ctx context.Context, id string) (*Brand, error)
	GetBySlug(ctx context.Context, slug string) (*Brand, error)
	Create(ctx context.Context, in CreateInput) (*Brand, error)
	Update(ctx context.Context, id string, in UpdateInput) (*Brand, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo *Repository
	bus  events.EventBus
}

func NewService(repo *Repository, bus events.EventBus) Service {
	return &service{repo: repo, bus: bus}
}

func (s *service) List(ctx context.Context, q ListQuery) ([]Brand, error) {
	return s.repo.List(ctx, q)
}

func (s *service) GetByID(ctx context.Context, id string) (*Brand, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetBySlug(ctx context.Context, slug string) (*Brand, error) {
	return s.repo.GetBySlug(ctx, slug)
}

func (s *service) Create(ctx context.Context, in CreateInput) (*Brand, error) {
	if in.Name == "" {
		return nil, types.NewValidation("name is required")
	}
	if in.Slug == "" {
		return nil, types.NewValidation("slug is required")
	}
	now := time.Now().UTC()
	b := &Brand{
		ID:          types.GenerateID("brand"),
		Name:        in.Name,
		Slug:        in.Slug,
		LogoURL:     in.LogoURL,
		Description: in.Description,
		IsActive:    true,
		Metadata:    []byte("{}"),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.Insert(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *service) Update(ctx context.Context, id string, in UpdateInput) (*Brand, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if in.Name != nil {
		b.Name = *in.Name
	}
	if in.Slug != nil {
		b.Slug = *in.Slug
	}
	if in.LogoURL != nil {
		b.LogoURL = in.LogoURL
	}
	if in.Description != nil {
		b.Description = in.Description
	}
	if in.IsActive != nil {
		b.IsActive = *in.IsActive
	}
	if in.Rank != nil {
		b.Rank = *in.Rank
	}
	b.UpdatedAt = time.Now().UTC()
	b.LastUpdatedBy = actor.Get(ctx)
	if err := s.repo.Update(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}
