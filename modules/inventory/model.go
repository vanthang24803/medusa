package inventory

import "time"

// InventoryItem — separated from product_variant. requires_shipping=false for digital.
type InventoryItem struct {
	ID               string     `db:"id" json:"id"`
	SKU              *string    `db:"sku" json:"sku"`
	Title            *string    `db:"title" json:"title"`
	Description      *string    `db:"description" json:"description"`
	Thumbnail        *string    `db:"thumbnail" json:"thumbnail"`
	RequiresShipping bool       `db:"requires_shipping" json:"requiresShipping"`
	Weight           *float64   `db:"weight" json:"weight"`
	OriginCountry    *string    `db:"origin_country" json:"originCountry"`
	HsCode           *string    `db:"hs_code" json:"hsCode"`
	Material         *string    `db:"material" json:"material"`
	Metadata         []byte     `db:"metadata" json:"-"`
	CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt        *time.Time `db:"deleted_at" json:"deletedAt"`
}

// InventoryLevel — stock per location. available = stocked - reserved.
type InventoryLevel struct {
	ID               string    `db:"id" json:"id"`
	InventoryItemID  string    `db:"inventory_item_id" json:"inventoryItemId"`
	LocationID       string    `db:"location_id" json:"locationId"`
	StockedQuantity  int       `db:"stocked_quantity" json:"stockedQuantity"`
	ReservedQuantity int       `db:"reserved_quantity" json:"reservedQuantity"`
	IncomingQuantity int       `db:"incoming_quantity" json:"incomingQuantity"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

func (l InventoryLevel) Available() int {
	return l.StockedQuantity - l.ReservedQuantity
}

// ReservationItem — lock qty when order placed; release when fulfilled/cancelled.
type ReservationItem struct {
	ID              string    `db:"id" json:"id"`
	InventoryItemID string    `db:"inventory_item_id" json:"inventoryItemId"`
	LocationID      string    `db:"location_id" json:"locationId"`
	LineItemID      *string   `db:"line_item_id" json:"lineItemId"`
	Quantity        int       `db:"quantity" json:"quantity"`
	AllowBackorder  bool      `db:"allow_backorder" json:"allowBackorder"`
	Description     *string   `db:"description" json:"description"`
	CreatedBy       *string   `db:"created_by" json:"createdBy"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`
}
