package product

import "time"

type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "draft"
	ProductStatusPublished ProductStatus = "published"
	ProductStatusProposed  ProductStatus = "proposed"
	ProductStatusRejected  ProductStatus = "rejected"
)

// Product — main catalog table. No FK to pricing/inventory.
type Product struct {
	ID            string        `db:"id" json:"id"`
	Title         string        `db:"title" json:"title"`
	Subtitle      *string       `db:"subtitle" json:"subtitle"`
	Description   *string       `db:"description" json:"description"`
	Handle        string        `db:"handle" json:"handle"`
	IsGiftcard    bool          `db:"is_giftcard" json:"is_giftcard"`
	Status        ProductStatus `db:"status" json:"status"`
	Thumbnail     *string       `db:"thumbnail" json:"thumbnail"`
	Weight        *float64      `db:"weight" json:"weight"`
	Length        *float64      `db:"length" json:"length"`
	Height        *float64      `db:"height" json:"height"`
	Width         *float64      `db:"width" json:"width"`
	OriginCountry *string       `db:"origin_country" json:"origin_country"`
	HsCode        *string       `db:"hs_code" json:"hs_code"`
	MidCode       *string       `db:"mid_code" json:"mid_code"`
	Material      *string       `db:"material" json:"material"`
	Discountable  bool          `db:"discountable" json:"discountable"`
	ExternalID    *string       `db:"external_id" json:"external_id"`
	CollectionID  *string       `db:"collection_id" json:"collection_id"`
	TypeID        *string       `db:"type_id" json:"type_id"`
	Metadata      []byte        `db:"metadata" json:"-"`
	CreatedAt     time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time    `db:"deleted_at" json:"deleted_at"`

	// Joined relations (not mapped directly to DB) — loaded separately in service.
	Variants []ProductVariant `db:"-" json:"variants,omitempty"`
	Options  []ProductOption  `db:"-" json:"options,omitempty"`
	Images   []ProductImage   `db:"-" json:"images,omitempty"`
}

// ProductVariant — SKU, barcode, manage_inventory. No price.
type ProductVariant struct {
	ID              string     `db:"id" json:"id"`
	ProductID       string     `db:"product_id" json:"product_id"`
	Title           string     `db:"title" json:"title"`
	SKU             *string    `db:"sku" json:"sku"`
	Barcode         *string    `db:"barcode" json:"barcode"`
	EAN             *string    `db:"ean" json:"ean"`
	UPC             *string    `db:"upc" json:"upc"`
	AllowBackorder  bool       `db:"allow_backorder" json:"allow_backorder"`
	ManageInventory bool       `db:"manage_inventory" json:"manage_inventory"`
	Weight          *float64   `db:"weight" json:"weight"`
	Rank            int        `db:"rank" json:"rank"`
	Metadata        []byte     `db:"metadata" json:"-"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deleted_at"`
}

// ProductOption — "Size", "Color".
type ProductOption struct {
	ID        string    `db:"id" json:"id"`
	ProductID string    `db:"product_id" json:"product_id"`
	Title     string    `db:"title" json:"title"`
	Rank      int       `db:"rank" json:"rank"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	Values []ProductOptionValue `db:"-" json:"values,omitempty"`
}

// ProductOptionValue — "S", "M", "L" / "Red".
type ProductOptionValue struct {
	ID        string    `db:"id" json:"id"`
	OptionID  string    `db:"option_id" json:"option_id"`
	Value     string    `db:"value" json:"value"`
	Rank      int       `db:"rank" json:"rank"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// ProductCategory — tree (self-referencing parent_category_id).
type ProductCategory struct {
	ID               string    `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`
	Description      *string   `db:"description" json:"description"`
	Handle           string    `db:"handle" json:"handle"`
	IsActive         bool      `db:"is_active" json:"is_active"`
	IsInternal       bool      `db:"is_internal" json:"is_internal"`
	Rank             int       `db:"rank" json:"rank"`
	ParentCategoryID *string   `db:"parent_category_id" json:"parent_category_id"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type ProductCollection struct {
	ID        string     `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	Handle    string     `db:"handle" json:"handle"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type ProductImage struct {
	ID        string    `db:"id" json:"id"`
	ProductID string    `db:"product_id" json:"product_id"`
	URL       string    `db:"url" json:"url"`
	Rank      int       `db:"rank" json:"rank"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type ProductTag struct {
	ID        string    `db:"id" json:"id"`
	Value     string    `db:"value" json:"value"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type ProductType struct {
	ID        string    `db:"id" json:"id"`
	Value     string    `db:"value" json:"value"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
