package handler

import (
	"net/http"

	"ecommerce/modules/brand"
	"ecommerce/packages/httpx"

	"github.com/go-chi/chi/v5"
)

type BrandHandler struct {
	svc brand.Service
}

func NewBrandHandler(svc brand.Service) *BrandHandler {
	return &BrandHandler{svc: svc}
}

func (h *BrandHandler) Routes(r chi.Router) {
	r.Get("/", h.List)
	r.Get("/{brandId}", h.GetByID)
	r.Get("/slug/{slug}", h.GetBySlug)
	r.Post("/", h.Create)
	r.Put("/{brandId}", h.Update)
	r.Delete("/{brandId}", h.Delete)
}

func (h *BrandHandler) List(w http.ResponseWriter, r *http.Request) {
	var q brand.ListQuery
	if err := httpx.DecodeQuery(r, &q); err != nil {
		httpx.Error(w, r, err)
		return
	}
	items, err := h.svc.List(r.Context(), q)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"brands": items})
}

func (h *BrandHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "brandId")
	b, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"brand": b})
}

func (h *BrandHandler) GetBySlug(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	b, err := h.svc.GetBySlug(r.Context(), slug)
	if err != nil {
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
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"brand": b})
}

func (h *BrandHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "brandId")
	if err := h.svc.Delete(r.Context(), id); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"message": "brand deleted"})
}
