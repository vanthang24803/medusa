-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS brand.brand (
    id          VARCHAR(255) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL UNIQUE,
    logo_url    VARCHAR(1000),
    description TEXT,
    is_active   BOOLEAN      NOT NULL DEFAULT true,
    rank        INTEGER      NOT NULL DEFAULT 0,
    metadata    JSONB        NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at  TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_brand_slug ON brand.brand(slug);
CREATE INDEX IF NOT EXISTS idx_brand_active ON brand.brand(is_active);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS brand.brand;
-- +goose StatementEnd
