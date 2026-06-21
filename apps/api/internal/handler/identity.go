package handler

import (
	"encoding/json"
	"net/http"

	"ecommerce/apps/api/internal/middleware"
	"ecommerce/modules/auth"
	"ecommerce/modules/iam"
	"ecommerce/modules/identity"
	"ecommerce/packages/httpx"
	"ecommerce/packages/types"

	"github.com/go-chi/chi/v5"
)

type IdentityHandler struct {
	svc     identity.Service
	authSvc auth.Service
	iamSvc  iam.Service
}

func NewIdentityHandler(svc identity.Service, authSvc auth.Service, iamSvc iam.Service) *IdentityHandler {
	return &IdentityHandler{svc: svc, authSvc: authSvc, iamSvc: iamSvc}
}

func (h *IdentityHandler) Routes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(h.authSvc))
		r.Get("/me", h.GetProfile)
		r.Post("/me/update", h.UpdateProfile)
		r.Post("/me/avatar", h.UploadAvatar)
	})

	// User management — requires user:Manage (assignable to any role)
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(h.authSvc))
		r.Use(middleware.RequirePermission(h.iamSvc, iam.ActionUserManage))
		r.Get("/admin", h.ListAdmins)
		r.Post("/admin/{userId}/ban", h.BanUser)
		r.Post("/admin/{userId}/unban", h.UnbanUser)
	})

	// Super-admin only — requires identity:Manage
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(h.authSvc))
		r.Use(middleware.RequirePermission(h.iamSvc, iam.ActionIdentityManage))
		r.Post("/admin", h.CreateAdmin)
		r.Delete("/admin/{userId}", h.RevokeAdmin)
	})
}

func (h *IdentityHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetAuthIdentityID(r.Context())
	if userID == "" {
		httpx.Error(w, r, types.ErrUnauthorized)
		return
	}
	u, err := h.svc.GetProfile(r.Context(), userID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"user": u})
}

func (h *IdentityHandler) CreateAdmin(w http.ResponseWriter, r *http.Request) {
	var req identity.CreateAdminReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, r, types.NewValidation("invalid request body"))
		return
	}

	u, err := h.svc.CreateAdmin(r.Context(), req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusCreated, map[string]any{"user": u})
}

func (h *IdentityHandler) ListAdmins(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.ListAdmins(r.Context())
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"users": users})
}

func (h *IdentityHandler) RevokeAdmin(w http.ResponseWriter, r *http.Request) {
	callerID := middleware.GetAuthIdentityID(r.Context())
	targetID := chi.URLParam(r, "userId")
	if targetID == "" {
		httpx.Error(w, r, types.NewValidation("userId is required"))
		return
	}
	if err := h.svc.RevokeAdmin(r.Context(), callerID, targetID); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "admin account closed"})
}

func (h *IdentityHandler) BanUser(w http.ResponseWriter, r *http.Request) {
	callerID := middleware.GetAuthIdentityID(r.Context())
	targetID := chi.URLParam(r, "userId")
	if targetID == "" {
		httpx.Error(w, r, types.NewValidation("userId is required"))
		return
	}
	if err := h.svc.BanUser(r.Context(), callerID, targetID); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "user banned"})
}

func (h *IdentityHandler) UnbanUser(w http.ResponseWriter, r *http.Request) {
	targetID := chi.URLParam(r, "userId")
	if targetID == "" {
		httpx.Error(w, r, types.NewValidation("userId is required"))
		return
	}
	if err := h.svc.UnbanUser(r.Context(), targetID); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "user unbanned"})
}

func (h *IdentityHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetAuthIdentityID(r.Context())
	if userID == "" {
		httpx.Error(w, r, types.ErrUnauthorized)
		return
	}

	var req identity.UpdateProfileReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, r, types.NewValidation("invalid request body"))
		return
	}

	u, err := h.svc.UpdateProfile(r.Context(), userID, &req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"user": u})
}

func (h *IdentityHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetAuthIdentityID(r.Context())
	if userID == "" {
		httpx.Error(w, r, types.ErrUnauthorized)
		return
	}

	if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
		httpx.Error(w, r, types.NewValidation("invalid multipart form"))
		return
	}

	file, header, err := r.FormFile("avatar")
	if err != nil {
		httpx.Error(w, r, types.NewValidation("field 'avatar' is required"))
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}

	u, err := h.svc.UploadAvatar(r.Context(), userID, file, header.Size, contentType)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"user": u})
}
