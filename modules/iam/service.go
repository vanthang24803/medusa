package iam

import (
	"context"
	"fmt"
	"time"

	"ecommerce/packages/types"
)

type Service interface {
	// Evaluate returns true if the auth identity is allowed to perform action on resource.
	// resource follows service-level format: "product/*", "order/*", "*".
	Evaluate(ctx context.Context, authIdentityID, action, resource string) (bool, error)

	// EvaluateAPIKey returns true if the API key is allowed to perform action on resource.
	EvaluateAPIKey(ctx context.Context, apiKeyID, action, resource string) (bool, error)

	// Role management
	ListRoles(ctx context.Context) ([]Role, error)
	GetRole(ctx context.Context, id string) (*Role, error)
	CreateRole(ctx context.Context, req CreateRoleReq) (*Role, error)
	DeleteRole(ctx context.Context, id string) error
	AttachPolicy(ctx context.Context, roleID, policyID string) error
	DetachPolicy(ctx context.Context, roleID, policyID string) error

	// Policy management
	ListPolicies(ctx context.Context) ([]Policy, error)
	GetPolicy(ctx context.Context, id string) (*Policy, error)
	CreatePolicy(ctx context.Context, req CreatePolicyReq) (*Policy, error)
	UpdatePolicy(ctx context.Context, id string, req UpdatePolicyReq) (*Policy, error)
	DeletePolicy(ctx context.Context, id string) error

	// Principal ↔ Role assignment
	AssignRole(ctx context.Context, authIdentityID, roleID string) error
	UnassignRole(ctx context.Context, authIdentityID, roleID string) error
	GetRolesForPrincipal(ctx context.Context, authIdentityID string) ([]Role, error)

	// API Key ↔ Policy
	AttachPolicyToAPIKey(ctx context.Context, apiKeyID, policyID string) error
	DetachPolicyFromAPIKey(ctx context.Context, apiKeyID, policyID string) error
}

type service struct {
	repo *Repository
}

func NewService(repo *Repository) Service {
	return &service{repo: repo}
}

// ── Evaluation ───────────────────────────────────────────────────────────────

func (s *service) Evaluate(ctx context.Context, authIdentityID, action, resource string) (bool, error) {
	docs, err := s.repo.GetDocumentsForPrincipal(ctx, authIdentityID)
	if err != nil {
		return false, fmt.Errorf("load policies: %w", err)
	}
	return Evaluate(docs, action, resource), nil
}

func (s *service) EvaluateAPIKey(ctx context.Context, apiKeyID, action, resource string) (bool, error) {
	docs, err := s.repo.GetDocumentsForAPIKey(ctx, apiKeyID)
	if err != nil {
		return false, fmt.Errorf("load api key policies: %w", err)
	}
	return Evaluate(docs, action, resource), nil
}

// ── Roles ────────────────────────────────────────────────────────────────────

func (s *service) ListRoles(ctx context.Context) ([]Role, error) {
	return s.repo.GetRoles(ctx)
}

func (s *service) GetRole(ctx context.Context, id string) (*Role, error) {
	return s.repo.GetRoleByID(ctx, id)
}

func (s *service) CreateRole(ctx context.Context, req CreateRoleReq) (*Role, error) {
	if req.Name == "" {
		return nil, types.NewValidation("role name is required")
	}
	now := time.Now().UTC()
	role := &Role{
		ID:          types.GenerateID("role"),
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.InsertRole(ctx, role); err != nil {
		return nil, fmt.Errorf("insert role: %w", err)
	}
	return role, nil
}

func (s *service) DeleteRole(ctx context.Context, id string) error {
	role, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}
	if role.IsSystem {
		return types.NewForbidden("system roles cannot be deleted")
	}
	return s.repo.DeleteRole(ctx, id)
}

func (s *service) AttachPolicy(ctx context.Context, roleID, policyID string) error {
	if _, err := s.repo.GetRoleByID(ctx, roleID); err != nil {
		return err
	}
	if _, err := s.repo.GetPolicyByID(ctx, policyID); err != nil {
		return err
	}
	return s.repo.AttachPolicyToRole(ctx, roleID, policyID)
}

func (s *service) DetachPolicy(ctx context.Context, roleID, policyID string) error {
	return s.repo.DetachPolicyFromRole(ctx, roleID, policyID)
}

// ── Policies ─────────────────────────────────────────────────────────────────

func (s *service) ListPolicies(ctx context.Context) ([]Policy, error) {
	return s.repo.GetPolicies(ctx)
}

func (s *service) GetPolicy(ctx context.Context, id string) (*Policy, error) {
	return s.repo.GetPolicyByID(ctx, id)
}

func (s *service) CreatePolicy(ctx context.Context, req CreatePolicyReq) (*Policy, error) {
	if req.Name == "" {
		return nil, types.NewValidation("policy name is required")
	}
	if len(req.Document.Statements) == 0 {
		return nil, types.NewValidation("policy document must have at least one statement")
	}
	now := time.Now().UTC()
	p := &Policy{
		ID:          types.GenerateID("pol"),
		Name:        req.Name,
		Description: req.Description,
		Document:    req.Document,
		IsSystem:    false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.InsertPolicy(ctx, p); err != nil {
		return nil, fmt.Errorf("insert policy: %w", err)
	}
	return p, nil
}

func (s *service) UpdatePolicy(ctx context.Context, id string, req UpdatePolicyReq) (*Policy, error) {
	existing, err := s.repo.GetPolicyByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing.IsSystem {
		return nil, types.NewForbidden("system policies cannot be modified")
	}
	if err := s.repo.UpdatePolicy(ctx, id, req.Description, req.Document); err != nil {
		return nil, fmt.Errorf("update policy: %w", err)
	}
	return s.repo.GetPolicyByID(ctx, id)
}

func (s *service) DeletePolicy(ctx context.Context, id string) error {
	existing, err := s.repo.GetPolicyByID(ctx, id)
	if err != nil {
		return err
	}
	if existing.IsSystem {
		return types.NewForbidden("system policies cannot be deleted")
	}
	return s.repo.DeletePolicy(ctx, id)
}

// ── Principal assignment ─────────────────────────────────────────────────────

func (s *service) AssignRole(ctx context.Context, authIdentityID, roleID string) error {
	if _, err := s.repo.GetRoleByID(ctx, roleID); err != nil {
		return err
	}
	return s.repo.AssignRoleToPrincipal(ctx, authIdentityID, roleID)
}

func (s *service) UnassignRole(ctx context.Context, authIdentityID, roleID string) error {
	return s.repo.UnassignRoleFromPrincipal(ctx, authIdentityID, roleID)
}

func (s *service) GetRolesForPrincipal(ctx context.Context, authIdentityID string) ([]Role, error) {
	return s.repo.GetRolesForPrincipal(ctx, authIdentityID)
}

// ── API Key policies ─────────────────────────────────────────────────────────

func (s *service) AttachPolicyToAPIKey(ctx context.Context, apiKeyID, policyID string) error {
	if _, err := s.repo.GetPolicyByID(ctx, policyID); err != nil {
		return err
	}
	return s.repo.AttachPolicyToAPIKey(ctx, apiKeyID, policyID)
}

func (s *service) DetachPolicyFromAPIKey(ctx context.Context, apiKeyID, policyID string) error {
	return s.repo.DetachPolicyFromAPIKey(ctx, apiKeyID, policyID)
}
