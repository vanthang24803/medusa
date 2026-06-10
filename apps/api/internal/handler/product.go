package handler

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"ecommerce/modules/product"
	"ecommerce/packages/httpx"
)

// ProductHandler — HTTP layer cho product domain.
type ProductHandler struct {
	svc product.Service
}

func NewProductHandler(svc product.Service) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// Routes mounts the product endpoints onto the router.
func (h *ProductHandler) Routes(r chi.Router) {
	r.Get("/products", h.List)
	r.Post("/products", h.Create)
	r.Get("/products/{id}", h.Get)
	r.Delete("/products/{id}", h.Delete)
	r.Get("/products/{id}/variants", h.ListVariants)
	r.Post("/products/{id}/variants", h.CreateVariant)
	r.Get("/products/handle/{handle}", h.GetByHandle)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	q := product.ListQuery{
		Status:       r.URL.Query().Get("status"),
		CollectionID: r.URL.Query().Get("collection_id"),
		Search:       r.URL.Query().Get("q"),
	}
	q.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	q.PerPage, _ = strconv.Atoi(r.URL.Query().Get("per_page"))

	res, err := h.svc.List(r.Context(), q)
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, res)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.GetByID(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"product": p})
}

func (h *ProductHandler) GetByHandle(w http.ResponseWriter, r *http.Request) {
	p, err := h.svc.GetByHandle(r.Context(), chi.URLParam(r, "handle"))
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"product": p})
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req product.CreateProductReq
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	p, err := h.svc.Create(r.Context(), product.CreateInput{
		Title:        req.Title,
		Subtitle:     req.Subtitle,
		Description:  req.Description,
		Handle:       req.Handle,
		Thumbnail:    req.Thumbnail,
		Status:       product.ProductStatus(req.Status),
		CollectionID: req.CollectionID,
	})
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusCreated, map[string]any{"product": p})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.svc.Delete(r.Context(), chi.URLParam(r, "id")); err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"id": chi.URLParam(r, "id"), "deleted": true})
}

func (h *ProductHandler) ListVariants(w http.ResponseWriter, r *http.Request) {
	variants, err := h.svc.ListVariants(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusOK, map[string]any{"variants": variants})
}

func (h *ProductHandler) CreateVariant(w http.ResponseWriter, r *http.Request) {
	var req product.CreateVariantReq
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, r, err)
		return
	}
	v, err := h.svc.CreateVariant(r.Context(), chi.URLParam(r, "id"), product.CreateVariantInput{
		Title:           req.Title,
		SKU:             req.SKU,
		ManageInventory: req.ManageInventory,
		AllowBackorder:  req.AllowBackorder,
	})
	if err != nil {
		httpx.Error(w, r, err)
		return
	}
	httpx.JSON(w, r, http.StatusCreated, map[string]any{"variant": v})
}
