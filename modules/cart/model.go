package cart

import "time"

// Cart — CustomerID nullable (guest checkout). CompletedAt set when → order.
type Cart struct {
	ID                string     `db:"id" json:"id"`
	Email             *string    `db:"email" json:"email"`
	CustomerID        *string    `db:"customer_id" json:"customerId"` // soft ref
	CurrencyCode      string     `db:"currency_code" json:"currencyCode"`
	RegionID          *string    `db:"region_id" json:"regionId"`           // soft ref
	SalesChannelID    *string    `db:"sales_channel_id" json:"salesChannelId"` // soft ref
	ShippingAddressID *string    `db:"shipping_address_id" json:"shippingAddressId"`
	BillingAddressID  *string    `db:"billing_address_id" json:"billingAddressId"`
	CompletedAt       *time.Time `db:"completed_at" json:"completedAt"`
	Metadata          []byte     `db:"metadata" json:"-"`
	CreatedAt         time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deletedAt"`

	Items []CartLineItem `db:"-" json:"items,omitempty"`
}

// CartLineItem — VariantID soft ref; product_title/thumbnail snapshot.
type CartLineItem struct {
	ID                 string    `db:"id" json:"id"`
	CartID             string    `db:"cart_id" json:"cartId"`
	VariantID          *string   `db:"variant_id" json:"variantId"` // soft ref
	ProductID          *string   `db:"product_id" json:"productId"` // soft ref
	ProductTitle       *string   `db:"product_title" json:"productTitle"`
	ProductDescription *string   `db:"product_description" json:"productDescription"`
	VariantTitle       *string   `db:"variant_title" json:"variantTitle"`
	VariantSKU         *string   `db:"variant_sku" json:"variantSku"`
	Thumbnail          *string   `db:"thumbnail" json:"thumbnail"`
	Quantity           int       `db:"quantity" json:"quantity"`
	UnitPrice          int64     `db:"unit_price" json:"unitPrice"`
	IsDiscountable     bool      `db:"is_discountable" json:"isDiscountable"`
	IsTaxInclusive     bool      `db:"is_tax_inclusive" json:"isTaxInclusive"`
	Metadata           []byte    `db:"metadata" json:"-"`
	CreatedAt          time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt          time.Time `db:"updated_at" json:"updatedAt"`
}

type CartAddress struct {
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

type CartShippingMethod struct {
	ID               string    `db:"id" json:"id"`
	CartID           string    `db:"cart_id" json:"cartId"`
	ShippingOptionID *string   `db:"shipping_option_id" json:"shippingOptionId"`
	Name             string    `db:"name" json:"name"`
	Amount           int64     `db:"amount" json:"amount"`
	IsTaxInclusive   bool      `db:"is_tax_inclusive" json:"isTaxInclusive"`
	Data             []byte    `db:"data" json:"-"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

// CartAdjustment — applied promotion. Negative amount = discount.
type CartAdjustment struct {
	ID          string    `db:"id" json:"id"`
	CartID      string    `db:"cart_id" json:"cartId"`
	ItemID      *string   `db:"item_id" json:"itemId"`
	PromotionID *string   `db:"promotion_id" json:"promotionId"`
	Code        *string   `db:"code" json:"code"`
	Amount      int64     `db:"amount" json:"amount"`
	Description *string   `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}
