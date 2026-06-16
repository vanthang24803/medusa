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

const (
	ReturnStatusRequested ReturnStatus = "requested"
	ReturnStatusReceived  ReturnStatus = "received"
	ReturnStatusRefunded  ReturnStatus = "refunded"
	ReturnStatusCancelled ReturnStatus = "cancelled"
)

// Order — immutable after creation. CustomerEmail/currency SNAPSHOT. CustomerID soft ref.
type Order struct {
	ID                string      `db:"id" json:"id"`
	DisplayID         int         `db:"display_id" json:"displayId"`
	Status            OrderStatus `db:"status" json:"status"`
	CustomerID        *string     `db:"customer_id" json:"customerId"` // soft ref
	CustomerEmail     string      `db:"customer_email" json:"customerEmail"`
	CurrencyCode      string      `db:"currency_code" json:"currencyCode"`
	RegionID          *string     `db:"region_id" json:"regionId"`
	SalesChannelID    *string     `db:"sales_channel_id" json:"salesChannelId"`
	ShippingAddressID *string     `db:"shipping_address_id" json:"shippingAddressId"`
	BillingAddressID  *string     `db:"billing_address_id" json:"billingAddressId"`
	NoNotification    *bool       `db:"no_notification" json:"noNotification"`
	Metadata          []byte      `db:"metadata" json:"-"`
	CreatedAt         time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time   `db:"updated_at" json:"updatedAt"`
	DeletedAt         *time.Time  `db:"deleted_at" json:"deletedAt"`

	Items []OrderLineItem `db:"-" json:"items,omitempty"`
}

// OrderLineItem — complete SNAPSHOT. variant_id/product_id are soft ref only.
type OrderLineItem struct {
	ID                string    `db:"id" json:"id"`
	OrderID           string    `db:"order_id" json:"orderId"`
	VariantID         *string   `db:"variant_id" json:"variantId"` // soft ref
	ProductID         *string   `db:"product_id" json:"productId"` // soft ref
	Title             string    `db:"title" json:"title"`
	Subtitle          *string   `db:"subtitle" json:"subtitle"`
	VariantTitle      *string   `db:"variant_title" json:"variantTitle"`
	VariantSKU        *string   `db:"variant_sku" json:"variantSku"`
	Thumbnail         *string   `db:"thumbnail" json:"thumbnail"`
	UnitPrice         int64     `db:"unit_price" json:"unitPrice"`
	Quantity          int       `db:"quantity" json:"quantity"`
	FulfilledQuantity int       `db:"fulfilled_quantity" json:"fulfilledQuantity"`
	ReturnedQuantity  int       `db:"returned_quantity" json:"returnedQuantity"`
	CancelledQuantity int       `db:"cancelled_quantity" json:"cancelledQuantity"`
	IsDiscountable    bool      `db:"is_discountable" json:"isDiscountable"`
	IsTaxInclusive    bool      `db:"is_tax_inclusive" json:"isTaxInclusive"`
	CreatedAt         time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time `db:"updated_at" json:"updatedAt"`
}

type OrderAddress struct {
	ID          string    `db:"id" json:"id"`
	FirstName   *string   `db:"first_name" json:"firstName"`
	LastName    *string   `db:"last_name" json:"lastName"`
	Company     *string   `db:"company" json:"company"`
	Address1    *string   `db:"address_1" json:"address1"`
	Address2    *string   `db:"address_2" json:"address2"`
	City        *string   `db:"city" json:"city"`
	Province    *string   `db:"province" json:"province"`
	PostalCode  *string   `db:"postal_code" json:"postalCode"`
	CountryCode *string   `db:"country_code" json:"countryCode"`
	Phone       *string   `db:"phone" json:"phone"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type OrderShippingMethod struct {
	ID               string    `db:"id" json:"id"`
	OrderID          string    `db:"order_id" json:"orderId"`
	ShippingOptionID *string   `db:"shipping_option_id" json:"shippingOptionId"`
	Name             string    `db:"name" json:"name"`
	Amount           int64     `db:"amount" json:"amount"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

type ReturnStatus string

type Return struct {
	ID           string       `db:"id" json:"id"`
	OrderID      string       `db:"order_id" json:"orderId"`
	Status       ReturnStatus `db:"status" json:"status"`
	RefundAmount *int64       `db:"refund_amount" json:"refundAmount"`
	Note         *string      `db:"note" json:"note"`
	CreatedAt    time.Time    `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updatedAt"`
}

type ReturnItem struct {
	ID         string  `db:"id" json:"id"`
	ReturnID   string  `db:"return_id" json:"returnId"`
	LineItemID string  `db:"line_item_id" json:"lineItemId"`
	Quantity   int     `db:"quantity" json:"quantity"`
	Note       *string `db:"note" json:"note"`
	ReasonID   *string `db:"reason_id" json:"reasonId"`
}
