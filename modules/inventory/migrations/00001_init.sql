-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS inventory.inventory_item (
    id VARCHAR(255) PRIMARY KEY, sku VARCHAR(255) UNIQUE, title VARCHAR(255),
    description TEXT, thumbnail VARCHAR(1000), requires_shipping BOOLEAN NOT NULL DEFAULT true,
    weight DOUBLE PRECISION, origin_country VARCHAR(2), hs_code VARCHAR(255), material VARCHAR(255),
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS inventory.inventory_level (
    id VARCHAR(255) PRIMARY KEY,
    inventory_item_id VARCHAR(255) NOT NULL REFERENCES inventory.inventory_item(id) ON DELETE CASCADE,
    location_id VARCHAR(255) NOT NULL,
    stocked_quantity INTEGER NOT NULL DEFAULT 0,
    reserved_quantity INTEGER NOT NULL DEFAULT 0,
    incoming_quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(inventory_item_id, location_id)
);
CREATE TABLE IF NOT EXISTS inventory.reservation_item (
    id VARCHAR(255) PRIMARY KEY,
    inventory_item_id VARCHAR(255) NOT NULL REFERENCES inventory.inventory_item(id) ON DELETE CASCADE,
    location_id VARCHAR(255) NOT NULL, line_item_id VARCHAR(255), quantity INTEGER NOT NULL,
    allow_backorder BOOLEAN NOT NULL DEFAULT false, description VARCHAR(255), created_by VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_reservation_line_item ON inventory.reservation_item(line_item_id);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS inventory.reservation_item;
DROP TABLE IF EXISTS inventory.inventory_level;
DROP TABLE IF EXISTS inventory.inventory_item;
-- +goose StatementEnd
