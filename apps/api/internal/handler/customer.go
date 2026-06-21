package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"ecommerce/apps/api/internal/middleware"
	"ecommerce/modules/auth"
	"ecommerce/modules/customer"
	"ecommerce/packages/httpx"
	"ecommerce/packages/types"

	"github.com/go-chi/chi/v5"
)

const maxUploadMemory = 5 << 20 // 5 MB multipart buffer

type CustomerHandler struct {
	svc     customer.Service
	authSvc auth.Service
	log     *zap.Logger
}

func NewCustomerHandler(svc customer.Service, authSvc auth.Service, log *zap.Logger) *CustomerHandler {
	return &CustomerHandler{svc: svc, authSvc: authSvc, log: log}
}

func (h *CustomerHandler) Routes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(h.authSvc))
		r.Get("/me", h.GetProfile)
		r.Post("/me/update", h.UpdateProfile)
		r.Post("/me/avatar", h.UploadAvatar)
		r.Post("/me/password", h.UpdatePassword)
	})
}

func (h *CustomerHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	customerID := middleware.GetCustomerID(r.Context())
	if customerID == "" {
		httpx.Error(w, r, types.ErrUnauthorized)
		return
	}
	c, err := h.svc.GetProfile(r.Context(), customerID)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, httpx.Response("customer", c))
}

func (h *CustomerHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	customerID := middleware.GetCustomerID(r.Context())
	if customerID == "" {
		httpx.Error(w, r, types.ErrUnauthorized)
		return
	}

	var req customer.UpdateCustomerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, r, types.NewValidation("invalid request body"))
		return
	}

	c, err := h.svc.UpdateProfile(r.Context(), customerID, &req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, httpx.Response("customer", c))
}

func (h *CustomerHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	authIdentityID := middleware.GetAuthIdentityID(r.Context())
	if authIdentityID == "" {
		httpx.Error(w, r, types.ErrUnauthorized)
		return
	}

	var req auth.UpdatePasswordReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpx.Error(w, r, types.NewValidation("invalid request body"))
		return
	}

	if err := h.authSvc.UpdatePassword(r.Context(), authIdentityID, req.CurrentPassword, req.NewPassword); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]string{"message": "password updated successfully"})
}

func (h *CustomerHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	customerID := middleware.GetCustomerID(r.Context())
	if customerID == "" {
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

	c, err := h.svc.UploadAvatar(r.Context(), customerID, file, header.Size, contentType)
	if err != nil {
		h.log.Error("upload avatar failed", zap.String("customerID", customerID), zap.Error(err))
		httpx.Error(w, r, err)
		return
	}

	httpx.JSON(w, r, http.StatusOK, httpx.Response("customer", c))
}
