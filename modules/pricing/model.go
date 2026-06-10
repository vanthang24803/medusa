package pricing

import "time"

// PriceSet — pivot. 1 variant links 1 PriceSet, PriceSet has multiple Prices.
type PriceSet struct {
	ID        string    `db:"id" json:"id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`

	Prices []Price `db:"-" json:"prices,omitempty"`
}

// Price — amount in smallest unit (cents). Tier via min/max quantity.
type Price struct {
	ID           string    `db:"id" json:"id"`
	PriceSetID   string    `db:"price_set_id" json:"price_set_id"`
	CurrencyCode string    `db:"currency_code" json:"currency_code"`
	Amount       int64     `db:"amount" json:"amount"`
	MinQuantity  *int      `db:"min_quantity" json:"min_quantity"`
	MaxQuantity  *int      `db:"max_quantity" json:"max_quantity"`
	PriceListID  *string   `db:"price_list_id" json:"price_list_id"`
	RulesCount   int       `db:"rules_count" json:"rules_count"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

type PriceListType string
type PriceListStatus string

const (
	PriceListTypeSale     PriceListType   = "sale"
	PriceListTypeOverride PriceListType   = "override"
	PriceListStatusActive PriceListStatus = "active"
	PriceListStatusDraft  PriceListStatus = "draft"
)

// PriceList — sale (timed discount) / override (B2B).
type PriceList struct {
	ID          string          `db:"id" json:"id"`
	Title       string          `db:"title" json:"title"`
	Description *string         `db:"description" json:"description"`
	Type        PriceListType   `db:"type" json:"type"`
	Status      PriceListStatus `db:"status" json:"status"`
	StartsAt    *time.Time      `db:"starts_at" json:"starts_at"`
	EndsAt      *time.Time      `db:"ends_at" json:"ends_at"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
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
	PriceID   string            `db:"price_id" json:"price_id"`
	Attribute string            `db:"attribute" json:"attribute"`
	Operator  PriceRuleOperator `db:"operator" json:"operator"`
	Value     string            `db:"value" json:"value"`
	CreatedAt time.Time         `db:"created_at" json:"created_at"`
	UpdatedAt time.Time         `db:"updated_at" json:"updated_at"`
}
