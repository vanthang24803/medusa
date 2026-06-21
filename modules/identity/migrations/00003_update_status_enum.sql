-- +goose Up
-- Migrate old inactive status to closed
UPDATE identity."user" SET status = 'closed' WHERE status = 'inactive';

ALTER TABLE identity."user"
    ADD CONSTRAINT user_status_check CHECK (status IN ('init', 'active', 'ban', 'closed'));

-- +goose Down
ALTER TABLE identity."user" DROP CONSTRAINT IF EXISTS user_status_check;
UPDATE identity."user" SET status = 'inactive' WHERE status = 'closed';
UPDATE identity."user" SET status = 'active' WHERE status IN ('init', 'ban');
