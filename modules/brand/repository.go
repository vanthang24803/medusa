package brand

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ecommerce/packages/db"
	"ecommerce/packages/types"
)

type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

const brandColumns = `id, name, slug, logo_url, description, is_active, rank,
	metadata, created_at, updated_at, deleted_at`

func (r *Repository) Insert(ctx context.Context, b *Brand) error {
	query := fmt.Sprintf(`
		INSERT INTO brand.brand (%s)
		VALUES (:id, :name, :slug, :logo_url, :description, :is_active, :rank,
			:metadata, :created_at, :updated_at, :deleted_at)`, brandColumns)
	_, err := r.db.Writer(ctx).NamedExec(query, b)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Brand, error) {
	var b Brand
	query := fmt.Sprintf(
		`SELECT %s FROM brand.brand WHERE id = $1 AND deleted_at IS NULL`,
		brandColumns)
	err := r.db.Reader(ctx).GetContext(ctx, &b, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("brand")
	}
	return &b, err
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*Brand, error) {
	var b Brand
	query := fmt.Sprintf(
		`SELECT %s FROM brand.brand WHERE slug = $1 AND deleted_at IS NULL`,
		brandColumns)
	err := r.db.Reader(ctx).GetContext(ctx, &b, query, slug)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("brand")
	}
	return &b, err
}

func (r *Repository) List(ctx context.Context, q ListQuery) ([]Brand, error) {
	where := `WHERE deleted_at IS NULL`
	args := map[string]any{}
	if q.Search != "" {
		where += ` AND name ILIKE :search`
		args["search"] = "%" + q.Search + "%"
	}
	if q.Status == "active" {
		where += ` AND is_active = true`
	} else if q.Status == "inactive" {
		where += ` AND is_active = false`
	}

	query := fmt.Sprintf(`SELECT %s FROM brand.brand %s ORDER BY rank, name`, brandColumns, where)
	bound, bargs, _ := r.db.Reader(ctx).BindNamed(query, args)

	var brands []Brand
	if err := r.db.Reader(ctx).SelectContext(ctx, &brands, bound, bargs...); err != nil {
		return nil, err
	}
	return brands, nil
}

func (r *Repository) Update(ctx context.Context, b *Brand) error {
	query := `UPDATE brand.brand SET
		name = :name, slug = :slug, logo_url = :logo_url, description = :description,
		is_active = :is_active, rank = :rank, metadata = :metadata, updated_at = :updated_at
		WHERE id = :id AND deleted_at IS NULL`
	_, err := r.db.Writer(ctx).NamedExec(query, b)
	return err
}

func (r *Repository) SoftDelete(ctx context.Context, id string) error {
	res, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE brand.brand SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return types.NewNotFound("brand")
	}
	return nil
}
