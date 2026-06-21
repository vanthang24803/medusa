-- +goose Up
ALTER TABLE iam.role   ADD COLUMN IF NOT EXISTS last_updated_by VARCHAR(255);
ALTER TABLE iam.policy ADD COLUMN IF NOT EXISTS last_updated_by VARCHAR(255);

-- +goose Down
ALTER TABLE iam.policy DROP COLUMN IF EXISTS last_updated_by;
ALTER TABLE iam.role   DROP COLUMN IF EXISTS last_updated_by;
