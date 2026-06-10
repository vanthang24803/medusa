-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS fulfillment.stock_location (
    id VARCHAR(255) PRIMARY KEY, name VARCHAR(255) NOT NULL, address_id VARCHAR(255),
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS fulfillment.shipping_profile (
    id VARCHAR(255) PRIMARY KEY, name VARCHAR(255) NOT NULL, type VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS fulfillment.shipping_option (
    id VARCHAR(255) PRIMARY KEY, name VARCHAR(255) NOT NULL, provider_id VARCHAR(255) NOT NULL,
    service_zone_id VARCHAR(255) NOT NULL,
    shipping_profile_id VARCHAR(255) REFERENCES fulfillment.shipping_profile(id),
    price_type VARCHAR(50) NOT NULL, data JSONB NOT NULL DEFAULT '{}', metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS fulfillment.fulfillment (
    id VARCHAR(255) PRIMARY KEY, order_id VARCHAR(255), location_id VARCHAR(255),
    provider_id VARCHAR(255) NOT NULL, tracking_number VARCHAR(255), tracking_url VARCHAR(1000),
    shipped_at TIMESTAMPTZ, delivered_at TIMESTAMPTZ, canceled_at TIMESTAMPTZ,
    data JSONB NOT NULL DEFAULT '{}', metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_fulfillment_order ON fulfillment.fulfillment(order_id);
CREATE TABLE IF NOT EXISTS fulfillment.fulfillment_item (
    id VARCHAR(255) PRIMARY KEY,
    fulfillment_id VARCHAR(255) NOT NULL REFERENCES fulfillment.fulfillment(id) ON DELETE CASCADE,
    line_item_id VARCHAR(255), inventory_item_id VARCHAR(255), title VARCHAR(255),
    sku VARCHAR(255), barcode VARCHAR(255), quantity INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fulfillment.fulfillment_item;
DROP TABLE IF EXISTS fulfillment.fulfillment;
DROP TABLE IF EXISTS fulfillment.shipping_option;
DROP TABLE IF EXISTS fulfillment.shipping_profile;
DROP TABLE IF EXISTS fulfillment.stock_location;
-- +goose StatementEnd
