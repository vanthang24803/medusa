package promotion

import (
	"context"

	"ecommerce/packages/db"
)

// Repository — data access for promotion schema (handwritten SQL with sqlx).
// STUB: add Insert/Get/List/Update/Delete when fully implemented,
// following the pattern in modules/product/repository.go.
type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

// Ping — verify connection reachable.
func (r *Repository) Ping(ctx context.Context) error {
	var one int
	return r.db.Reader(ctx).GetContext(ctx, &one, "SELECT 1")
}
