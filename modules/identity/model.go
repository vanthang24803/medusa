package identity

import "time"

// User — admin dashboard users, profile only. Credentials in auth schema.
type User struct {
	ID        string     `db:"id" json:"id"`
	Email     string     `db:"email" json:"email"`
	FirstName *string    `db:"first_name" json:"firstName"`
	LastName  *string    `db:"last_name" json:"lastName"`
	AvatarURL *string    `db:"avatar_url" json:"avatarUrl"`
	Metadata  []byte     `db:"metadata" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
}

// Invite — token to invite a new admin.
type Invite struct {
	ID        string    `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Token     string    `db:"token" json:"token"`
	Accepted  bool      `db:"accepted" json:"accepted"`
	ExpiresAt time.Time `db:"expires_at" json:"expiresAt"`
	CreatedBy *string   `db:"created_by" json:"createdBy"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
