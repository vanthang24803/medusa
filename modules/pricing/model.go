package pricing

import "time"

const (
	PriceListTypeSale     PriceListType   = "sale"
	PriceListTypeOverride PriceListType   = "override"
	PriceListStatusActive PriceListStatus = "active"
	PriceListStatusDraft  PriceListStatus = "draft"
)

// PriceSet — pivot. 1 variant links 1 PriceSet, PriceSet has multiple Prices.
type PriceSet struct {
	ID        string    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`

	Prices []Price `db:"-" json:"prices,omitempty"`
}

// Price — amount in smallest unit (cents). Tier via min/max quantity.
type Price struct {
	ID           string    `db:"id" json:"id"`
	PriceSetID   string    `db:"price_set_id" json:"priceSetId"`
	CurrencyCode string    `db:"currency_code" json:"currencyCode"`
	Amount       int64     `db:"amount" json:"amount"`
	MinQuantity  *int      `db:"min_quantity" json:"minQuantity"`
	MaxQuantity  *int      `db:"max_quantity" json:"maxQuantity"`
	PriceListID  *string   `db:"price_list_id" json:"priceListId"`
	RulesCount   int       `db:"rules_count" json:"rulesCount"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

type PriceListType string
type PriceListStatus string

// PriceList — sale (timed discount) / override (B2B).
type PriceList struct {
	ID          string          `db:"id" json:"id"`
	Title       string          `db:"title" json:"title"`
	Description *string         `db:"description" json:"description"`
	Type        PriceListType   `db:"type" json:"type"`
	Status      PriceListStatus `db:"status" json:"status"`
	StartsAt    *time.Time      `db:"starts_at" json:"startsAt"`
	EndsAt      *time.Time      `db:"ends_at" json:"endsAt"`
	CreatedAt   time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updatedAt"`
}

type PriceRuleOperator string

const (
	OperatorEq PriceRuleOperator = "eq"
	OperatorIn PriceRuleOperator = "in"
	OperatorGt PriceRuleOperator = "gt"
	OperatorLt PriceRuleOperator = "lt"
)

// PriceRule — filter price by context (region_id, customer_group_id).
type PriceRule struct {
	ID        string            `db:"id" json:"id"`
	PriceID   string            `db:"price_id" json:"priceId"`
	Attribute string            `db:"attribute" json:"attribute"`
	Operator  PriceRuleOperator `db:"operator" json:"operator"`
	Value     string            `db:"value" json:"value"`
	CreatedAt time.Time         `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time         `db:"updated_at" json:"updatedAt"`
}
