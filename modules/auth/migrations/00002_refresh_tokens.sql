-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS auth.auth_token (
    id                VARCHAR(255) PRIMARY KEY,
    auth_identity_id  VARCHAR(255) NOT NULL REFERENCES auth.auth_identity(id) ON DELETE CASCADE,
    token_hash        VARCHAR(255) NOT NULL UNIQUE,
    type              VARCHAR(50)  NOT NULL DEFAULT 'refresh',
    expires_at        TIMESTAMPTZ  NOT NULL,
    revoked_at        TIMESTAMPTZ,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_auth_token_identity ON auth.auth_token(auth_identity_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth.auth_token;
-- +goose StatementEnd
