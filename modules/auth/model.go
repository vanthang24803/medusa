package auth

import "time"

type AuthIdentity struct {
	ID          string    `db:"id" json:"id"`
	AppMetadata []byte    `db:"app_metadata" json:"-"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

type ProviderIdentity struct {
	ID               string    `db:"id" json:"id"`
	AuthIdentityID   string    `db:"auth_identity_id" json:"authIdentityId"`
	Provider         string    `db:"provider" json:"provider"`
	EntityID         string    `db:"entity_id" json:"entityId"`
	ProviderMetadata []byte    `db:"provider_metadata" json:"-"`
	UserMetadata     []byte    `db:"user_metadata" json:"userMetadata"`
	CreatedAt        time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time `db:"updated_at" json:"updatedAt"`
}

type APIKeyType string

const (
	APIKeyTypePublishable APIKeyType = "publishable"
	APIKeyTypeSecret      APIKeyType = "secret"
)

type APIKey struct {
	ID         string     `db:"id" json:"id"`
	Token      string     `db:"token" json:"token"`
	Type       APIKeyType `db:"type" json:"type"`
	Title      string     `db:"title" json:"title"`
	LastUsedAt *time.Time `db:"last_used_at" json:"lastUsedAt"`
	RevokedAt  *time.Time `db:"revoked_at" json:"revokedAt"`
	CreatedBy  *string    `db:"created_by" json:"createdBy"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}

type AuthToken struct {
	ID             string     `db:"id" json:"id"`
	AuthIdentityID string     `db:"auth_identity_id" json:"authIdentityId"`
	TokenHash      string     `db:"token_hash" json:"-"`
	Type           string     `db:"type" json:"type"`
	ExpiresAt      time.Time  `db:"expires_at" json:"expiresAt"`
	RevokedAt      *time.Time `db:"revoked_at" json:"revokedAt"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
}

type RegisterReq struct {
	Email     string `json:"email"     validate:"required,email"`
	Password  string `json:"password"  validate:"required,min=6"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"  validate:"required"`
}

type LoginReq struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshReq struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type CustomerInfo struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Phone     *string `json:"phone"`
	CreatedAt string  `json:"createdAt"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type RefreshResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
