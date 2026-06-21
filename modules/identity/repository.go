package identity

import (
	"context"
	"database/sql"
	"errors"

	"ecommerce/packages/db"
	"ecommerce/packages/types"
)

type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

const userColumns = `id, email, first_name, last_name, avatar_url, metadata, created_at, updated_at, deleted_at`

func (r *Repository) Ping(ctx context.Context) error {
	var one int
	return r.db.Reader(ctx).GetContext(ctx, &one, "SELECT 1")
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	var u User
	err := r.db.Reader(ctx).GetContext(ctx, &u,
		`SELECT `+userColumns+` FROM identity.user WHERE id = $1 AND deleted_at IS NULL`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("user")
	}
	return &u, err
}

func (r *Repository) UpdateAvatarURL(ctx context.Context, id, avatarURL string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE identity.user SET avatar_url = $1, updated_at = now() WHERE id = $2 AND deleted_at IS NULL`,
		avatarURL, id)
	return err
}
