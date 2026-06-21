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
	ID              string        `db:"id" json:"id"`
	Title           string        `db:"title" json:"title"`
	Subtitle        *string       `db:"subtitle" json:"subtitle"`
	Description     *string       `db:"description" json:"description"`
	Handle          string        `db:"handle" json:"handle"`
	IsGiftcard      bool          `db:"is_giftcard" json:"isGiftcard"`
	Status          ProductStatus `db:"status" json:"status"`
	Thumbnail       *string       `db:"thumbnail" json:"thumbnail"`
	Weight          *float64      `db:"weight" json:"weight"`
	Length          *float64      `db:"length" json:"length"`
	Height          *float64      `db:"height" json:"height"`
	Width           *float64      `db:"width" json:"width"`
	OriginCountry   *string       `db:"origin_country" json:"originCountry"`
	HsCode          *string       `db:"hs_code" json:"hsCode"`
	MidCode         *string       `db:"mid_code" json:"midCode"`
	Material        *string       `db:"material" json:"material"`
	Discountable    bool          `db:"discountable" json:"discountable"`
	ExternalID      *string       `db:"external_id" json:"externalId"`
	CollectionID    *string       `db:"collection_id" json:"collectionId"`
	TypeID          *string       `db:"type_id" json:"typeId"`
	BrandID         *string       `db:"brand_id" json:"brandId"`
	Author          *string       `db:"author" json:"author"`
	ISBN            *string       `db:"isbn" json:"isbn"`
	PageCount       *int          `db:"page_count" json:"pageCount"`
	CompareAtPrice  *int64        `db:"compare_at_price" json:"compareAtPrice"`
	Quantity        *int          `db:"quantity" json:"quantity"`
	Rating          *float64      `db:"rating" json:"rating"`
	ReviewCount     int           `db:"review_count" json:"reviewCount"`
	IsFeatured      bool          `db:"is_featured" json:"isFeatured"`
	PublishedAt     *time.Time    `db:"published_at" json:"publishedAt"`
	Metadata        []byte        `db:"metadata" json:"-"`
	LastUpdatedBy   *string       `db:"last_updated_by" json:"lastUpdatedBy"`
	CreatedAt       time.Time     `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time     `db:"updated_at" json:"updatedAt"`
	DeletedAt       *time.Time    `db:"deleted_at" json:"deletedAt"`

	// Joined relations (not mapped directly to DB) — loaded separately in service.
	Variants []ProductVariant `db:"-" json:"variants,omitempty"`
	Options  []ProductOption  `db:"-" json:"options,omitempty"`
	Images   []ProductImage   `db:"-" json:"images,omitempty"`
	Brand    *BrandRef        `db:"-" json:"brand,omitempty"`
}

// BrandRef — lightweight brand snapshot joined from brand.brand.
type BrandRef struct {
	ID      string  `db:"id" json:"id"`
	Name    string  `db:"name" json:"name"`
	Slug    string  `db:"slug" json:"slug"`
	LogoURL *string `db:"logo_url" json:"logoUrl"`
}

// ProductVariant — SKU, barcode, manage_inventory. No price.
type ProductVariant struct {
	ID              string     `db:"id" json:"id"`
	ProductID       string     `db:"product_id" json:"productId"`
	Title           string     `db:"title" json:"title"`
	SKU             *string    `db:"sku" json:"sku"`
	Barcode         *string    `db:"barcode" json:"barcode"`
	EAN             *string    `db:"ean" json:"ean"`
	UPC             *string    `db:"upc" json:"upc"`
	AllowBackorder  bool       `db:"allow_backorder" json:"allowBackorder"`
	ManageInventory bool       `db:"manage_inventory" json:"manageInventory"`
	Weight          *float64   `db:"weight" json:"weight"`
	Rank            int        `db:"rank" json:"rank"`
	Metadata        []byte     `db:"metadata" json:"-"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deletedAt"`
}

// ProductOption — "Size", "Color".
type ProductOption struct {
	ID        string    `db:"id" json:"id"`
	ProductID string    `db:"product_id" json:"productId"`
	Title     string    `db:"title" json:"title"`
	Rank      int       `db:"rank" json:"rank"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	Values []ProductOptionValue `db:"-" json:"values,omitempty"`
}

// ProductOptionValue — "S", "M", "L" / "Red".
type ProductOptionValue struct {
	ID        string    `db:"id" json:"id"`
	OptionID  string    `db:"option_id" json:"optionId"`
	Value     string    `db:"value" json:"value"`
	Rank      int       `db:"rank" json:"rank"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// ProductCategory — tree (self-referencing parent_category_id).
type ProductCategory struct {
	ID               string    `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`
	Description      *string   `db:"description" json:"description"`
	Handle           string    `db:"handle" json:"handle"`
	IsActive         bool      `db:"is_active" json:"isActive"`
	IsInternal       bool      `db:"is_internal" json:"isInternal"`
	Rank             int       `db:"rank" json:"rank"`
	ParentCategoryID *string   `db:"parent_category_id" json:"parentCategoryId"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

type ProductCollection struct {
	ID        string     `db:"id" json:"id"`
	Title     string     `db:"title" json:"title"`
	Handle    string     `db:"handle" json:"handle"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

type ProductImage struct {
	ID        string    `db:"id" json:"id"`
	ProductID string    `db:"product_id" json:"productId"`
	URL       string    `db:"url" json:"url"`
	Rank      int       `db:"rank" json:"rank"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type ProductTag struct {
	ID        string    `db:"id" json:"id"`
	Value     string    `db:"value" json:"value"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type ProductType struct {
	ID        string    `db:"id" json:"id"`
	Value     string    `db:"value" json:"value"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
