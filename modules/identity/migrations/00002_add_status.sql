-- +goose Up
-- +goose StatementBegin
ALTER TABLE identity."user"
    ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'active';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE identity."user" DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
