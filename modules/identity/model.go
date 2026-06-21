package identity

import "time"

const (
	StatusInit   = "init"   // created, awaiting activation
	StatusActive = "active" // can log in
	StatusBan    = "ban"    // banned by user manager
	StatusClosed = "closed" // closed by super_admin
)

// User — admin dashboard users, profile only. Credentials in auth schema.
type User struct {
	ID        string     `db:"id" json:"id"`
	Email     string     `db:"email" json:"email"`
	FirstName *string    `db:"first_name" json:"firstName"`
	LastName  *string    `db:"last_name" json:"lastName"`
	AvatarURL *string    `db:"avatar_url" json:"avatarUrl"`
	Status          string     `db:"status" json:"status"` // init | active | ban | closed
	Metadata        []byte     `db:"metadata" json:"-"`
	LastUpdatedBy   *string    `db:"last_updated_by" json:"lastUpdatedBy"`
	CreatedAt       time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt       *time.Time `db:"deleted_at" json:"deletedAt"`
}

type CreateAdminReq struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	RoleID    *string `json:"roleId"` // optional — assign an existing role immediately
}

type UpdateProfileReq struct {
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
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
