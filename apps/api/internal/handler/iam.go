package handler

import (
	"net/http"

	"ecommerce/apps/api/internal/middleware"
	"ecommerce/modules/auth"
	"ecommerce/modules/iam"
	"ecommerce/packages/httpx"
	"ecommerce/packages/types"

	"github.com/go-chi/chi/v5"
)

type IAMHandler struct {
	svc     iam.Service
	authSvc auth.Service
}

func NewIAMHandler(svc iam.Service, authSvc auth.Service) *IAMHandler {
	return &IAMHandler{svc: svc, authSvc: authSvc}
}

func (h *IAMHandler) Routes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(h.authSvc))
		r.Use(middleware.RequirePermission(h.svc, iam.ActionIAMManage))

		// Roles
		r.Get("/roles", h.ListRoles)
		r.Post("/roles", h.CreateRole)
		r.Delete("/roles/{roleId}", h.DeleteRole)
		r.Post("/roles/{roleId}/policies/{policyId}", h.AttachPolicy)
		r.Delete("/roles/{roleId}/policies/{policyId}", h.DetachPolicy)

		// Policies
		r.Get("/policies", h.ListPolicies)
		r.Post("/policies", h.CreatePolicy)
		r.Post("/policies/{policyId}/update", h.UpdatePolicy)
		r.Delete("/policies/{policyId}", h.DeletePolicy)

		// Principal role assignment
		r.Get("/principals/{authId}/roles", h.GetPrincipalRoles)
		r.Post("/principals/{authId}/roles/{roleId}", h.AssignRole)
		r.Delete("/principals/{authId}/roles/{roleId}", h.UnassignRole)

		// API Key policy
		r.Post("/api-keys/{keyId}/policies/{policyId}", h.AttachPolicyToAPIKey)
		r.Delete("/api-keys/{keyId}/policies/{policyId}", h.DetachPolicyFromAPIKey)
	})

	// API key CRUD — requires auth + iam:Manage
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(h.authSvc))
		r.Use(middleware.RequirePermission(h.svc, iam.ActionIAMManage))
		r.Get("/api-keys", h.ListAPIKeys)
		r.Post("/api-keys", h.CreateAPIKey)
		r.Delete("/api-keys/{keyId}", h.RevokeAPIKey)
	})
}

// ── Roles ────────────────────────────────────────────────────────────────────

func (h *IAMHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.svc.ListRoles(r.Context())
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, httpx.Response("roles", roles))
}

func (h *IAMHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req iam.CreateRoleReq
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	role, err := h.svc.CreateRole(r.Context(), req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusCreated, httpx.Response("role", role))
}

func (h *IAMHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roleId")
	if err := h.svc.DeleteRole(r.Context(), id); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "role deleted"})
}

func (h *IAMHandler) AttachPolicy(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleId")
	policyID := chi.URLParam(r, "policyId")
	if err := h.svc.AttachPolicy(r.Context(), roleID, policyID); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "policy attached"})
}

func (h *IAMHandler) DetachPolicy(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleId")
	policyID := chi.URLParam(r, "policyId")
	if err := h.svc.DetachPolicy(r.Context(), roleID, policyID); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "policy detached"})
}

// ── Policies ─────────────────────────────────────────────────────────────────

func (h *IAMHandler) ListPolicies(w http.ResponseWriter, r *http.Request) {
	policies, err := h.svc.ListPolicies(r.Context())
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, httpx.Response("policies", policies))
}

func (h *IAMHandler) CreatePolicy(w http.ResponseWriter, r *http.Request) {
	var req iam.CreatePolicyReq
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	p, err := h.svc.CreatePolicy(r.Context(), req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusCreated, httpx.Response("policy", p))
}

func (h *IAMHandler) UpdatePolicy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "policyId")
	var req iam.UpdatePolicyReq
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	p, err := h.svc.UpdatePolicy(r.Context(), id, req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, httpx.Response("policy", p))
}

func (h *IAMHandler) DeletePolicy(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "policyId")
	if err := h.svc.DeletePolicy(r.Context(), id); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "policy deleted"})
}

// ── Principal role assignment ─────────────────────────────────────────────────

func (h *IAMHandler) GetPrincipalRoles(w http.ResponseWriter, r *http.Request) {
	authID := chi.URLParam(r, "authId")
	roles, err := h.svc.GetRolesForPrincipal(r.Context(), authID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, httpx.Response("roles", roles))
}

func (h *IAMHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	authID := chi.URLParam(r, "authId")
	roleID := chi.URLParam(r, "roleId")
	if err := h.svc.AssignRole(r.Context(), authID, roleID); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "role assigned"})
}

func (h *IAMHandler) UnassignRole(w http.ResponseWriter, r *http.Request) {
	authID := chi.URLParam(r, "authId")
	roleID := chi.URLParam(r, "roleId")
	if err := h.svc.UnassignRole(r.Context(), authID, roleID); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "role unassigned"})
}

// ── API Key ───────────────────────────────────────────────────────────────────

func (h *IAMHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	authID := middleware.GetAuthIdentityID(r.Context())
	keys, err := h.authSvc.ListAPIKeys(r.Context(), authID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, httpx.Response("apiKeys", keys))
}

func (h *IAMHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	authID := middleware.GetAuthIdentityID(r.Context())
	if authID == "" {
		httpx.Error(w, r, types.ErrUnauthorized)
		return
	}
	var req auth.CreateAPIKeyReq
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	k, err := h.authSvc.CreateAPIKey(r.Context(), authID, req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusCreated, httpx.Response("apiKey", k))
}

func (h *IAMHandler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "keyId")
	if err := h.authSvc.RevokeAPIKey(r.Context(), id); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "api key revoked"})
}

func (h *IAMHandler) AttachPolicyToAPIKey(w http.ResponseWriter, r *http.Request) {
	keyID := chi.URLParam(r, "keyId")
	policyID := chi.URLParam(r, "policyId")
	if err := h.svc.AttachPolicyToAPIKey(r.Context(), keyID, policyID); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "policy attached to api key"})
}

func (h *IAMHandler) DetachPolicyFromAPIKey(w http.ResponseWriter, r *http.Request) {
	keyID := chi.URLParam(r, "keyId")
	policyID := chi.URLParam(r, "policyId")
	if err := h.svc.DetachPolicyFromAPIKey(r.Context(), keyID, policyID); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "policy detached from api key"})
}
