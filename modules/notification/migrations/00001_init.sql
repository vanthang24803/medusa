-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification.notification_provider (
    id VARCHAR(255) PRIMARY KEY, name VARCHAR(255) NOT NULL, handle VARCHAR(100) NOT NULL UNIQUE,
    is_enabled BOOLEAN NOT NULL DEFAULT true, config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS notification.notification (
    id VARCHAR(255) PRIMARY KEY, to_address VARCHAR(255) NOT NULL, channel VARCHAR(50) NOT NULL,
    template_id VARCHAR(255) NOT NULL, provider_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending', data JSONB NOT NULL DEFAULT '{}',
    external_id VARCHAR(255), resource_id VARCHAR(255), resource_type VARCHAR(100),
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_notif_resource ON notification.notification(resource_id);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification.notification;
DROP TABLE IF EXISTS notification.notification_provider;
-- +goose StatementEnd
