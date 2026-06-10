package customer

import (
	"context"

	"ecommerce/packages/events"
)

// Service — contract cho customer domain.
// STUB: skeleton structure; add methods when fully implemented.
type Service interface {
	Ping(ctx context.Context) error
}

type service struct {
	repo *Repository
	bus  events.EventBus
}

func NewService(repo *Repository, bus events.EventBus) Service {
	return &service{repo: repo, bus: bus}
}

// Ping — placeholder to verify wiring; replace with real business methods.
func (s *service) Ping(ctx context.Context) error {
	return s.repo.Ping(ctx)
}
