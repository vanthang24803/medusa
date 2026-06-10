-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS region.sales_channel (
    id VARCHAR(255) PRIMARY KEY, name VARCHAR(255) NOT NULL, description VARCHAR(1000),
    is_disabled BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS region.region (
    id VARCHAR(255) PRIMARY KEY, name VARCHAR(255) NOT NULL, currency_code VARCHAR(10) NOT NULL,
    automatic_taxes BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE TABLE IF NOT EXISTS region.store (
    id VARCHAR(255) PRIMARY KEY, name VARCHAR(255) NOT NULL,
    supported_currencies JSONB NOT NULL DEFAULT '[]',
    default_sales_channel_id VARCHAR(255), default_region_id VARCHAR(255), default_location_id VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS region.region_country (
    id VARCHAR(255) PRIMARY KEY,
    region_id VARCHAR(255) NOT NULL REFERENCES region.region(id) ON DELETE CASCADE,
    iso_2 VARCHAR(2) NOT NULL UNIQUE, iso_3 VARCHAR(3) NOT NULL, name VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL, num_code VARCHAR(10)
);
CREATE TABLE IF NOT EXISTS region.tax_rate (
    id VARCHAR(255) PRIMARY KEY,
    region_id VARCHAR(255) NOT NULL REFERENCES region.region(id) ON DELETE CASCADE,
    country_code VARCHAR(2), province_code VARCHAR(10), name VARCHAR(255) NOT NULL, code VARCHAR(50),
    rate DOUBLE PRECISION NOT NULL, is_default BOOLEAN NOT NULL DEFAULT false,
    is_combinable BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS region.tax_rate;
DROP TABLE IF EXISTS region.region_country;
DROP TABLE IF EXISTS region.store;
DROP TABLE IF EXISTS region.region;
DROP TABLE IF EXISTS region.sales_channel;
-- +goose StatementEnd
