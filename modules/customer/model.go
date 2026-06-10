package customer

import "time"

// Customer — storefront customers. HasAccount=false → guest checkout.
type Customer struct {
	ID          string     `db:"id" json:"id"`
	Email       string     `db:"email" json:"email"`
	FirstName   *string    `db:"first_name" json:"first_name"`
	LastName    *string    `db:"last_name" json:"last_name"`
	Phone       *string    `db:"phone" json:"phone"`
	CompanyName *string    `db:"company_name" json:"company_name"`
	HasAccount  bool       `db:"has_account" json:"has_account"`
	Metadata    []byte     `db:"metadata" json:"-"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}

// CustomerAddress — multiple addresses per customer.
type CustomerAddress struct {
	ID                string    `db:"id" json:"id"`
	CustomerID        string    `db:"customer_id" json:"customer_id"`
	AddressName       *string   `db:"address_name" json:"address_name"`
	FirstName         *string   `db:"first_name" json:"first_name"`
	LastName          *string   `db:"last_name" json:"last_name"`
	Company           *string   `db:"company" json:"company"`
	Address1          *string   `db:"address_1" json:"address_1"`
	Address2          *string   `db:"address_2" json:"address_2"`
	City              *string   `db:"city" json:"city"`
	Province          *string   `db:"province" json:"province"`
	PostalCode        *string   `db:"postal_code" json:"postal_code"`
	CountryCode       *string   `db:"country_code" json:"country_code"`
	Phone             *string   `db:"phone" json:"phone"`
	IsDefaultShipping bool      `db:"is_default_shipping" json:"is_default_shipping"`
	IsDefaultBilling  bool      `db:"is_default_billing" json:"is_default_billing"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// CustomerGroup — grouping for pricing rules / promotions.
type CustomerGroup struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Metadata  []byte    `db:"metadata" json:"-"`
	CreatedBy *string   `db:"created_by" json:"created_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
