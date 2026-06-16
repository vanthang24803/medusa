-- +goose Up
-- +goose StatementBegin
ALTER TABLE product.product
    ADD COLUMN IF NOT EXISTS brand_id         VARCHAR(255),
    ADD COLUMN IF NOT EXISTS author           VARCHAR(500),
    ADD COLUMN IF NOT EXISTS isbn             VARCHAR(50),
    ADD COLUMN IF NOT EXISTS page_count       INTEGER,
    ADD COLUMN IF NOT EXISTS compare_at_price BIGINT,
    ADD COLUMN IF NOT EXISTS quantity         INTEGER,
    ADD COLUMN IF NOT EXISTS rating           DOUBLE PRECISION,
    ADD COLUMN IF NOT EXISTS review_count     INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS is_featured      BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN IF NOT EXISTS published_at     TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_product_brand ON product.product(brand_id);
CREATE INDEX IF NOT EXISTS idx_product_featured ON product.product(is_featured) WHERE is_featured = true;
CREATE INDEX IF NOT EXISTS idx_product_published ON product.product(published_at DESC NULLS LAST);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE product.product
    DROP COLUMN IF EXISTS brand_id,
    DROP COLUMN IF EXISTS author,
    DROP COLUMN IF EXISTS isbn,
    DROP COLUMN IF EXISTS page_count,
    DROP COLUMN IF EXISTS compare_at_price,
    DROP COLUMN IF EXISTS quantity,
    DROP COLUMN IF EXISTS rating,
    DROP COLUMN IF EXISTS review_count,
    DROP COLUMN IF EXISTS is_featured,
    DROP COLUMN IF EXISTS published_at;
-- +goose StatementEnd
