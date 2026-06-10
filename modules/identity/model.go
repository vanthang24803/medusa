package identity

import "time"

// User — admin dashboard users, profile only. Credentials in auth schema.
type User struct {
	ID        string     `db:"id" json:"id"`
	Email     string     `db:"email" json:"email"`
	FirstName *string    `db:"first_name" json:"first_name"`
	LastName  *string    `db:"last_name" json:"last_name"`
	AvatarURL *string    `db:"avatar_url" json:"avatar_url"`
	Metadata  []byte     `db:"metadata" json:"-"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at"`
}

// Invite — token to invite a new admin.
type Invite struct {
	ID        string    `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Token     string    `db:"token" json:"token"`
	Accepted  bool      `db:"accepted" json:"accepted"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	CreatedBy *string   `db:"created_by" json:"created_by"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
