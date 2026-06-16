package fulfillment

import "time"

type StockLocation struct {
	ID        string     `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	AddressID *string    `db:"address_id" json:"addressId"`
	Metadata  []byte     `db:"metadata" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

type ShippingProfile struct {
	ID        string     `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	Type      string     `db:"type" json:"type"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

type ShippingPriceType string

const (
	ShippingPriceTypeFlat       ShippingPriceType = "flat"
	ShippingPriceTypeCalculated ShippingPriceType = "calculated"
)

type ShippingOption struct {
	ID                string            `db:"id" json:"id"`
	Name              string            `db:"name" json:"name"`
	ProviderID        string            `db:"provider_id" json:"providerId"`
	ServiceZoneID     string            `db:"service_zone_id" json:"serviceZoneId"`
	ShippingProfileID *string           `db:"shipping_profile_id" json:"shippingProfileId"`
	PriceType         ShippingPriceType `db:"price_type" json:"priceType"`
	Data              []byte            `db:"data" json:"-"`
	Metadata          []byte            `db:"metadata" json:"-"`
	CreatedAt         time.Time         `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time         `db:"updated_at" json:"updatedAt"`
	DeletedAt         *time.Time        `db:"deleted_at" json:"deletedAt"`
}

// Fulfillment — OrderID/LocationID soft ref.
type Fulfillment struct {
	ID             string     `db:"id" json:"id"`
	OrderID        *string    `db:"order_id" json:"orderId"`       // soft ref
	LocationID     *string    `db:"location_id" json:"locationId"` // soft ref
	ProviderID     string     `db:"provider_id" json:"providerId"`
	TrackingNumber *string    `db:"tracking_number" json:"trackingNumber"`
	TrackingURL    *string    `db:"tracking_url" json:"trackingUrl"`
	ShippedAt      *time.Time `db:"shipped_at" json:"shippedAt"`
	DeliveredAt    *time.Time `db:"delivered_at" json:"deliveredAt"`
	CanceledAt     *time.Time `db:"canceled_at" json:"canceledAt"`
	Data           []byte     `db:"data" json:"-"`
	Metadata       []byte     `db:"metadata" json:"-"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updatedAt"`
}

// FulfillmentItem — line_item_id soft ref; title/sku snapshot.
type FulfillmentItem struct {
	ID              string    `db:"id" json:"id"`
	FulfillmentID   string    `db:"fulfillment_id" json:"fulfillmentId"`
	LineItemID      *string   `db:"line_item_id" json:"lineItemId"`
	InventoryItemID *string   `db:"inventory_item_id" json:"inventoryItemId"`
	Title           *string   `db:"title" json:"title"`
	SKU             *string   `db:"sku" json:"sku"`
	Barcode         *string   `db:"barcode" json:"barcode"`
	Quantity        int       `db:"quantity" json:"quantity"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`
}
