-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payment.payment_collection (
    id VARCHAR(255) PRIMARY KEY, order_id VARCHAR(255), cart_id VARCHAR(255),
    currency_code VARCHAR(10) NOT NULL, amount BIGINT NOT NULL,
    authorized_amount BIGINT, captured_amount BIGINT, refunded_amount BIGINT,
    status VARCHAR(50) NOT NULL DEFAULT 'not_paid',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_paycol_order ON payment.payment_collection(order_id);
CREATE TABLE IF NOT EXISTS payment.payment_session (
    id VARCHAR(255) PRIMARY KEY,
    payment_collection_id VARCHAR(255) NOT NULL REFERENCES payment.payment_collection(id) ON DELETE CASCADE,
    provider_id VARCHAR(255) NOT NULL, status VARCHAR(50) NOT NULL DEFAULT 'pending',
    amount BIGINT NOT NULL, currency_code VARCHAR(10) NOT NULL, data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS payment.payment (
    id VARCHAR(255) PRIMARY KEY,
    payment_collection_id VARCHAR(255) NOT NULL REFERENCES payment.payment_collection(id) ON DELETE CASCADE,
    provider_id VARCHAR(255) NOT NULL, currency_code VARCHAR(10) NOT NULL, amount BIGINT NOT NULL,
    authorized_at TIMESTAMPTZ, captured_at TIMESTAMPTZ, cancelled_at TIMESTAMPTZ, data JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS payment.refund (
    id VARCHAR(255) PRIMARY KEY,
    payment_id VARCHAR(255) NOT NULL REFERENCES payment.payment(id) ON DELETE CASCADE,
    amount BIGINT NOT NULL, note VARCHAR(1000), created_by VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS payment.refund;
DROP TABLE IF EXISTS payment.payment;
DROP TABLE IF EXISTS payment.payment_session;
DROP TABLE IF EXISTS payment.payment_collection;
-- +goose StatementEnd
