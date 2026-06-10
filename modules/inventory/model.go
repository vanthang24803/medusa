package inventory

import "time"

// InventoryItem — separated from product_variant. requires_shipping=false for digital.
type InventoryItem struct {
	ID               string     `db:"id" json:"id"`
	SKU              *string    `db:"sku" json:"sku"`
	Title            *string    `db:"title" json:"title"`
	Description      *string    `db:"description" json:"description"`
	Thumbnail        *string    `db:"thumbnail" json:"thumbnail"`
	RequiresShipping bool       `db:"requires_shipping" json:"requires_shipping"`
	Weight           *float64   `db:"weight" json:"weight"`
	OriginCountry    *string    `db:"origin_country" json:"origin_country"`
	HsCode           *string    `db:"hs_code" json:"hs_code"`
	Material         *string    `db:"material" json:"material"`
	Metadata         []byte     `db:"metadata" json:"-"`
	CreatedAt        time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at" json:"deleted_at"`
}

// InventoryLevel — stock per location. available = stocked - reserved.
type InventoryLevel struct {
	ID               string    `db:"id" json:"id"`
	InventoryItemID  string    `db:"inventory_item_id" json:"inventory_item_id"`
	LocationID       string    `db:"location_id" json:"location_id"`
	StockedQuantity  int       `db:"stocked_quantity" json:"stocked_quantity"`
	ReservedQuantity int       `db:"reserved_quantity" json:"reserved_quantity"`
	IncomingQuantity int       `db:"incoming_quantity" json:"incoming_quantity"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

func (l InventoryLevel) Available() int {
	return l.StockedQuantity - l.ReservedQuantity
}

// ReservationItem — lock qty when order placed; release when fulfilled/cancelled.
type ReservationItem struct {
	ID              string    `db:"id" json:"id"`
	InventoryItemID string    `db:"inventory_item_id" json:"inventory_item_id"`
	LocationID      string    `db:"location_id" json:"location_id"`
	LineItemID      *string   `db:"line_item_id" json:"line_item_id"`
	Quantity        int       `db:"quantity" json:"quantity"`
	AllowBackorder  bool      `db:"allow_backorder" json:"allow_backorder"`
	Description     *string   `db:"description" json:"description"`
	CreatedBy       *string   `db:"created_by" json:"created_by"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
