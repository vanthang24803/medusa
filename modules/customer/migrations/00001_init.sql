-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS customer.customer (
    id           VARCHAR(255) PRIMARY KEY,
    email        VARCHAR(255) NOT NULL UNIQUE,
    first_name   VARCHAR(255),
    last_name    VARCHAR(255),
    phone        VARCHAR(50),
    company_name VARCHAR(255),
    has_account  BOOLEAN      NOT NULL DEFAULT false,
    metadata     JSONB        NOT NULL DEFAULT '{}',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS customer.customer_address (
    id                  VARCHAR(255) PRIMARY KEY,
    customer_id         VARCHAR(255) NOT NULL REFERENCES customer.customer(id) ON DELETE CASCADE,
    address_name        VARCHAR(255),
    first_name          VARCHAR(255),
    last_name           VARCHAR(255),
    company             VARCHAR(255),
    address_1           VARCHAR(255),
    address_2           VARCHAR(255),
    city                VARCHAR(255),
    province            VARCHAR(255),
    postal_code         VARCHAR(50),
    country_code        VARCHAR(2),
    phone               VARCHAR(50),
    is_default_shipping BOOLEAN      NOT NULL DEFAULT false,
    is_default_billing  BOOLEAN      NOT NULL DEFAULT false,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_customer_address_customer ON customer.customer_address(customer_id);

CREATE TABLE IF NOT EXISTS customer.customer_group (
    id         VARCHAR(255) PRIMARY KEY,
    name       VARCHAR(255) NOT NULL UNIQUE,
    metadata   JSONB        NOT NULL DEFAULT '{}',
    created_by VARCHAR(255),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS customer.customer_group;
DROP TABLE IF EXISTS customer.customer_address;
DROP TABLE IF EXISTS customer.customer;
-- +goose StatementEnd
