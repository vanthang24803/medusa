package customer

import (
	"context"

	"ecommerce/packages/events"
)

type Service interface {
	Ping(ctx context.Context) error
	GetProfile(ctx context.Context, id string) (*Customer, error)
}

type service struct {
	repo *Repository
	bus  events.EventBus
}

func NewService(repo *Repository, bus events.EventBus) Service {
	return &service{repo: repo, bus: bus}
}

func (s *service) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}

func (s *service) GetProfile(ctx context.Context, id string) (*Customer, error) {
	return s.repo.GetByID(ctx, id)
}
