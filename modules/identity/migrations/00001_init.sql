-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS identity."user" (
    id         VARCHAR(255) PRIMARY KEY,
    email      VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(255),
    last_name  VARCHAR(255),
    avatar_url VARCHAR(1000),
    metadata   JSONB        NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS identity.invite (
    id         VARCHAR(255) PRIMARY KEY,
    email      VARCHAR(255) NOT NULL,
    token      VARCHAR(255) NOT NULL UNIQUE,
    accepted   BOOLEAN      NOT NULL DEFAULT false,
    expires_at TIMESTAMPTZ  NOT NULL,
    created_by VARCHAR(255),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_invite_email ON identity.invite(email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS identity.invite;
DROP TABLE IF EXISTS identity."user";
-- +goose StatementEnd
