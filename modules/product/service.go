package product

import (
	"context"
	"time"

	"ecommerce/packages/events"
	"ecommerce/packages/types"
)

// ListQuery — filter params cho List.
type ListQuery struct {
	types.PaginationQuery
	Status       string `form:"status"`
	CollectionID string `form:"collection_id"`
	Search       string `form:"q"`
}

// Service interface — contract for product domain, easy to mock in tests.
type Service interface {
	List(ctx context.Context, q ListQuery) (types.PaginatedResponse[Product], error)
	GetByID(ctx context.Context, id string) (*Product, error)
	GetByHandle(ctx context.Context, handle string) (*Product, error)
	Create(ctx context.Context, in CreateInput) (*Product, error)
	Delete(ctx context.Context, id string) error
	ListVariants(ctx context.Context, productID string) ([]ProductVariant, error)
	CreateVariant(ctx context.Context, productID string, in CreateVariantInput) (*ProductVariant, error)
}

// CreateInput — payload for creating a product (validated in handler/dto).
type CreateInput struct {
	Title        string
	Subtitle     *string
	Description  *string
	Handle       string
	Thumbnail    *string
	Status       ProductStatus
	CollectionID *string
}

type CreateVariantInput struct {
	Title           string
	SKU             *string
	ManageInventory bool
	AllowBackorder  bool
}

type service struct {
	repo *Repository
	bus  events.EventBus
}

func NewService(repo *Repository, bus events.EventBus) Service {
	return &service{repo: repo, bus: bus}
}

func (s *service) List(ctx context.Context, q ListQuery) (types.PaginatedResponse[Product], error) {
	q.Normalize()
	items, total, err := s.repo.List(ctx, q)
	if err != nil {
		return types.PaginatedResponse[Product]{}, err
	}
	return types.NewPaginated(items, total, q.Offset(), q.Limit()), nil
}

func (s *service) GetByID(ctx context.Context, id string) (*Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	variants, err := s.repo.ListVariants(ctx, id)
	if err != nil {
		return nil, err
	}
	p.Variants = variants
	return p, nil
}

func (s *service) GetByHandle(ctx context.Context, handle string) (*Product, error) {
	return s.repo.GetByHandle(ctx, handle)
}

func (s *service) Create(ctx context.Context, in CreateInput) (*Product, error) {
	if in.Title == "" {
		return nil, types.NewValidation("title is required")
	}
	if in.Handle == "" {
		return nil, types.NewValidation("handle is required")
	}
	status := in.Status
	if status == "" {
		status = ProductStatusDraft
	}
	now := time.Now().UTC()
	p := &Product{
		ID:           types.GenerateID("prod"),
		Title:        in.Title,
		Subtitle:     in.Subtitle,
		Description:  in.Description,
		Handle:       in.Handle,
		Status:       status,
		Thumbnail:    in.Thumbnail,
		Discountable: true,
		CollectionID: in.CollectionID,
		Metadata:     []byte("{}"),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.repo.Insert(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *service) ListVariants(ctx context.Context, productID string) ([]ProductVariant, error) {
	return s.repo.ListVariants(ctx, productID)
}

func (s *service) CreateVariant(ctx context.Context, productID string, in CreateVariantInput) (*ProductVariant, error) {
	if _, err := s.repo.GetByID(ctx, productID); err != nil {
		return nil, err
	}
	if in.Title == "" {
		return nil, types.NewValidation("variant title is required")
	}
	now := time.Now().UTC()
	v := &ProductVariant{
		ID:              types.GenerateID("variant"),
		ProductID:       productID,
		Title:           in.Title,
		SKU:             in.SKU,
		ManageInventory: in.ManageInventory,
		AllowBackorder:  in.AllowBackorder,
		Metadata:        []byte("{}"),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := s.repo.InsertVariant(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}
