package iam

import (
	"context"
	"database/sql"
	"errors"

	"ecommerce/packages/actor"
	"ecommerce/packages/db"
	"ecommerce/packages/types"
)

type Repository struct {
	db *db.DB
}

func NewRepository(database *db.DB) *Repository {
	return &Repository{db: database}
}

// ── Roles ────────────────────────────────────────────────────────────────────

const roleColumns = `id, name, description, is_system, last_updated_by, created_at, updated_at`

func (r *Repository) GetRoles(ctx context.Context) ([]Role, error) {
	var roles []Role
	err := r.db.Reader(ctx).SelectContext(ctx, &roles,
		`SELECT `+roleColumns+` FROM iam.role ORDER BY name`)
	return roles, err
}

func (r *Repository) GetRoleByID(ctx context.Context, id string) (*Role, error) {
	var role Role
	err := r.db.Reader(ctx).GetContext(ctx, &role,
		`SELECT `+roleColumns+` FROM iam.role WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("role")
	}
	return &role, err
}

func (r *Repository) GetRoleByName(ctx context.Context, name string) (*Role, error) {
	var role Role
	err := r.db.Reader(ctx).GetContext(ctx, &role,
		`SELECT `+roleColumns+` FROM iam.role WHERE name = $1`, name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("role")
	}
	return &role, err
}

func (r *Repository) InsertRole(ctx context.Context, role *Role) error {
	_, err := r.db.Writer(ctx).NamedExec(
		`INSERT INTO iam.role (id, name, description, is_system, last_updated_by, created_at, updated_at)
		 VALUES (:id, :name, :description, :is_system, :last_updated_by, :created_at, :updated_at)`, role)
	return err
}

func (r *Repository) DeleteRole(ctx context.Context, id string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx, `DELETE FROM iam.role WHERE id = $1 AND is_system = false`, id)
	return err
}

// ── Policies ─────────────────────────────────────────────────────────────────

const policyColumns = `id, name, description, document, is_system, last_updated_by, created_at, updated_at`

func (r *Repository) GetPolicies(ctx context.Context) ([]Policy, error) {
	var policies []Policy
	err := r.db.Reader(ctx).SelectContext(ctx, &policies,
		`SELECT `+policyColumns+` FROM iam.policy ORDER BY name`)
	return policies, err
}

func (r *Repository) GetPolicyByID(ctx context.Context, id string) (*Policy, error) {
	var policy Policy
	err := r.db.Reader(ctx).GetContext(ctx, &policy,
		`SELECT `+policyColumns+` FROM iam.policy WHERE id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, types.NewNotFound("policy")
	}
	return &policy, err
}

func (r *Repository) InsertPolicy(ctx context.Context, p *Policy) error {
	_, err := r.db.Writer(ctx).NamedExec(
		`INSERT INTO iam.policy (id, name, description, document, is_system, last_updated_by, created_at, updated_at)
		 VALUES (:id, :name, :description, :document, :is_system, :last_updated_by, :created_at, :updated_at)`, p)
	return err
}

