package handler

import (
	"net/http"

	"ecommerce/apps/api/internal/middleware"
	"ecommerce/modules/auth"
	"ecommerce/modules/customer"
	"ecommerce/packages/httpx"
	"ecommerce/packages/types"

	"github.com/go-chi/chi/v5"
)

type CustomerHandler struct {
	svc     customer.Service
	authSvc auth.Service
}

func NewCustomerHandler(svc customer.Service, authSvc auth.Service) *CustomerHandler {
	return &CustomerHandler{svc: svc, authSvc: authSvc}
}

func (h *CustomerHandler) Routes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(h.authSvc))
		r.Get("/me", h.GetProfile)
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
	httpx.JSON(w, r, http.StatusOK, map[string]any{
		"customer": c,
	})
}
