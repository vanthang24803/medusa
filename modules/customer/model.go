package customer

import "time"

// Customer — storefront customers. HasAccount=false → guest checkout.
type Customer struct {
	ID          string     `db:"id" json:"id"`
	Email       string     `db:"email" json:"email"`
	FirstName   *string    `db:"first_name" json:"firstName"`
	LastName    *string    `db:"last_name" json:"lastName"`
	Phone       *string    `db:"phone" json:"phone"`
	CompanyName *string    `db:"company_name" json:"companyName"`
	AvatarURL   *string    `db:"avatar_url" json:"avatarUrl"`
	HasAccount      bool       `db:"has_account" json:"hasAccount"`
	Metadata        []byte     `db:"metadata" json:"-"`
	LastUpdatedBy   *string    `db:"last_updated_by" json:"lastUpdatedBy"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deletedAt"`
}

// CustomerAddress — multiple addresses per customer.
type CustomerAddress struct {
	ID                string    `db:"id" json:"id"`
	CustomerID        string    `db:"customer_id" json:"customerId"`
	AddressName       *string   `db:"address_name" json:"addressName"`
	FirstName         *string   `db:"first_name" json:"firstName"`
	LastName          *string   `db:"last_name" json:"lastName"`
	Company           *string   `db:"company" json:"company"`
	Address1          *string   `db:"address_1" json:"address1"`
	Address2          *string   `db:"address_2" json:"address2"`
	City              *string   `db:"city" json:"city"`
	Province          *string   `db:"province" json:"province"`
	PostalCode        *string   `db:"postal_code" json:"postalCode"`
	CountryCode       *string   `db:"country_code" json:"countryCode"`
	Phone             *string   `db:"phone" json:"phone"`
	IsDefaultShipping bool      `db:"is_default_shipping" json:"isDefaultShipping"`
	IsDefaultBilling  bool      `db:"is_default_billing" json:"isDefaultBilling"`
	CreatedAt         time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time `db:"updated_at" json:"updatedAt"`
}

// CustomerGroup — grouping for pricing rules / promotions.
type CustomerGroup struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Metadata  []byte    `db:"metadata" json:"-"`
	CreatedBy *string   `db:"created_by" json:"createdBy"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

type UpdateCustomerReq struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Phone     *string `json:"phone"`
	Metadata  []byte  `json:"metadata"`
}
