-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS auth.auth_identity (
    id            VARCHAR(255) PRIMARY KEY,
    app_metadata  JSONB        NOT NULL DEFAULT '{}',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS auth.auth_provider_identity (
    id                VARCHAR(255) PRIMARY KEY,
    auth_identity_id  VARCHAR(255) NOT NULL REFERENCES auth.auth_identity(id) ON DELETE CASCADE,
    provider          VARCHAR(100) NOT NULL,
    entity_id         VARCHAR(255) NOT NULL UNIQUE,
    provider_metadata JSONB        NOT NULL DEFAULT '{}',
    user_metadata     JSONB        NOT NULL DEFAULT '{}',
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_provider_identity_auth ON auth.auth_provider_identity(auth_identity_id);

CREATE TABLE IF NOT EXISTS auth.api_key (
    id           VARCHAR(255) PRIMARY KEY,
    token        VARCHAR(255) NOT NULL UNIQUE,
    type         VARCHAR(50)  NOT NULL,
    title        VARCHAR(255) NOT NULL,
    last_used_at TIMESTAMPTZ,
    revoked_at   TIMESTAMPTZ,
    created_by   VARCHAR(255),
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS auth.api_key;
DROP TABLE IF EXISTS auth.auth_provider_identity;
DROP TABLE IF EXISTS auth.auth_identity;
-- +goose StatementEnd
