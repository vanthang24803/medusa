package iam

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ── Action constants ─────────────────────────────────────────────────────────

const (
	ActionWildcard = "*"

	ActionProductRead   = "product:Read"
	ActionProductCreate = "product:Create"
	ActionProductUpdate = "product:Update"
	ActionProductDelete = "product:Delete"
	ActionProductManage = "product:*"

	ActionBrandRead   = "brand:Read"
	ActionBrandCreate = "brand:Create"
	ActionBrandUpdate = "brand:Update"
	ActionBrandDelete = "brand:Delete"
	ActionBrandManage = "brand:*"

	ActionOrderRead   = "order:Read"
	ActionOrderCreate = "order:Create"
	ActionOrderUpdate = "order:Update"
	ActionOrderDelete = "order:Delete"
	ActionOrderManage = "order:*"

	ActionCustomerRead   = "customer:Read"
	ActionCustomerCreate = "customer:Create"
	ActionCustomerUpdate = "customer:Update"
	ActionCustomerDelete = "customer:Delete"
	ActionCustomerManage = "customer:*"

	ActionInventoryRead   = "inventory:Read"
	ActionInventoryCreate = "inventory:Create"
	ActionInventoryUpdate = "inventory:Update"
	ActionInventoryDelete = "inventory:Delete"
	ActionInventoryManage = "inventory:*"

	ActionPricingRead   = "pricing:Read"
	ActionPricingCreate = "pricing:Create"
	ActionPricingUpdate = "pricing:Update"
	ActionPricingDelete = "pricing:Delete"
	ActionPricingManage = "pricing:*"

	ActionCartRead   = "cart:Read"
	ActionCartManage = "cart:*"

	ActionFulfillmentRead   = "fulfillment:Read"
	ActionFulfillmentCreate = "fulfillment:Create"
	ActionFulfillmentUpdate = "fulfillment:Update"
	ActionFulfillmentManage = "fulfillment:*"

	ActionPromotionRead   = "promotion:Read"
	ActionPromotionCreate = "promotion:Create"
	ActionPromotionUpdate = "promotion:Update"
	ActionPromotionDelete = "promotion:Delete"
	ActionPromotionManage = "promotion:*"

	ActionRegionRead   = "region:Read"
	ActionRegionCreate = "region:Create"
	ActionRegionUpdate = "region:Update"
	ActionRegionDelete = "region:Delete"
	ActionRegionManage = "region:*"

	ActionNotificationRead   = "notification:Read"
	ActionNotificationCreate = "notification:Create"
	ActionNotificationManage = "notification:*"

	ActionIAMManage      = "iam:Manage"
	ActionIdentityManage = "identity:Manage" // super_admin only (wildcard covers it)
	ActionUserManage     = "user:Manage"     // list users, ban/unban — assignable to any role
)

// ── DB models ────────────────────────────────────────────────────────────────

type Role struct {
	ID            string    `db:"id" json:"id"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	IsSystem      bool      `db:"is_system" json:"isSystem"`
	LastUpdatedBy *string   `db:"last_updated_by" json:"lastUpdatedBy"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time `db:"updated_at" json:"updatedAt"`
}

type Policy struct {
	ID            string         `db:"id" json:"id"`
	Name          string         `db:"name" json:"name"`
	Description   string         `db:"description" json:"description"`
	Document      PolicyDocument `db:"document" json:"document"`
	IsSystem      bool           `db:"is_system" json:"isSystem"`
	LastUpdatedBy *string        `db:"last_updated_by" json:"lastUpdatedBy"`
	CreatedAt     time.Time      `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time      `db:"updated_at" json:"updatedAt"`
}

type PrincipalRole struct {
	AuthIdentityID string    `db:"auth_identity_id" json:"authIdentityId"`
	RoleID         string    `db:"role_id" json:"roleId"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
}

type APIKeyPolicy struct {
	APIKeyID  string    `db:"api_key_id" json:"apiKeyId"`
	PolicyID  string    `db:"policy_id" json:"policyId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

// ── Policy document ──────────────────────────────────────────────────────────

type Effect string

const (
	EffectAllow Effect = "Allow"
	EffectDeny  Effect = "Deny"
)

type Statement struct {
	SID       string   `json:"sid,omitempty"`
	Effect    Effect   `json:"effect"`
	Actions   []string `json:"actions"`
	Resources []string `json:"resources"`
}

// PolicyDocument is the JSON structure stored in iam.policy.document (JSONB).
type PolicyDocument struct {
	Version    string      `json:"version"`
	Statements []Statement `json:"statements"`
}

func (d PolicyDocument) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *PolicyDocument) Scan(src any) error {
	var b []byte
	switch v := src.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return fmt.Errorf("unsupported type: %T", src)
	}
	return json.Unmarshal(b, d)
}

// ── Request/response DTOs ────────────────────────────────────────────────────

type CreateRoleReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreatePolicyReq struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Document    PolicyDocument `json:"document"`
}

type UpdatePolicyReq struct {
	Description *string         `json:"description"`
	Document    *PolicyDocument `json:"document"`
}

// ── Action matching ──────────────────────────────────────────────────────────

// matchesAction returns true if the policy action pattern matches the requested action.
// Supports: "*", "service:*", "service:Action"
func matchesAction(pattern, action string) bool {
	if pattern == ActionWildcard {
		return true
	}
	if strings.HasSuffix(pattern, ":*") {
		prefix := strings.TrimSuffix(pattern, "*")
		return strings.HasPrefix(action, prefix)
	}
	return pattern == action
}

// matchesResource returns true if the policy resource pattern matches.
// Current service-level granularity: "service/*" or "*".
func matchesResource(pattern, resource string) bool {
	if pattern == ActionWildcard {
		return true
	}
	if before, ok := strings.CutSuffix(pattern, "/*"); ok {
		prefix := before
		parts := strings.SplitN(resource, "/", 2)
		return parts[0] == prefix
	}
	return pattern == resource
}

// Evaluate runs AWS-style policy evaluation across a slice of documents.
// Returns true only if there is at least one Allow and no Deny for (action, resource).
func Evaluate(docs []PolicyDocument, action, resource string) bool {
	allowed := false
	for _, doc := range docs {
		for _, stmt := range doc.Statements {
			actionMatch := false
			for _, a := range stmt.Actions {
				if matchesAction(a, action) {
					actionMatch = true
					break
				}
			}
			if !actionMatch {
				continue
			}
			resourceMatch := false
			for _, res := range stmt.Resources {
				if matchesResource(res, resource) {
					resourceMatch = true
					break
				}
			}
			if !resourceMatch {
				continue
			}
			if stmt.Effect == EffectDeny {
				return false // explicit deny wins immediately
			}
			if stmt.Effect == EffectAllow {
				allowed = true
			}
		}
	}
	return allowed
}
