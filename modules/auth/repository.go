package auth

import (
	"context"
	"database/sql"
	"errors"

	"ecommerce/packages/db"
	"ecommerce/packages/types"
)

type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

func (r *Repository) Ping(ctx context.Context) error {
	var one int
	return r.db.Reader(ctx).GetContext(ctx, &one, "SELECT 1")
}

func (r *Repository) InsertAuthIdentity(ctx context.Context, a *AuthIdentity) error {
	_, err := r.db.Writer(ctx).NamedExec(`
		INSERT INTO auth.auth_identity (id, app_metadata, created_at, updated_at)
		VALUES (:id, :app_metadata, :created_at, :updated_at)`, a)
	return err
}

func (r *Repository) GetAuthIdentityByID(ctx context.Context, id string) (*AuthIdentity, error) {
	var a AuthIdentity
	err := r.db.Reader(ctx).GetContext(ctx, &a,
		`SELECT id, app_metadata, created_at, updated_at FROM auth.auth_identity WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("auth identity")
	}
	return &a, err
}

func (r *Repository) InsertProviderIdentity(ctx context.Context, p *ProviderIdentity) error {
	_, err := r.db.Writer(ctx).NamedExec(`
		INSERT INTO auth.auth_provider_identity
			(id, auth_identity_id, provider, entity_id, provider_metadata, user_metadata, created_at, updated_at)
		VALUES (:id, :auth_identity_id, :provider, :entity_id, :provider_metadata, :user_metadata, :created_at, :updated_at)`, p)
	return err
}

func (r *Repository) GetProviderByAuthIdentityID(ctx context.Context, authIdentityID string) (*ProviderIdentity, error) {
	var p ProviderIdentity
	err := r.db.Reader(ctx).GetContext(ctx, &p,
		`SELECT id, auth_identity_id, provider, entity_id, provider_metadata, user_metadata, created_at, updated_at
		 FROM auth.auth_provider_identity WHERE auth_identity_id = $1 AND provider = 'emailpass'`, authIdentityID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("provider identity")
	}
	return &p, err
}

// IsAdminActive checks identity.user.status directly via SQL to avoid circular import.
// Returns true if no identity user record exists (regular customer) or status = 'active'.
func (r *Repository) IsAdminActive(ctx context.Context, authIdentityID string) (bool, error) {
	var status string
	err := r.db.Reader(ctx).GetContext(ctx, &status,
		`SELECT status FROM identity."user" WHERE id = $1 AND deleted_at IS NULL`, authIdentityID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return true, nil // not an admin (customer), allow
		}
		return false, err
	}
	return status == "active", nil
}

func (r *Repository) UpdateProviderMetadata(ctx context.Context, id string, metadata []byte) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE auth.auth_provider_identity SET provider_metadata = $1, updated_at = now() WHERE id = $2`,
		metadata, id)
	return err
}

func (r *Repository) GetProviderByEntityID(ctx context.Context, provider, entityID string) (*ProviderIdentity, error) {
	var p ProviderIdentity
	err := r.db.Reader(ctx).GetContext(ctx, &p,
		`SELECT id, auth_identity_id, provider, entity_id, provider_metadata, user_metadata, created_at, updated_at
		 FROM auth.auth_provider_identity WHERE provider = $1 AND entity_id = $2`, provider, entityID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("provider identity")
	}
	return &p, err
}

func (r *Repository) InsertAuthToken(ctx context.Context, t *AuthToken) error {
	_, err := r.db.Writer(ctx).NamedExec(`
		INSERT INTO auth.auth_token (id, auth_identity_id, token_hash, type, expires_at, created_at)
		VALUES (:id, :auth_identity_id, :token_hash, :type, :expires_at, :created_at)`, t)
	return err
}

func (r *Repository) GetAuthTokenByHash(ctx context.Context, hash string) (*AuthToken, error) {
	var t AuthToken
	err := r.db.Reader(ctx).GetContext(ctx, &t,
		`SELECT id, auth_identity_id, token_hash, type, expires_at, revoked_at, created_at
		 FROM auth.auth_token WHERE token_hash = $1`, hash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("token")
	}
	return &t, err
}

func (r *Repository) RevokeAuthToken(ctx context.Context, id string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE auth.auth_token SET revoked_at = now() WHERE id = $1 AND revoked_at IS NULL`, id)
	return err
}

func (r *Repository) RevokeAuthTokensByAuthID(ctx context.Context, authIdentityID string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE auth.auth_token SET revoked_at = now() WHERE auth_identity_id = $1 AND revoked_at IS NULL`, authIdentityID)
	return err
}

// ── API Key ──────────────────────────────────────────────────────────────────

func (r *Repository) InsertAPIKey(ctx context.Context, k *APIKey) error {
	_, err := r.db.Writer(ctx).NamedExec(`
		INSERT INTO auth.api_key (id, token, type, title, created_by, created_at, updated_at)
		VALUES (:id, :token, :type, :title, :created_by, :created_at, :updated_at)`, k)
	return err
}

func (r *Repository) GetAPIKeyByToken(ctx context.Context, token string) (*APIKey, error) {
	var k APIKey
	err := r.db.Reader(ctx).GetContext(ctx, &k,
		`SELECT id, token, type, title, last_used_at, revoked_at, created_by, created_at, updated_at
		 FROM auth.api_key WHERE token = $1 AND revoked_at IS NULL`, token)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("api key")
	}
	return &k, err
}

func (r *Repository) GetAPIKeyByID(ctx context.Context, id string) (*APIKey, error) {
	var k APIKey
	err := r.db.Reader(ctx).GetContext(ctx, &k,
		`SELECT id, token, type, title, last_used_at, revoked_at, created_by, created_at, updated_at
		 FROM auth.api_key WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("api key")
	}
	return &k, err
}

func (r *Repository) ListAPIKeysByCreator(ctx context.Context, createdBy string) ([]APIKey, error) {
	var keys []APIKey
	err := r.db.Reader(ctx).SelectContext(ctx, &keys,
		`SELECT id, token, type, title, last_used_at, revoked_at, created_by, created_at, updated_at
		 FROM auth.api_key WHERE created_by = $1 ORDER BY created_at DESC`, createdBy)
	return keys, err
}

func (r *Repository) RevokeAPIKey(ctx context.Context, id string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE auth.api_key SET revoked_at = now(), updated_at = now() WHERE id = $1 AND revoked_at IS NULL`, id)
	return err
}

func (r *Repository) TouchAPIKeyLastUsed(ctx context.Context, id string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`UPDATE auth.api_key SET last_used_at = now() WHERE id = $1`, id)
	return err
}
