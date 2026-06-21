-- +goose Up
ALTER TABLE customer.customer         ADD COLUMN IF NOT EXISTS last_updated_by VARCHAR(255);
ALTER TABLE customer.customer_address ADD COLUMN IF NOT EXISTS last_updated_by VARCHAR(255);

-- +goose Down
ALTER TABLE customer.customer_address DROP COLUMN IF EXISTS last_updated_by;
ALTER TABLE customer.customer         DROP COLUMN IF EXISTS last_updated_by;
