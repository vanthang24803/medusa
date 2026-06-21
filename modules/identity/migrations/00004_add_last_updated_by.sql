-- +goose Up
ALTER TABLE identity."user" ADD COLUMN IF NOT EXISTS last_updated_by VARCHAR(255);

-- +goose Down
ALTER TABLE identity."user" DROP COLUMN IF EXISTS last_updated_by;
