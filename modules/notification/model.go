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
	TemplateID   string              `db:"template_id" json:"templateId"`
	ProviderID   string              `db:"provider_id" json:"providerId"`
	Status       NotificationStatus  `db:"status" json:"status"`
	Data         []byte              `db:"data" json:"-"`
	ExternalID   *string             `db:"external_id" json:"externalId"`
	ResourceID   *string             `db:"resource_id" json:"resourceId"`
	ResourceType *string             `db:"resource_type" json:"resourceType"`
	SentAt       *time.Time          `db:"sent_at" json:"sentAt"`
	CreatedAt    time.Time           `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time           `db:"updated_at" json:"updatedAt"`
}

type NotificationProvider struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Handle    string    `db:"handle" json:"handle"`
	IsEnabled bool      `db:"is_enabled" json:"isEnabled"`
	Config    []byte    `db:"config" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
