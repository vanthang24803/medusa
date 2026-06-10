package region

import "time"

// Store — global config, usually only 1 row.
type Store struct {
	ID                    string    `db:"id" json:"id"`
	Name                  string    `db:"name" json:"name"`
	SupportedCurrencies   []byte    `db:"supported_currencies" json:"-"`
	DefaultSalesChannelID *string   `db:"default_sales_channel_id" json:"default_sales_channel_id"`
	DefaultRegionID       *string   `db:"default_region_id" json:"default_region_id"`
	DefaultLocationID     *string   `db:"default_location_id" json:"default_location_id"`
	CreatedAt             time.Time `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time `db:"updated_at" json:"updated_at"`
}

type SalesChannel struct {
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	IsDisabled  bool       `db:"is_disabled" json:"is_disabled"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at"`
}

// Region — 1 region = 1 currency + set countries.
type Region struct {
	ID             string     `db:"id" json:"id"`
	Name           string     `db:"name" json:"name"`
	CurrencyCode   string     `db:"currency_code" json:"currency_code"`
	AutomaticTaxes bool       `db:"automatic_taxes" json:"automatic_taxes"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deleted_at"`
}

type RegionCountry struct {
	ID          string  `db:"id" json:"id"`
	RegionID    string  `db:"region_id" json:"region_id"`
	Iso2        string  `db:"iso_2" json:"iso_2"`
	Iso3        string  `db:"iso_3" json:"iso_3"`
	Name        string  `db:"name" json:"name"`
	DisplayName string  `db:"display_name" json:"display_name"`
	NumCode     *string `db:"num_code" json:"num_code"`
}

type TaxRate struct {
	ID           string    `db:"id" json:"id"`
	RegionID     string    `db:"region_id" json:"region_id"`
	CountryCode  *string   `db:"country_code" json:"country_code"`
	ProvinceCode *string   `db:"province_code" json:"province_code"`
	Name         string    `db:"name" json:"name"`
	Code         *string   `db:"code" json:"code"`
	Rate         float64   `db:"rate" json:"rate"`
	IsDefault    bool      `db:"is_default" json:"is_default"`
	IsCombinable bool      `db:"is_combinable" json:"is_combinable"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}
