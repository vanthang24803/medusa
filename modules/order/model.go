package order

import "time"

type OrderStatus string

const (
	OrderStatusPending        OrderStatus = "pending"
	OrderStatusProcessing     OrderStatus = "processing"
	OrderStatusCompleted      OrderStatus = "completed"
	OrderStatusCancelled      OrderStatus = "cancelled"
	OrderStatusRequiresAction OrderStatus = "requires_action"
)

// Order — immutable after creation. CustomerEmail/currency SNAPSHOT. CustomerID soft ref.
type Order struct {
	ID                string      `db:"id" json:"id"`
	DisplayID         int         `db:"display_id" json:"display_id"`
	Status            OrderStatus `db:"status" json:"status"`
	CustomerID        *string     `db:"customer_id" json:"customer_id"` // soft ref
	CustomerEmail     string      `db:"customer_email" json:"customer_email"`
	CurrencyCode      string      `db:"currency_code" json:"currency_code"`
	RegionID          *string     `db:"region_id" json:"region_id"`
	SalesChannelID    *string     `db:"sales_channel_id" json:"sales_channel_id"`
	ShippingAddressID *string     `db:"shipping_address_id" json:"shipping_address_id"`
	BillingAddressID  *string     `db:"billing_address_id" json:"billing_address_id"`
	NoNotification    *bool       `db:"no_notification" json:"no_notification"`
	Metadata          []byte      `db:"metadata" json:"-"`
	CreatedAt         time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time   `db:"updated_at" json:"updated_at"`
	DeletedAt         *time.Time  `db:"deleted_at" json:"deleted_at"`

	Items []OrderLineItem `db:"-" json:"items,omitempty"`
}

// OrderLineItem — complete SNAPSHOT. variant_id/product_id are soft ref only.
type OrderLineItem struct {
	ID                string    `db:"id" json:"id"`
	OrderID           string    `db:"order_id" json:"order_id"`
	VariantID         *string   `db:"variant_id" json:"variant_id"` // soft ref
	ProductID         *string   `db:"product_id" json:"product_id"` // soft ref
	Title             string    `db:"title" json:"title"`
	Subtitle          *string   `db:"subtitle" json:"subtitle"`
	VariantTitle      *string   `db:"variant_title" json:"variant_title"`
	VariantSKU        *string   `db:"variant_sku" json:"variant_sku"`
	Thumbnail         *string   `db:"thumbnail" json:"thumbnail"`
	UnitPrice         int64     `db:"unit_price" json:"unit_price"`
	Quantity          int       `db:"quantity" json:"quantity"`
	FulfilledQuantity int       `db:"fulfilled_quantity" json:"fulfilled_quantity"`
	ReturnedQuantity  int       `db:"returned_quantity" json:"returned_quantity"`
	CancelledQuantity int       `db:"cancelled_quantity" json:"cancelled_quantity"`
	IsDiscountable    bool      `db:"is_discountable" json:"is_discountable"`
	IsTaxInclusive    bool      `db:"is_tax_inclusive" json:"is_tax_inclusive"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

type OrderAddress struct {
	ID          string    `db:"id" json:"id"`
	FirstName   *string   `db:"first_name" json:"first_name"`
	LastName    *string   `db:"last_name" json:"last_name"`
	Company     *string   `db:"company" json:"company"`
	Address1    *string   `db:"address_1" json:"address_1"`
	Address2    *string   `db:"address_2" json:"address_2"`
	City        *string   `db:"city" json:"city"`
	Province    *string   `db:"province" json:"province"`
	PostalCode  *string   `db:"postal_code" json:"postal_code"`
	CountryCode *string   `db:"country_code" json:"country_code"`
	Phone       *string   `db:"phone" json:"phone"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

type OrderShippingMethod struct {
	ID               string    `db:"id" json:"id"`
	OrderID          string    `db:"order_id" json:"order_id"`
	ShippingOptionID *string   `db:"shipping_option_id" json:"shipping_option_id"`
	Name             string    `db:"name" json:"name"`
	Amount           int64     `db:"amount" json:"amount"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type ReturnStatus string

const (
	ReturnStatusRequested ReturnStatus = "requested"
	ReturnStatusReceived  ReturnStatus = "received"
	ReturnStatusRefunded  ReturnStatus = "refunded"
	ReturnStatusCancelled ReturnStatus = "cancelled"
)

type Return struct {
	ID           string       `db:"id" json:"id"`
	OrderID      string       `db:"order_id" json:"order_id"`
	Status       ReturnStatus `db:"status" json:"status"`
	RefundAmount *int64       `db:"refund_amount" json:"refund_amount"`
	Note         *string      `db:"note" json:"note"`
	CreatedAt    time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updated_at"`
}

type ReturnItem struct {
	ID         string  `db:"id" json:"id"`
	ReturnID   string  `db:"return_id" json:"return_id"`
	LineItemID string  `db:"line_item_id" json:"line_item_id"`
	Quantity   int     `db:"quantity" json:"quantity"`
	Note       *string `db:"note" json:"note"`
	ReasonID   *string `db:"reason_id" json:"reason_id"`
}
