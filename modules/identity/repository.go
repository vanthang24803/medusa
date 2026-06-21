package identity

import (
	"context"
	"database/sql"
	"errors"

	"ecommerce/packages/actor"
	"ecommerce/packages/db"
	"ecommerce/packages/types"
)

type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

const userColumns = `id, email, first_name, last_name, avatar_url, status, metadata, last_updated_by, created_at, updated_at, deleted_at`

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

func (r *Repository) ListUsers(ctx context.Context) ([]User, error) {
	var users []User
	err := r.db.Reader(ctx).SelectContext(ctx, &users,
		`SELECT `+userColumns+` FROM identity.user WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	return users, err
}

func (r *Repository) SoftDeleteUser(ctx context.Context, id string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE identity.user SET deleted_at = now(), updated_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	return err
}

func (r *Repository) InsertUser(ctx context.Context, u *User) error {
	_, err := r.db.Writer(ctx).NamedExec(`
		INSERT INTO identity.user (id, email, first_name, last_name, avatar_url, status, metadata, last_updated_by, created_at, updated_at)
		VALUES (:id, :email, :first_name, :last_name, :avatar_url, :status, :metadata, :last_updated_by, :created_at, :updated_at)`, u)
	return err
}

func (r *Repository) SetStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE identity.user SET status = $1, last_updated_by = $2, updated_at = now() WHERE id = $3 AND deleted_at IS NULL`,
		status, actor.Get(ctx), id)
	return err
}

func (r *Repository) UpdateProfile(ctx context.Context, id string, req *UpdateProfileReq) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE identity.user SET first_name = COALESCE($1, first_name), last_name = COALESCE($2, last_name), last_updated_by = $3, updated_at = now() WHERE id = $4 AND deleted_at IS NULL`,
		req.FirstName, req.LastName, actor.Get(ctx), id)
	return err
}

func (r *Repository) UpdateAvatarURL(ctx context.Context, id, avatarURL string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE identity.user SET avatar_url = $1, last_updated_by = $2, updated_at = now() WHERE id = $3 AND deleted_at IS NULL`,
		avatarURL, actor.Get(ctx), id)
	return err
}
