package brand

import "time"

// Brand — publisher / vendor / NXB for products.
type Brand struct {
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Slug        string     `db:"slug" json:"slug"`
	LogoURL     *string    `db:"logo_url" json:"logoUrl"`
	Description *string    `db:"description" json:"description"`
	IsActive    bool       `db:"is_active" json:"isActive"`
	Rank        int        `db:"rank" json:"rank"`
	Metadata    []byte     `db:"metadata" json:"-"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
}

type ListQuery struct {
	Search string `form:"q"`
	Status string `form:"status"`
}

type CreateInput struct {
	Name        string
	Slug        string
	LogoURL     *string
	Description *string
}

type UpdateInput struct {
	Name        *string
	Slug        *string
	LogoURL     *string
	Description *string
	IsActive    *bool
	Rank        *int
}
