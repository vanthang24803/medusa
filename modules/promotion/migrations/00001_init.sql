-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS promotion.campaign (
    id VARCHAR(255) PRIMARY KEY, name VARCHAR(255) NOT NULL, description TEXT,
    identifier VARCHAR(255) NOT NULL UNIQUE, starts_at TIMESTAMPTZ, ends_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS promotion.campaign_budget (
    id VARCHAR(255) PRIMARY KEY,
    campaign_id VARCHAR(255) NOT NULL UNIQUE REFERENCES promotion.campaign(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL, budget_limit BIGINT, used BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS promotion.promotion (
    id VARCHAR(255) PRIMARY KEY, code VARCHAR(255) NOT NULL UNIQUE, type VARCHAR(50) NOT NULL,
    is_automatic BOOLEAN NOT NULL DEFAULT false, status VARCHAR(50) NOT NULL DEFAULT 'inactive',
    campaign_id VARCHAR(255) REFERENCES promotion.campaign(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS promotion.promotion_rule (
    id VARCHAR(255) PRIMARY KEY,
    promotion_id VARCHAR(255) NOT NULL REFERENCES promotion.promotion(id) ON DELETE CASCADE,
    attribute VARCHAR(255) NOT NULL, operator VARCHAR(20) NOT NULL, values TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS promotion.promotion_rule;
DROP TABLE IF EXISTS promotion.promotion;
DROP TABLE IF EXISTS promotion.campaign_budget;
DROP TABLE IF EXISTS promotion.campaign;
-- +goose StatementEnd
