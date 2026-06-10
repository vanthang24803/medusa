-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ordering.order_address (
    id VARCHAR(255) PRIMARY KEY, first_name VARCHAR(255), last_name VARCHAR(255),
    company VARCHAR(255), address_1 VARCHAR(255), address_2 VARCHAR(255), city VARCHAR(255),
    province VARCHAR(255), postal_code VARCHAR(50), country_code VARCHAR(2), phone VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS ordering."order" (
    id VARCHAR(255) PRIMARY KEY,
    display_id BIGSERIAL UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    customer_id VARCHAR(255), customer_email VARCHAR(255) NOT NULL,
    currency_code VARCHAR(10) NOT NULL, region_id VARCHAR(255), sales_channel_id VARCHAR(255),
    shipping_address_id VARCHAR(255) REFERENCES ordering.order_address(id),
    billing_address_id VARCHAR(255) REFERENCES ordering.order_address(id),
    no_notification BOOLEAN, metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_order_customer ON ordering."order"(customer_id);
CREATE TABLE IF NOT EXISTS ordering.order_line_item (
    id VARCHAR(255) PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL REFERENCES ordering."order"(id) ON DELETE CASCADE,
    variant_id VARCHAR(255), product_id VARCHAR(255), title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255), variant_title VARCHAR(255), variant_sku VARCHAR(255),
    thumbnail VARCHAR(1000), unit_price BIGINT NOT NULL, quantity INTEGER NOT NULL,
    fulfilled_quantity INTEGER NOT NULL DEFAULT 0, returned_quantity INTEGER NOT NULL DEFAULT 0,
    cancelled_quantity INTEGER NOT NULL DEFAULT 0,
    is_discountable BOOLEAN NOT NULL DEFAULT true, is_tax_inclusive BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_order_line_order ON ordering.order_line_item(order_id);
CREATE TABLE IF NOT EXISTS ordering.order_shipping_method (
    id VARCHAR(255) PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL REFERENCES ordering."order"(id) ON DELETE CASCADE,
    shipping_option_id VARCHAR(255), name VARCHAR(255) NOT NULL, amount BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS ordering.return (
    id VARCHAR(255) PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL REFERENCES ordering."order"(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL DEFAULT 'requested', refund_amount BIGINT, note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS ordering.return_item (
    id VARCHAR(255) PRIMARY KEY,
    return_id VARCHAR(255) NOT NULL REFERENCES ordering.return(id) ON DELETE CASCADE,
    line_item_id VARCHAR(255) NOT NULL, quantity INTEGER NOT NULL, note VARCHAR(255), reason_id VARCHAR(255)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ordering.return_item;
DROP TABLE IF EXISTS ordering.return;
DROP TABLE IF EXISTS ordering.order_shipping_method;
DROP TABLE IF EXISTS ordering.order_line_item;
DROP TABLE IF EXISTS ordering."order";
DROP TABLE IF EXISTS ordering.order_address;
-- +goose StatementEnd
