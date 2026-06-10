package notification

import "time"

type NotificationChannel string
type NotificationStatus string

const (
	ChannelEmail NotificationChannel = "email"
	ChannelSMS   NotificationChannel = "sms"
	ChannelPush  NotificationChannel = "push"

	StatusPending NotificationStatus = "pending"
	StatusSent    NotificationStatus = "sent"
	StatusFailed  NotificationStatus = "failed"
)

// Notification — logs all sent notifications.
type Notification struct {
	ID           string              `db:"id" json:"id"`
	To           string              `db:"to_address" json:"to"`
	Channel      NotificationChannel `db:"channel" json:"channel"`
	TemplateID   string              `db:"template_id" json:"template_id"`
	ProviderID   string              `db:"provider_id" json:"provider_id"`
	Status       NotificationStatus  `db:"status" json:"status"`
	Data         []byte              `db:"data" json:"-"`
	ExternalID   *string             `db:"external_id" json:"external_id"`
	ResourceID   *string             `db:"resource_id" json:"resource_id"`
	ResourceType *string             `db:"resource_type" json:"resource_type"`
	SentAt       *time.Time          `db:"sent_at" json:"sent_at"`
	CreatedAt    time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time           `db:"updated_at" json:"updated_at"`
}

type NotificationProvider struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Handle    string    `db:"handle" json:"handle"`
	IsEnabled bool      `db:"is_enabled" json:"is_enabled"`
	Config    []byte    `db:"config" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
