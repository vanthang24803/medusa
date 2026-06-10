-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pricing.price_set (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS pricing.price_list (
    id VARCHAR(255) PRIMARY KEY, title VARCHAR(255) NOT NULL, description VARCHAR(1000),
    type VARCHAR(50) NOT NULL, status VARCHAR(50) NOT NULL DEFAULT 'draft',
    starts_at TIMESTAMPTZ, ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS pricing.price (
    id VARCHAR(255) PRIMARY KEY,
    price_set_id VARCHAR(255) NOT NULL REFERENCES pricing.price_set(id) ON DELETE CASCADE,
    currency_code VARCHAR(10) NOT NULL, amount BIGINT NOT NULL,
    min_quantity INTEGER, max_quantity INTEGER,
    price_list_id VARCHAR(255) REFERENCES pricing.price_list(id),
    rules_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_price_set ON pricing.price(price_set_id);
CREATE TABLE IF NOT EXISTS pricing.price_rule (
    id VARCHAR(255) PRIMARY KEY,
    price_id VARCHAR(255) NOT NULL REFERENCES pricing.price(id) ON DELETE CASCADE,
    attribute VARCHAR(255) NOT NULL, operator VARCHAR(20) NOT NULL, value VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pricing.price_rule;
DROP TABLE IF EXISTS pricing.price;
DROP TABLE IF EXISTS pricing.price_list;
DROP TABLE IF EXISTS pricing.price_set;
-- +goose StatementEnd
