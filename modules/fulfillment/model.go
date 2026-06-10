package fulfillment

import "time"

type StockLocation struct {
	ID        string     `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	AddressID *string    `db:"address_id" json:"address_id"`
	Metadata  []byte     `db:"metadata" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type ShippingProfile struct {
	ID        string     `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Type      string     `db:"type" json:"type"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

type ShippingPriceType string

const (
	ShippingPriceTypeFlat       ShippingPriceType = "flat"
	ShippingPriceTypeCalculated ShippingPriceType = "calculated"
)

type ShippingOption struct {
	ID                string            `db:"id" json:"id"`
	Name              string            `db:"name" json:"name"`
	ProviderID        string            `db:"provider_id" json:"provider_id"`
	ServiceZoneID     string            `db:"service_zone_id" json:"service_zone_id"`
	ShippingProfileID *string           `db:"shipping_profile_id" json:"shipping_profile_id"`
	PriceType         ShippingPriceType `db:"price_type" json:"price_type"`
	Data              []byte            `db:"data" json:"-"`
	Metadata          []byte            `db:"metadata" json:"-"`
	CreatedAt         time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time         `db:"updated_at" json:"updated_at"`
	DeletedAt         *time.Time        `db:"deleted_at" json:"deleted_at"`
}

// Fulfillment — OrderID/LocationID soft ref.
type Fulfillment struct {
	ID             string     `db:"id" json:"id"`
	OrderID        *string    `db:"order_id" json:"order_id"`       // soft ref
	LocationID     *string    `db:"location_id" json:"location_id"` // soft ref
	ProviderID     string     `db:"provider_id" json:"provider_id"`
	TrackingNumber *string    `db:"tracking_number" json:"tracking_number"`
	TrackingURL    *string    `db:"tracking_url" json:"tracking_url"`
	ShippedAt      *time.Time `db:"shipped_at" json:"shipped_at"`
	DeliveredAt    *time.Time `db:"delivered_at" json:"delivered_at"`
	CanceledAt     *time.Time `db:"canceled_at" json:"canceled_at"`
	Data           []byte     `db:"data" json:"-"`
	Metadata       []byte     `db:"metadata" json:"-"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
}

// FulfillmentItem — line_item_id soft ref; title/sku snapshot.
type FulfillmentItem struct {
	ID              string    `db:"id" json:"id"`
	FulfillmentID   string    `db:"fulfillment_id" json:"fulfillment_id"`
	LineItemID      *string   `db:"line_item_id" json:"line_item_id"`
	InventoryItemID *string   `db:"inventory_item_id" json:"inventory_item_id"`
	Title           *string   `db:"title" json:"title"`
	SKU             *string   `db:"sku" json:"sku"`
	Barcode         *string   `db:"barcode" json:"barcode"`
	Quantity        int       `db:"quantity" json:"quantity"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}
