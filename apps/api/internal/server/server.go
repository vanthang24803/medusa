package server

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"ecommerce/apps/api/internal/handler"
	"ecommerce/apps/api/internal/middleware"
	"ecommerce/packages/httpx"
	"ecommerce/packages/types"

	"ecommerce/modules/auth"
	"ecommerce/modules/brand"
	"ecommerce/modules/cart"
	"ecommerce/modules/customer"
	"ecommerce/modules/fulfillment"
	"ecommerce/modules/iam"
	"ecommerce/modules/identity"
	"ecommerce/modules/inventory"
	"ecommerce/modules/notification"
	"ecommerce/modules/order"
	"ecommerce/modules/payment"
	"ecommerce/modules/pricing"
	"ecommerce/modules/product"
	"ecommerce/modules/promotion"
	"ecommerce/modules/region"
)

type Modules struct {
	Auth         auth.Service
	IAM          iam.Service
	Identity     identity.Service
	Customer     customer.Service
	Brand        brand.Service
	Product      product.Service
	Pricing      pricing.Service
	Inventory    inventory.Service
	Cart         cart.Service
	Order        order.Service
	Payment      payment.Service
	Fulfillment  fulfillment.Service
	Promotion    promotion.Service
	Region       region.Service
	Notification notification.Service
}

func New(log *zap.Logger, mods *Modules) http.Handler {
	r := chi.NewRouter()

	middlewares := []func(http.Handler) http.Handler{
		middleware.CORS(middleware.DefaultCORSConfig),
		httpx.RequestIDUUIDMiddleware,
		chimw.RequestID,
		chimw.RealIP,
		httpx.RecordStartTimeMiddleware,
		middleware.Logger(log),
		middleware.Recovery(log),
		middleware.RateLimit(100),
	}

	r.Use(middlewares...)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, r, types.ErrNotFound)
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, r, http.StatusOK, httpx.Response("status", "ok"))
	})

	r.Get("/docs", http.RedirectHandler("/docs/", http.StatusMovedPermanently).ServeHTTP)
	r.Handle("/docs/*", http.StripPrefix("/docs/", noCacheForSpec(http.FileServer(http.Dir("docs")))))

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/auth", handler.NewAuthHandler(mods.Auth, log).Routes)
		r.Route("/customers", handler.NewCustomerHandler(mods.Customer, mods.Auth, log).Routes)
		r.Route("/identity", handler.NewIdentityHandler(mods.Identity, mods.Auth, mods.IAM).Routes)
		r.Route("/products", handler.NewProductHandler(mods.Product, mods.Auth, mods.IAM).Routes)
		r.Route("/brands", handler.NewBrandHandler(mods.Brand, mods.Auth, mods.IAM, log).Routes)
		r.Route("/iam", handler.NewIAMHandler(mods.IAM, mods.Auth).Routes)
	})

	return r
}

// noCacheForSpec disables browser caching for .yml/.yaml files so spec updates
// are always reflected immediately without a hard refresh.
func noCacheForSpec(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, ".yml") || strings.HasSuffix(r.URL.Path, ".yaml") {
			w.Header().Set("Cache-Control", "no-store")
		}
		next.ServeHTTP(w, r)
	})
}