func (r *Repository) UpdatePolicy(ctx context.Context, id string, description *string, doc *PolicyDocument) error {
	updatedBy := actor.Get(ctx)
	if description != nil {
		if _, err := r.db.Writer(ctx).ExecContext(ctx,
			`UPDATE iam.policy SET description = $1, last_updated_by = $2, updated_at = now() WHERE id = $3 AND is_system = false`,
			*description, updatedBy, id); err != nil {
			return err
		}
	}
	if doc != nil {
		if _, err := r.db.Writer(ctx).ExecContext(ctx,
			`UPDATE iam.policy SET document = $1, last_updated_by = $2, updated_at = now() WHERE id = $3 AND is_system = false`,
			doc, updatedBy, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) DeletePolicy(ctx context.Context, id string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`DELETE FROM iam.policy WHERE id = $1 AND is_system = false`, id)
	return err
}

// ── Role ↔ Policy ────────────────────────────────────────────────────────────

func (r *Repository) AttachPolicyToRole(ctx context.Context, roleID, policyID string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`INSERT INTO iam.role_policy (role_id, policy_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		roleID, policyID)
	return err
}

func (r *Repository) DetachPolicyFromRole(ctx context.Context, roleID, policyID string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`DELETE FROM iam.role_policy WHERE role_id = $1 AND policy_id = $2`, roleID, policyID)
	return err
}

func (r *Repository) GetPoliciesForRole(ctx context.Context, roleID string) ([]Policy, error) {
	var policies []Policy
	err := r.db.Reader(ctx).SelectContext(ctx, &policies,
		`SELECT p.id, p.name, p.description, p.document, p.is_system, p.last_updated_by, p.created_at, p.updated_at
		 FROM iam.policy p
		 JOIN iam.role_policy rp ON rp.policy_id = p.id
		 WHERE rp.role_id = $1`, roleID)
	return policies, err
}

// ── Principal ↔ Role ─────────────────────────────────────────────────────────

func (r *Repository) AssignRoleToPrincipal(ctx context.Context, authIdentityID, roleID string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`INSERT INTO iam.principal_role (auth_identity_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		authIdentityID, roleID)
	return err
}

func (r *Repository) UnassignRoleFromPrincipal(ctx context.Context, authIdentityID, roleID string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`DELETE FROM iam.principal_role WHERE auth_identity_id = $1 AND role_id = $2`,
		authIdentityID, roleID)
	return err
}

// GetDocumentsForPrincipal loads all PolicyDocuments for an auth identity via role chain.
func (r *Repository) GetDocumentsForPrincipal(ctx context.Context, authIdentityID string) ([]PolicyDocument, error) {
	var policies []Policy
	err := r.db.Reader(ctx).SelectContext(ctx, &policies,
		`SELECT p.id, p.name, p.description, p.document, p.is_system, p.last_updated_by, p.created_at, p.updated_at
		 FROM iam.policy p
		 JOIN iam.role_policy rp ON rp.policy_id = p.id
		 JOIN iam.principal_role pr ON pr.role_id = rp.role_id
		 WHERE pr.auth_identity_id = $1`, authIdentityID)
	if err != nil {
		return nil, err
	}
	docs := make([]PolicyDocument, len(policies))
	for i, p := range policies {
		docs[i] = p.Document
	}
	return docs, nil
}

func (r *Repository) GetRolesForPrincipal(ctx context.Context, authIdentityID string) ([]Role, error) {
	var roles []Role
	err := r.db.Reader(ctx).SelectContext(ctx, &roles,
		`SELECT ro.id, ro.name, ro.description, ro.is_system, ro.last_updated_by, ro.created_at, ro.updated_at
		 FROM iam.role ro
		 JOIN iam.principal_role pr ON pr.role_id = ro.id
		 WHERE pr.auth_identity_id = $1`, authIdentityID)
	return roles, err
}

// ── API Key ↔ Policy ─────────────────────────────────────────────────────────

func (r *Repository) AttachPolicyToAPIKey(ctx context.Context, apiKeyID, policyID string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`INSERT INTO iam.api_key_policy (api_key_id, policy_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		apiKeyID, policyID)
	return err
}

func (r *Repository) DetachPolicyFromAPIKey(ctx context.Context, apiKeyID, policyID string) error {
	_, err := r.db.Writer(ctx).ExecContext(ctx,
		`DELETE FROM iam.api_key_policy WHERE api_key_id = $1 AND policy_id = $2`, apiKeyID, policyID)
	return err
}

// GetDocumentsForAPIKey loads all PolicyDocuments directly attached to an API key.
func (r *Repository) GetDocumentsForAPIKey(ctx context.Context, apiKeyID string) ([]PolicyDocument, error) {
	var policies []Policy
	err := r.db.Reader(ctx).SelectContext(ctx, &policies,
		`SELECT p.id, p.name, p.description, p.document, p.is_system, p.last_updated_by, p.created_at, p.updated_at
		 FROM iam.policy p
		 JOIN iam.api_key_policy akp ON akp.policy_id = p.id
		 WHERE akp.api_key_id = $1`, apiKeyID)
	if err != nil {
		return nil, err
	}
	docs := make([]PolicyDocument, len(policies))
	for i, p := range policies {
		docs[i] = p.Document
	}
	return docs, nil
}
