package handler

import (
	"net/http"

	"ecommerce/apps/api/internal/middleware"
	"ecommerce/modules/auth"
	"ecommerce/modules/identity"
	"ecommerce/packages/httpx"
	"ecommerce/packages/types"

	"github.com/go-chi/chi/v5"
)

type IdentityHandler struct {
	svc     identity.Service
	authSvc auth.Service
}

func NewIdentityHandler(svc identity.Service, authSvc auth.Service) *IdentityHandler {
	return &IdentityHandler{svc: svc, authSvc: authSvc}
}

func (h *IdentityHandler) Routes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(h.authSvc))
		r.Get("/me", h.GetProfile)
		r.Post("/me/avatar", h.UploadAvatar)
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
	httpx.JSON(w, r, http.StatusOK, map[string]any{
		"user": u,
	})
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

	httpx.JSON(w, r, http.StatusOK, map[string]any{
		"user": u,
	})
}
