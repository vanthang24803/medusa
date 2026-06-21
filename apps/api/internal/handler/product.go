package handler

import (
	"net/http"

	"ecommerce/apps/api/internal/middleware"
	"ecommerce/modules/auth"
	"ecommerce/modules/iam"
	"ecommerce/modules/product"
	"ecommerce/packages/httpx"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	svc     product.Service
	authSvc auth.Service
	iamSvc  iam.Service
}

func NewProductHandler(svc product.Service, authSvc auth.Service, iamSvc iam.Service) *ProductHandler {
	return &ProductHandler{svc: svc, authSvc: authSvc, iamSvc: iamSvc}
}

func (h *ProductHandler) Routes(r chi.Router) {
	// Public read endpoints
	r.Get("/", h.List)
	r.Get("/{productId}", h.GetByID)
	r.Get("/handle/{handle}", h.GetByHandle)

	// Admin-only write endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAnyAuth(h.authSvc, h.authSvc))
		r.Use(middleware.RequireAnyPermission(h.iamSvc, iam.ActionProductCreate, iam.ActionProductManage))
		r.Post("/", h.Create)
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAnyAuth(h.authSvc, h.authSvc))
		r.Use(middleware.RequireAnyPermission(h.iamSvc, iam.ActionProductDelete, iam.ActionProductManage))
		r.Delete("/{productId}", h.Delete)
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAnyAuth(h.authSvc, h.authSvc))
		r.Use(middleware.RequireAnyPermission(h.iamSvc, iam.ActionProductCreate, iam.ActionProductManage))
		r.Post("/{productId}/variants", h.CreateVariant)
	})
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	var q product.ListQuery
	if err := httpx.DecodeQuery(r, &q); err != nil {
		httpx.Error(w, r, err)
		return
	}
	resp, err := h.svc.List(r.Context(), q)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, resp)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "productId")
	p, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"product": p})
}

func (h *ProductHandler) GetByHandle(w http.ResponseWriter, r *http.Request) {
	handle := chi.URLParam(r, "handle")
	p, err := h.svc.GetByHandle(r.Context(), handle)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"product": p})
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req product.CreateInput
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	p, err := h.svc.Create(r.Context(), req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusCreated, map[string]any{"product": p})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "productId")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"message": "product deleted"})
}

func (h *ProductHandler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productId")
	var req product.CreateVariantInput
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	v, err := h.svc.CreateVariant(r.Context(), productID, req)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusCreated, map[string]any{"variant": v})
}
