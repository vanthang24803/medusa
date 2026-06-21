-- +goose Up
-- +goose StatementBegin

-- Auth identity
INSERT INTO auth.auth_identity (id, app_metadata, created_at, updated_at) VALUES (
    'auth_019ee8db-8b26-737b-b2f5-0cfa5b3edb30',
    '{"customerId":"cus_019ee8db-8b26-737f-a605-bab4212f29fe"}',
    '2026-06-21T06:26:02.662246107Z',
    '2026-06-21T06:26:02.662246107Z'
);

-- Provider identity (emailpass: root@admin.com / root)
INSERT INTO auth.auth_provider_identity (id, auth_identity_id, provider, entity_id, provider_metadata, user_metadata, created_at, updated_at) VALUES (
    'prov_019ee8db-8b26-737e-bfad-8fed929c8715',
    'auth_019ee8db-8b26-737b-b2f5-0cfa5b3edb30',
    'emailpass',
    'root@admin.com',
    '{"password":"$2a$10$A7Wi4OyWMhCXBoZOzgpELuxrbtw4o2jFDyTwvDZfmxOMUiZ59juD6"}',
    '{}',
    '2026-06-21T06:26:02.662246107Z',
    '2026-06-21T06:26:02.662246107Z'
);

-- Customer record
INSERT INTO customer.customer (id, email, first_name, last_name, phone, company_name, avatar_url, has_account, metadata, created_at, updated_at, deleted_at) VALUES (
    'cus_019ee8db-8b26-737f-a605-bab4212f29fe',
    'root@admin.com',
    'Root',
    'Admin',
    NULL, NULL, NULL,
    true,
    '{}',
    '2026-06-21T06:26:02.662246107Z',
    '2026-06-21T06:26:02.662246107Z',
    NULL
);

-- Identity user record (admin profile — id = auth_identity_id by design)
INSERT INTO identity.user (id, email, first_name, last_name, avatar_url, metadata, created_at, updated_at, deleted_at) VALUES (
    'auth_019ee8db-8b26-737b-b2f5-0cfa5b3edb30',
    'root@admin.com',
    'Root',
    'Admin',
    NULL,
    '{}',
    '2026-06-21T06:26:02.662246107Z',
    '2026-06-21T06:26:02.662246107Z',
    NULL
);

-- Assign super_admin role
INSERT INTO iam.principal_role (auth_identity_id, role_id, created_at) VALUES (
    'auth_019ee8db-8b26-737b-b2f5-0cfa5b3edb30',
    'role_system_super_admin',
    '2026-06-21T06:26:02.662246107Z'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM iam.principal_role          WHERE auth_identity_id = 'auth_019ee8db-8b26-737b-b2f5-0cfa5b3edb30';
DELETE FROM identity.user               WHERE id  = 'auth_019ee8db-8b26-737b-b2f5-0cfa5b3edb30';
DELETE FROM customer.customer           WHERE id  = 'cus_019ee8db-8b26-737f-a605-bab4212f29fe';
DELETE FROM auth.auth_provider_identity WHERE id  = 'prov_019ee8db-8b26-737e-bfad-8fed929c8715';
DELETE FROM auth.auth_identity          WHERE id  = 'auth_019ee8db-8b26-737b-b2f5-0cfa5b3edb30';
-- +goose StatementEnd
