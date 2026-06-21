-- +goose Up
-- +goose StatementBegin

CREATE TABLE iam.role (
    id          VARCHAR(255) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    is_system   BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE iam.policy (
    id          VARCHAR(255) PRIMARY KEY,
    name        VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    document    JSONB NOT NULL DEFAULT '{}',
    is_system   BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- M:N role ↔ policy
CREATE TABLE iam.role_policy (
    role_id    VARCHAR(255) NOT NULL REFERENCES iam.role(id) ON DELETE CASCADE,
    policy_id  VARCHAR(255) NOT NULL REFERENCES iam.policy(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (role_id, policy_id)
);

-- auth_identity → role assignments (admin users)
CREATE TABLE iam.principal_role (
    auth_identity_id VARCHAR(255) NOT NULL,
    role_id          VARCHAR(255) NOT NULL REFERENCES iam.role(id) ON DELETE CASCADE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (auth_identity_id, role_id)
);

CREATE INDEX idx_principal_role_auth_identity ON iam.principal_role (auth_identity_id);

-- API key → policy assignments (M2M scoped access)
CREATE TABLE iam.api_key_policy (
    api_key_id VARCHAR(255) NOT NULL,
    policy_id  VARCHAR(255) NOT NULL REFERENCES iam.policy(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (api_key_id, policy_id)
);

-- ── Seed built-in policies ────────────────────────────────────────────────────

INSERT INTO iam.policy (id, name, description, document, is_system) VALUES
(
    'pol_system_super_admin',
    'SuperAdminPolicy',
    'Full access to all resources',
    '{"version":"2025-01-01","statements":[{"sid":"AllowAll","effect":"Allow","actions":["*"],"resources":["*"]}]}',
    true
),
(
    'pol_system_admin',
    'AdminPolicy',
    'Full CRUD on catalog, orders, customers, and operations',
    '{"version":"2025-01-01","statements":[{"sid":"AllowAdmin","effect":"Allow","actions":["product:*","brand:*","order:*","customer:*","inventory:*","pricing:*","cart:*","fulfillment:*","promotion:*","region:*","notification:*","iam:Manage"],"resources":["*"]}]}',
    true
),
(
    'pol_system_operator',
    'OperatorPolicy',
    'Read, create and update orders, fulfillments and inventory',
    '{"version":"2025-01-01","statements":[{"sid":"AllowOperator","effect":"Allow","actions":["order:Read","order:Create","order:Update","fulfillment:Read","fulfillment:Create","fulfillment:Update","inventory:Read","inventory:Create","inventory:Update"],"resources":["*"]}]}',
    true
),
(
    'pol_system_viewer',
    'ViewerPolicy',
    'Read-only access to all resources',
    '{"version":"2025-01-01","statements":[{"sid":"AllowRead","effect":"Allow","actions":["product:Read","brand:Read","order:Read","customer:Read","inventory:Read","pricing:Read","cart:Read","fulfillment:Read","promotion:Read","region:Read","notification:Read"],"resources":["*"]}]}',
    true
);

-- ── Seed built-in roles ───────────────────────────────────────────────────────

INSERT INTO iam.role (id, name, description, is_system) VALUES
('role_system_super_admin', 'super_admin', 'Unrestricted access to everything', true),
('role_system_admin',       'admin',       'Manage catalog, orders and operations', true),
('role_system_operator',    'operator',    'Fulfill orders and manage stock', true),
('role_system_viewer',      'viewer',      'Read-only access', true);

-- ── Wire built-in role → policy ───────────────────────────────────────────────

INSERT INTO iam.role_policy (role_id, policy_id) VALUES
('role_system_super_admin', 'pol_system_super_admin'),
('role_system_admin',       'pol_system_admin'),
('role_system_operator',    'pol_system_operator'),
('role_system_viewer',      'pol_system_viewer');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS iam.api_key_policy;
DROP TABLE IF EXISTS iam.principal_role;
DROP TABLE IF EXISTS iam.role_policy;
DROP TABLE IF EXISTS iam.policy;
DROP TABLE IF EXISTS iam.role;
-- +goose StatementEnd
