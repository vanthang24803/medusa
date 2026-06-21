-- +goose Up
ALTER TABLE brand.brand ADD COLUMN IF NOT EXISTS last_updated_by VARCHAR(255);

-- +goose Down
ALTER TABLE brand.brand DROP COLUMN IF EXISTS last_updated_by;
