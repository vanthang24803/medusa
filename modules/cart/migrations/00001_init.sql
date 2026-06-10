-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cart.cart_address (
    id VARCHAR(255) PRIMARY KEY, first_name VARCHAR(255), last_name VARCHAR(255),
    company VARCHAR(255), address_1 VARCHAR(255), address_2 VARCHAR(255), city VARCHAR(255),
    province VARCHAR(255), postal_code VARCHAR(50), country_code VARCHAR(2), phone VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS cart.cart (
    id VARCHAR(255) PRIMARY KEY, email VARCHAR(255), customer_id VARCHAR(255),
    currency_code VARCHAR(10) NOT NULL, region_id VARCHAR(255), sales_channel_id VARCHAR(255),
    shipping_address_id VARCHAR(255) REFERENCES cart.cart_address(id),
    billing_address_id VARCHAR(255) REFERENCES cart.cart_address(id),
    completed_at TIMESTAMPTZ, metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_cart_customer ON cart.cart(customer_id);
CREATE TABLE IF NOT EXISTS cart.cart_line_item (
    id VARCHAR(255) PRIMARY KEY,
    cart_id VARCHAR(255) NOT NULL REFERENCES cart.cart(id) ON DELETE CASCADE,
    variant_id VARCHAR(255), product_id VARCHAR(255), product_title VARCHAR(255),
    product_description TEXT, variant_title VARCHAR(255), variant_sku VARCHAR(255),
    thumbnail VARCHAR(1000), quantity INTEGER NOT NULL, unit_price BIGINT NOT NULL,
    is_discountable BOOLEAN NOT NULL DEFAULT true, is_tax_inclusive BOOLEAN NOT NULL DEFAULT false,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_cart_line_cart ON cart.cart_line_item(cart_id);
CREATE TABLE IF NOT EXISTS cart.cart_shipping_method (
    id VARCHAR(255) PRIMARY KEY,
    cart_id VARCHAR(255) NOT NULL REFERENCES cart.cart(id) ON DELETE CASCADE,
    shipping_option_id VARCHAR(255), name VARCHAR(255) NOT NULL, amount BIGINT NOT NULL,
    is_tax_inclusive BOOLEAN NOT NULL DEFAULT false, data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS cart.cart_adjustment (
    id VARCHAR(255) PRIMARY KEY,
    cart_id VARCHAR(255) NOT NULL REFERENCES cart.cart(id) ON DELETE CASCADE,
    item_id VARCHAR(255), promotion_id VARCHAR(255), code VARCHAR(255),
    amount BIGINT NOT NULL, description VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cart.cart_adjustment;
DROP TABLE IF EXISTS cart.cart_shipping_method;
DROP TABLE IF EXISTS cart.cart_line_item;
DROP TABLE IF EXISTS cart.cart;
DROP TABLE IF EXISTS cart.cart_address;
-- +goose StatementEnd
