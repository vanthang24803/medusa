package product

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ecommerce/packages/db"
	"ecommerce/packages/types"
)

// Repository — data access for product schema, handwritten SQL with sqlx.
type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

const productColumns = `id, title, subtitle, description, handle, is_giftcard, status,
	thumbnail, weight, length, height, width, origin_country, hs_code, mid_code,
	material, discountable, external_id, collection_id, type_id, brand_id,
	author, isbn, page_count, compare_at_price, quantity, rating, review_count,
	is_featured, published_at, metadata, last_updated_by, created_at, updated_at, deleted_at`

// Insert creates a new product.
func (r *Repository) Insert(ctx context.Context, p *Product) error {
	query := fmt.Sprintf(`
		INSERT INTO product.product (%s)
		VALUES (:id, :title, :subtitle, :description, :handle, :is_giftcard, :status,
			:thumbnail, :weight, :length, :height, :width, :origin_country, :hs_code,
			:mid_code, :material, :discountable, :external_id, :collection_id, :type_id,
			:brand_id, :author, :isbn, :page_count, :compare_at_price, :quantity,
			:rating, :review_count, :is_featured, :published_at,
			:metadata, :last_updated_by, :created_at, :updated_at, :deleted_at)`, productColumns)
	_, err := r.db.Writer(ctx).NamedExec(query, p)
	return err
}

// GetByID reads 1 product (variants not loaded).
func (r *Repository) GetByID(ctx context.Context, id string) (*Product, error) {
	var p Product
	query := fmt.Sprintf(
		`SELECT %s FROM product.product WHERE id = $1 AND deleted_at IS NULL`,
		productColumns)
	err := r.db.Reader(ctx).GetContext(ctx, &p, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("product")
	}
	return &p, err
}

// GetByHandle reads product by handle (used for storefront URL).
func (r *Repository) GetByHandle(ctx context.Context, handle string) (*Product, error) {
	var p Product
	query := fmt.Sprintf(
		`SELECT %s FROM product.product WHERE handle = $1 AND deleted_at IS NULL`,
		productColumns)
	err := r.db.Reader(ctx).GetContext(ctx, &p, query, handle)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("product")
	}
	return &p, err
}

// List returns products with filtering + pagination.
func (r *Repository) List(ctx context.Context, q ListQuery) ([]Product, int, error) {
	where := `WHERE deleted_at IS NULL`
	args := map[string]any{}
	if q.Status != "" {
		where += ` AND status = :status`
		args["status"] = q.Status
	}
	if q.CollectionID != "" {
		where += ` AND collection_id = :collection_id`
		args["collection_id"] = q.CollectionID
	}
	if q.TypeID != "" {
		where += ` AND type_id = :type_id`
		args["type_id"] = q.TypeID
	}
	if q.BrandID != "" {
		where += ` AND brand_id = :brand_id`
		args["brand_id"] = q.BrandID
	}
	if q.IsFeatured != nil {
		where += ` AND is_featured = :is_featured`
		args["is_featured"] = *q.IsFeatured
	}
	if q.Search != "" {
		where += ` AND (title ILIKE :search OR author ILIKE :search)`
		args["search"] = "%" + q.Search + "%"
	}

	// Sort
	orderBy := `ORDER BY created_at DESC`
	switch q.Sort {
	case "price_asc":
		orderBy = `ORDER BY compare_at_price ASC NULLS LAST`
	case "price_desc":
		orderBy = `ORDER BY compare_at_price DESC NULLS LAST`
	case "rating":
		orderBy = `ORDER BY rating DESC NULLS LAST`
	case "newest":
		orderBy = `ORDER BY published_at DESC NULLS LAST`
	case "title":
		orderBy = `ORDER BY title ASC`
	}

	// Count
	var total int
	countQ, countArgs, _ := r.db.Reader(ctx).BindNamed(
		`SELECT COUNT(*) FROM product.product `+where, args)
	if err := r.db.Reader(ctx).GetContext(ctx, &total, countQ, countArgs...); err != nil {
		return nil, 0, err
	}

	// Page
	args["limit"] = q.Limit()
	args["offset"] = q.Offset()
	listQ := fmt.Sprintf(`SELECT %s FROM product.product %s
		%s LIMIT :limit OFFSET :offset`, productColumns, where, orderBy)
	bound, bargs, _ := r.db.Reader(ctx).BindNamed(listQ, args)

	products := []Product{}
	if err := r.db.Reader(ctx).SelectContext(ctx, &products, bound, bargs...); err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

// SoftDelete marks deleted_at.
func (r *Repository) SoftDelete(ctx context.Context, id string) error {
	res, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE product.product SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return types.NewNotFound("product")
	}
	return nil
}

// ── Variants ────────────────────────────────────────────────────────────────

func (r *Repository) InsertVariant(ctx context.Context, v *ProductVariant) error {
	query := `INSERT INTO product.product_variant
		(id, product_id, title, sku, barcode, ean, upc, allow_backorder,
		 manage_inventory, weight, rank, metadata, created_at, updated_at)
		VALUES (:id, :product_id, :title, :sku, :barcode, :ean, :upc, :allow_backorder,
		 :manage_inventory, :weight, :rank, :metadata, :created_at, :updated_at)`
	_, err := r.db.Writer(ctx).NamedExec(query, v)
	return err
}

func (r *Repository) ListVariants(ctx context.Context, productID string) ([]ProductVariant, error) {
	variants := []ProductVariant{}
	err := r.db.Reader(ctx).SelectContext(ctx, &variants,
		`SELECT id, product_id, title, sku, barcode, ean, upc, allow_backorder,
		 manage_inventory, weight, rank, metadata, created_at, updated_at, deleted_at
		 FROM product.product_variant
		 WHERE product_id = $1 AND deleted_at IS NULL ORDER BY rank`, productID)
	return variants, err
}

// ── Brand join ──────────────────────────────────────────────────────────────

func (r *Repository) GetBrand(ctx context.Context, brandID string) (*BrandRef, error) {
	var b BrandRef
	err := r.db.Reader(ctx).GetContext(ctx, &b,
		`SELECT id, name, slug, logo_url FROM brand.brand WHERE id = $1 AND deleted_at IS NULL`, brandID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &b, err
}

func (r *Repository) ListImages(ctx context.Context, productID string) ([]ProductImage, error) {
	var imgs []ProductImage
	err := r.db.Reader(ctx).SelectContext(ctx, &imgs,
		`SELECT id, product_id, url, rank, created_at, updated_at
		 FROM product.product_image
		 WHERE product_id = $1 ORDER BY rank`, productID)
	return imgs, err
}

func (r *Repository) ListOptions(ctx context.Context, productID string) ([]ProductOption, error) {
	var opts []ProductOption
	err := r.db.Reader(ctx).SelectContext(ctx, &opts,
		`SELECT id, product_id, title, rank, created_at, updated_at
		 FROM product.product_option
		 WHERE product_id = $1 ORDER BY rank`, productID)
	return opts, err
}
