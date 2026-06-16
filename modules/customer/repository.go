package customer

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

const customerColumns = `id, email, first_name, last_name, phone, company_name,
	has_account, metadata, created_at, updated_at, deleted_at`

func (r *Repository) Ping(ctx context.Context) error {
	var one int
	return r.db.Reader(ctx).GetContext(ctx, &one, "SELECT 1")
}

func (r *Repository) Insert(ctx context.Context, c *Customer) error {
	query := `INSERT INTO customer.customer (` + customerColumns + `)
		VALUES (:id, :email, :first_name, :last_name, :phone, :company_name,
			:has_account, :metadata, :created_at, :updated_at, :deleted_at)`
	_, err := r.db.Writer(ctx).NamedExec(query, c)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Customer, error) {
	var c Customer
	err := r.db.Reader(ctx).GetContext(ctx, &c,
		`SELECT `+customerColumns+` FROM customer.customer WHERE id = $1 AND deleted_at IS NULL`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("customer")
	}
	return &c, err
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*Customer, error) {
	var c Customer
	err := r.db.Reader(ctx).GetContext(ctx, &c,
		`SELECT `+customerColumns+` FROM customer.customer WHERE email = $1 AND deleted_at IS NULL`, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("customer")
	}
	return &c, err
}
