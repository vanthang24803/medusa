package server

import (
	"net/http"

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

	r.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("docs"))))

	r.Route("/api/v1", func(r chi.Router) {
		routes := map[string]func(chi.Router){
			"/auth":      handler.NewAuthHandler(mods.Auth).Routes,
			"/customers": handler.NewCustomerHandler(mods.Customer, mods.Auth).Routes,
			"/identity":  handler.NewIdentityHandler(mods.Identity, mods.Auth).Routes,
			"/products":  handler.NewProductHandler(mods.Product).Routes,
			"/brands":    handler.NewBrandHandler(mods.Brand).Routes,
		}

		for path, route := range routes {
			r.Route(path, route)
		}
	})

	return r
}
