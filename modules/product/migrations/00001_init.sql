-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS product.product_collection (
    id         VARCHAR(255) PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    handle     VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS product.product_type (
    id VARCHAR(255) PRIMARY KEY, value VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product.product (
    id            VARCHAR(255) PRIMARY KEY,
    title         VARCHAR(255) NOT NULL,
    subtitle      VARCHAR(255),
    description   TEXT,
    handle        VARCHAR(255) NOT NULL UNIQUE,
    is_giftcard   BOOLEAN      NOT NULL DEFAULT false,
    status        VARCHAR(50)  NOT NULL DEFAULT 'draft',
    thumbnail     VARCHAR(1000),
    weight        DOUBLE PRECISION,
    length        DOUBLE PRECISION,
    height        DOUBLE PRECISION,
    width         DOUBLE PRECISION,
    origin_country VARCHAR(2),
    hs_code       VARCHAR(255),
    mid_code      VARCHAR(255),
    material      VARCHAR(255),
    discountable  BOOLEAN      NOT NULL DEFAULT true,
    external_id   VARCHAR(255),
    collection_id VARCHAR(255) REFERENCES product.product_collection(id),
    type_id       VARCHAR(255) REFERENCES product.product_type(id),
    metadata      JSONB        NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at    TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_product_collection ON product.product(collection_id);
CREATE INDEX IF NOT EXISTS idx_product_status ON product.product(status);

CREATE TABLE IF NOT EXISTS product.product_variant (
    id               VARCHAR(255) PRIMARY KEY,
    product_id       VARCHAR(255) NOT NULL REFERENCES product.product(id) ON DELETE CASCADE,
    title            VARCHAR(255) NOT NULL,
    sku              VARCHAR(255) UNIQUE,
    barcode          VARCHAR(255),
    ean              VARCHAR(255),
    upc              VARCHAR(255),
    allow_backorder  BOOLEAN      NOT NULL DEFAULT false,
    manage_inventory BOOLEAN      NOT NULL DEFAULT true,
    weight           DOUBLE PRECISION,
    rank             INTEGER      NOT NULL DEFAULT 0,
    metadata         JSONB        NOT NULL DEFAULT '{}',
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at       TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_variant_product ON product.product_variant(product_id);

CREATE TABLE IF NOT EXISTS product.product_option (
    id VARCHAR(255) PRIMARY KEY,
    product_id VARCHAR(255) NOT NULL REFERENCES product.product(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL, rank INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product.product_option_value (
    id VARCHAR(255) PRIMARY KEY,
    option_id VARCHAR(255) NOT NULL REFERENCES product.product_option(id) ON DELETE CASCADE,
    value VARCHAR(255) NOT NULL, rank INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product.product_category (
    id                 VARCHAR(255) PRIMARY KEY,
    name               VARCHAR(255) NOT NULL,
    description        TEXT,
    handle             VARCHAR(255) NOT NULL UNIQUE,
    is_active          BOOLEAN      NOT NULL DEFAULT false,
    is_internal        BOOLEAN      NOT NULL DEFAULT false,
    rank               INTEGER      NOT NULL DEFAULT 0,
    parent_category_id VARCHAR(255) REFERENCES product.product_category(id),
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product.product_image (
    id VARCHAR(255) PRIMARY KEY,
    product_id VARCHAR(255) NOT NULL REFERENCES product.product(id) ON DELETE CASCADE,
    url VARCHAR(1000) NOT NULL, rank INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS product.product_tag (
    id VARCHAR(255) PRIMARY KEY, value VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product.product_tag;
DROP TABLE IF EXISTS product.product_image;
DROP TABLE IF EXISTS product.product_category;
DROP TABLE IF EXISTS product.product_option_value;
DROP TABLE IF EXISTS product.product_option;
DROP TABLE IF EXISTS product.product_variant;
DROP TABLE IF EXISTS product.product;
DROP TABLE IF EXISTS product.product_type;
DROP TABLE IF EXISTS product.product_collection;
-- +goose StatementEnd
