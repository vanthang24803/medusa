package promotion

import "time"

type PromotionType string
type PromotionStatus string
type RuleOperator string

const (
	PromotionTypeStandard PromotionType = "standard"
	PromotionTypeBuyGet   PromotionType = "buyget"

	PromotionStatusActive   PromotionStatus = "active"
	PromotionStatusInactive PromotionStatus = "inactive"
	PromotionStatusExpired  PromotionStatus = "expired"

	RuleOperatorEq  RuleOperator = "eq"
	RuleOperatorIn  RuleOperator = "in"
	RuleOperatorGt  RuleOperator = "gt"
	RuleOperatorLt  RuleOperator = "lt"
	RuleOperatorGte RuleOperator = "gte"
	RuleOperatorLte RuleOperator = "lte"
)

type Promotion struct {
	ID          string          `db:"id" json:"id"`
	Code        string          `db:"code" json:"code"`
	Type        PromotionType   `db:"type" json:"type"`
	IsAutomatic bool            `db:"is_automatic" json:"isAutomatic"`
	Status      PromotionStatus `db:"status" json:"status"`
	CampaignID  *string         `db:"campaign_id" json:"campaignId"`
	CreatedAt   time.Time       `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updatedAt"`
}

type PromotionRule struct {
	ID          string       `db:"id" json:"id"`
	PromotionID string       `db:"promotion_id" json:"promotionId"`
	Attribute   string       `db:"attribute" json:"attribute"`
	Operator    RuleOperator `db:"operator" json:"operator"`
	Values      string       `db:"values" json:"values"` // JSON array string
	CreatedAt   time.Time    `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time    `db:"updated_at" json:"updatedAt"`
}

type Campaign struct {
	ID          string     `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	Identifier  string     `db:"identifier" json:"identifier"`
	StartsAt    *time.Time `db:"starts_at" json:"startsAt"`
	EndsAt      *time.Time `db:"ends_at" json:"endsAt"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
}

type CampaignBudgetType string

const (
	CampaignBudgetTypeSpend CampaignBudgetType = "spend"
	CampaignBudgetTypeUsage CampaignBudgetType = "usage"
)

type CampaignBudget struct {
	ID         string             `db:"id" json:"id"`
	CampaignID string             `db:"campaign_id" json:"campaignId"`
	Type       CampaignBudgetType `db:"type" json:"type"`
	Limit      *int64             `db:"budget_limit" json:"limit"`
	Used       int64              `db:"used" json:"used"`
	CreatedAt  time.Time          `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time          `db:"updated_at" json:"updatedAt"`
}
