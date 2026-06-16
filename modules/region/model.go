package region

import "time"

// Store — global config, usually only 1 row.
type Store struct {
	ID                    string    `db:"id" json:"id"`
	Name                  string    `db:"name" json:"name"`
	SupportedCurrencies   []byte    `db:"supported_currencies" json:"-"`
	DefaultSalesChannelID *string   `db:"default_sales_channel_id" json:"defaultSalesChannelId"`
	DefaultRegionID       *string   `db:"default_region_id" json:"defaultRegionId"`
	DefaultLocationID     *string   `db:"default_location_id" json:"defaultLocationId"`
	CreatedAt             time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt             time.Time `db:"updated_at" json:"updatedAt"`
}

type SalesChannel struct {
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	IsDisabled  bool       `db:"is_disabled" json:"isDisabled"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
}

// Region — 1 region = 1 currency + set countries.
type Region struct {
	ID             string     `db:"id" json:"id"`
	Name           string     `db:"name" json:"name"`
	CurrencyCode   string     `db:"currency_code" json:"currencyCode"`
	AutomaticTaxes bool       `db:"automatic_taxes" json:"automaticTaxes"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deletedAt"`
}

type RegionCountry struct {
	ID          string  `db:"id" json:"id"`
	RegionID    string  `db:"region_id" json:"regionId"`
	Iso2        string  `db:"iso_2" json:"iso2"`
	Iso3        string  `db:"iso_3" json:"iso3"`
	Name        string  `db:"name" json:"name"`
	DisplayName string  `db:"display_name" json:"displayName"`
	NumCode     *string `db:"num_code" json:"numCode"`
}

type TaxRate struct {
	ID           string    `db:"id" json:"id"`
	RegionID     string    `db:"region_id" json:"regionId"`
	CountryCode  *string   `db:"country_code" json:"countryCode"`
	ProvinceCode *string   `db:"province_code" json:"provinceCode"`
	Name         string    `db:"name" json:"name"`
	Code         *string   `db:"code" json:"code"`
	Rate         float64   `db:"rate" json:"rate"`
	IsDefault    bool      `db:"is_default" json:"isDefault"`
	IsCombinable bool      `db:"is_combinable" json:"isCombinable"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}
