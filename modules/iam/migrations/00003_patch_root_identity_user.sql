-- +goose Up
-- +goose StatementBegin
INSERT INTO identity.user (id, email, first_name, last_name, avatar_url, metadata, created_at, updated_at, deleted_at)
VALUES (
    'auth_019ee8db-8b26-737b-b2f5-0cfa5b3edb30',
    'root@admin.com',
    'Root',
    'Admin',
    NULL,
    '{}',
    '2026-06-21T06:26:02.662246107Z',
    '2026-06-21T06:26:02.662246107Z',
    NULL
)
ON CONFLICT (id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM identity.user WHERE id = 'auth_019ee8db-8b26-737b-b2f5-0cfa5b3edb30';
-- +goose StatementEnd
