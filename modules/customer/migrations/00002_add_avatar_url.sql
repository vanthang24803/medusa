-- +goose Up
ALTER TABLE customer.customer ADD COLUMN IF NOT EXISTS avatar_url TEXT;

-- +goose Down
ALTER TABLE customer.customer DROP COLUMN IF EXISTS avatar_url;
