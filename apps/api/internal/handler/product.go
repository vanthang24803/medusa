package handler

import (
	"net/http"

	"ecommerce/modules/product"
	"ecommerce/packages/httpx"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	svc product.Service
}

func NewProductHandler(svc product.Service) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) Routes(r chi.Router) {
	r.Get("/", h.List)
	r.Get("/{productId}", h.GetByID)
	r.Get("/handle/{handle}", h.GetByHandle)
	r.Post("/", h.Create)
	r.Delete("/{productId}", h.Delete)
	r.Post("/{productId}/variants", h.CreateVariant)
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
