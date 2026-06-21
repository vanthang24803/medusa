package handler

import (
	"net/http"

	"ecommerce/apps/api/internal/middleware"
	"ecommerce/modules/auth"
	"ecommerce/modules/brand"
	"ecommerce/modules/iam"
	"ecommerce/packages/httpx"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type BrandHandler struct {
	svc     brand.Service
	authSvc auth.Service
	iamSvc  iam.Service
	log     *zap.Logger
}

func NewBrandHandler(svc brand.Service, authSvc auth.Service, iamSvc iam.Service, log *zap.Logger) *BrandHandler {
	return &BrandHandler{svc: svc, authSvc: authSvc, iamSvc: iamSvc, log: log}
}

func (h *BrandHandler) Routes(r chi.Router) {
	// Public read endpoints
	r.Get("/", h.List)
	r.Get("/{brandId}", h.GetByID)
	r.Get("/slug/{slug}", h.GetBySlug)

	// Admin-only write endpoints
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAnyAuth(h.authSvc, h.authSvc))
		r.Use(middleware.RequireAnyPermission(h.iamSvc, iam.ActionBrandCreate, iam.ActionBrandManage))
		r.Post("/", h.Create)
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAnyAuth(h.authSvc, h.authSvc))
		r.Use(middleware.RequireAnyPermission(h.iamSvc, iam.ActionBrandUpdate, iam.ActionBrandManage))
		r.Post("/{brandId}/update", h.Update)
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAnyAuth(h.authSvc, h.authSvc))
		r.Use(middleware.RequireAnyPermission(h.iamSvc, iam.ActionBrandDelete, iam.ActionBrandManage))
		r.Delete("/{brandId}", h.Delete)
	})
}

func (h *BrandHandler) List(w http.ResponseWriter, r *http.Request) {
	var q brand.ListQuery
	if err := httpx.DecodeQuery(r, &q); err != nil {
		h.log.Error("failed to decode list query", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	items, err := h.svc.List(r.Context(), q)
	if err != nil {
		h.log.Error("failed to list brands", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"brands": items})
}

func (h *BrandHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "brandId")
	b, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		h.log.Error("failed to get brand by ID", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"brand": b})
}

func (h *BrandHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	b, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
		h.log.Error("failed to get brand by slug", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"brand": b})
}

func (h *BrandHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req brand.CreateInput
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	b, err := h.svc.Create(r.Context(), req)
	if err != nil {
		h.log.Error("failed to create brand", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusCreated, map[string]any{"brand": b})
}

func (h *BrandHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "brandId")
	var req brand.UpdateInput
	if err := httpx.DecodeAndValidate(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	b, err := h.svc.Update(r.Context(), id, req)
	if err != nil {
		h.log.Error("failed to update brand", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"brand": b})
}

func (h *BrandHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "brandId")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		h.log.Error("failed to delete brand", zap.Error(err))
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"message": "brand deleted"})
}
