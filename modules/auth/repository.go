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
