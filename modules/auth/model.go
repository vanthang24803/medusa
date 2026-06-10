package auth

import "time"

// AuthIdentity — 1 person = 1 identity. AppMetadata stores { "customer_id": ... }
// or { "user_id": ... }. No FK to customer/user — lookup at app layer.
type AuthIdentity struct {
	ID          string    `db:"id" json:"id"`
	AppMetadata []byte    `db:"app_metadata" json:"-"` // jsonb raw
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// ProviderIdentity — 1 identity has multiple providers (email + google + github).
// ProviderMetadata stores hashed password / OAuth tokens (opaque).
type ProviderIdentity struct {
	ID               string    `db:"id" json:"id"`
	AuthIdentityID   string    `db:"auth_identity_id" json:"auth_identity_id"`
	Provider         string    `db:"provider" json:"provider"` // emailpass | google | github
	EntityID         string    `db:"entity_id" json:"entity_id"`
	ProviderMetadata []byte    `db:"provider_metadata" json:"-"`
	UserMetadata     []byte    `db:"user_metadata" json:"user_metadata"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

type APIKeyType string

const (
	APIKeyTypePublishable APIKeyType = "publishable"
	APIKeyTypeSecret      APIKeyType = "secret"
)

// APIKey — publishable cho storefront, secret cho server-to-server.
type APIKey struct {
	ID         string     `db:"id" json:"id"`
	Token      string     `db:"token" json:"token"`
	Type       APIKeyType `db:"type" json:"type"`
	Title      string     `db:"title" json:"title"`
	LastUsedAt *time.Time `db:"last_used_at" json:"last_used_at"`
	RevokedAt  *time.Time `db:"revoked_at" json:"revoked_at"`
	CreatedBy  *string    `db:"created_by" json:"created_by"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updated_at"`
}
