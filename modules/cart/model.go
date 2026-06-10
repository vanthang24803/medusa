package cart

import "time"

// Cart — CustomerID nullable (guest checkout). CompletedAt set when → order.
type Cart struct {
	ID                string     `db:"id" json:"id"`
	Email             *string    `db:"email" json:"email"`
	CustomerID        *string    `db:"customer_id" json:"customer_id"` // soft ref
	CurrencyCode      string     `db:"currency_code" json:"currency_code"`
	RegionID          *string    `db:"region_id" json:"region_id"`           // soft ref
	SalesChannelID    *string    `db:"sales_channel_id" json:"sales_channel_id"` // soft ref
	ShippingAddressID *string    `db:"shipping_address_id" json:"shipping_address_id"`
	BillingAddressID  *string    `db:"billing_address_id" json:"billing_address_id"`
	CompletedAt       *time.Time `db:"completed_at" json:"completed_at"`
	Metadata          []byte     `db:"metadata" json:"-"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deleted_at"`

	Items []CartLineItem `db:"-" json:"items,omitempty"`
}

// CartLineItem — VariantID soft ref; product_title/thumbnail snapshot.
type CartLineItem struct {
	ID                 string    `db:"id" json:"id"`
	CartID             string    `db:"cart_id" json:"cart_id"`
	VariantID          *string   `db:"variant_id" json:"variant_id"` // soft ref
	ProductID          *string   `db:"product_id" json:"product_id"` // soft ref
	ProductTitle       *string   `db:"product_title" json:"product_title"`
	ProductDescription *string   `db:"product_description" json:"product_description"`
	VariantTitle       *string   `db:"variant_title" json:"variant_title"`
	VariantSKU         *string   `db:"variant_sku" json:"variant_sku"`
	Thumbnail          *string   `db:"thumbnail" json:"thumbnail"`
	Quantity           int       `db:"quantity" json:"quantity"`
	UnitPrice          int64     `db:"unit_price" json:"unit_price"`
	IsDiscountable     bool      `db:"is_discountable" json:"is_discountable"`
	IsTaxInclusive     bool      `db:"is_tax_inclusive" json:"is_tax_inclusive"`
	Metadata           []byte    `db:"metadata" json:"-"`
	CreatedAt          time.Time `db:"created_at" json:"created_at"`
	UpdatedAt          time.Time `db:"updated_at" json:"updated_at"`
}

type CartAddress struct {
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

type CartShippingMethod struct {
	ID               string    `db:"id" json:"id"`
	CartID           string    `db:"cart_id" json:"cart_id"`
	ShippingOptionID *string   `db:"shipping_option_id" json:"shipping_option_id"`
	Name             string    `db:"name" json:"name"`
	Amount           int64     `db:"amount" json:"amount"`
	IsTaxInclusive   bool      `db:"is_tax_inclusive" json:"is_tax_inclusive"`
	Data             []byte    `db:"data" json:"-"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// CartAdjustment — applied promotion. Negative amount = discount.
type CartAdjustment struct {
	ID          string    `db:"id" json:"id"`
	CartID      string    `db:"cart_id" json:"cart_id"`
	ItemID      *string   `db:"item_id" json:"item_id"`
	PromotionID *string   `db:"promotion_id" json:"promotion_id"`
	Code        *string   `db:"code" json:"code"`
	Amount      int64     `db:"amount" json:"amount"`
	Description *string   `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
