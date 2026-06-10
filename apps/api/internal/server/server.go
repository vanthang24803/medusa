package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"ecommerce/apps/api/internal/handler"
	"ecommerce/apps/api/internal/middleware"
	"ecommerce/packages/db"
	"ecommerce/packages/events"
	"ecommerce/packages/httpx"
	"ecommerce/packages/types"

	"ecommerce/modules/auth"
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

// Modules groups all wired services — to inject into handlers or tests.
type Modules struct {
	Auth         auth.Service
	Identity     identity.Service
	Customer     customer.Service
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

// WireModules initializes repositories + services for all 13 modules.
func WireModules(database *db.DB, bus events.EventBus) *Modules {
	return &Modules{
		Auth:         auth.NewService(auth.NewRepository(database), bus),
		Identity:     identity.NewService(identity.NewRepository(database), bus),
		Customer:     customer.NewService(customer.NewRepository(database), bus),
		Product:      product.NewService(product.NewRepository(database), bus),
		Pricing:      pricing.NewService(pricing.NewRepository(database), bus),
		Inventory:    inventory.NewService(inventory.NewRepository(database), bus),
		Cart:         cart.NewService(cart.NewRepository(database), bus),
		Order:        order.NewService(order.NewRepository(database), bus),
		Payment:      payment.NewService(payment.NewRepository(database), bus),
		Fulfillment:  fulfillment.NewService(fulfillment.NewRepository(database), bus),
		Promotion:    promotion.NewService(promotion.NewRepository(database), bus),
		Region:       region.NewService(region.NewRepository(database), bus),
		Notification: notification.NewService(notification.NewRepository(database), bus),
	}
}

// New creates an HTTP handler with all routes mounted.
func New(log *zap.Logger, mods *Modules) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(httpx.RequestIDUUIDMiddleware)
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(httpx.RecordStartTimeMiddleware)
	r.Use(middleware.Logger(log))
	r.Use(middleware.Recovery(log))

	// Custom 404 handler for unmatched routes
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.Error(w, r, types.ErrNotFound)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.JSON(w, r, http.StatusOK, map[string]string{"status": "ok"})
	})

	// API v1
	r.Route("/api/v1", func(r chi.Router) {
		// Product domain — fully implemented
		handler.NewProductHandler(mods.Product).Routes(r)

		// Other domains: handlers will be mounted similarly when implemented.
		// handler.NewCartHandler(mods.Cart).Routes(r)
		// handler.NewOrderHandler(mods.Order).Routes(r)
		// ...
	})

	return r
}
