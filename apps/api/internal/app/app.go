package app

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"ecommerce/apps/api/internal/server"
	"ecommerce/packages/config"
	"ecommerce/packages/db"
	"ecommerce/packages/events"
	"ecommerce/packages/logger"
	"ecommerce/packages/upload"

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

func New() *fx.App {
	return fx.New(
		fx.Provide(config.Load),
		fx.Provide(newLogger),
		fx.Provide(newDB),
		fx.Provide(newEventBus),
		fx.Provide(wireModules),
		fx.Provide(server.New),
		fx.Invoke(startServer),
	)
}

func newLogger(cfg *config.Config, lc fx.Lifecycle) *zap.Logger {
	log := logger.New(cfg.Env == "development")
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			_ = log.Sync()
			return nil
		},
	})
	return log
}

func newDB(cfg *config.Config, lc fx.Lifecycle) *db.DB {
	database := db.MustConnect(db.Config{
		MasterDSN: cfg.DatabaseURL,
		SlaveDSN:  cfg.DatabaseReplicaURL,
	})
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			database.Close()
			return nil
		},
	})
	return database
}

func newEventBus(cfg *config.Config, log *zap.Logger, lc fx.Lifecycle) events.EventBus {
	bus, err := events.NewEventBus(events.Config{
		Provider: events.ProviderType(cfg.EventsProvider),
		URL:      cfg.EventsURL,
	}, log)
	if err != nil {
		log.Fatal("failed to initialize event bus", zap.Error(err))
	}
	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			bus.Close()
			return nil
		},
	})
	return bus
}

func wireModules(database *db.DB, bus events.EventBus, cfg *config.Config, log *zap.Logger) *server.Modules {
	uploader, err := upload.NewUploader(upload.Config{
		Provider:        upload.ProviderType(cfg.UploadProvider),
		Endpoint:        cfg.UploadEndpoint,
		Region:          cfg.UploadRegion,
		AccessKeyID:     cfg.UploadAccessKey,
		SecretAccessKey: cfg.UploadSecretKey,
		UseSSL:          cfg.UploadSSL,
		PublicURL:       cfg.UploadPublicURL,
	})
	if err != nil {
		log.Warn("upload provider unavailable; avatar uploads disabled", zap.Error(err))
		uploader = upload.NewNopUploader()
	}

	custRepo := customer.NewRepository(database)
	return &server.Modules{
		Auth:         auth.NewService(auth.NewRepository(database), custRepo, bus, cfg.JWTSecret),
		Identity:     identity.NewService(identity.NewRepository(database), bus, uploader),
		Customer:     customer.NewService(custRepo, bus, uploader),
		Brand:        brand.NewService(brand.NewRepository(database), bus),
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

func startServer(lc fx.Lifecycle, handler http.Handler, cfg *config.Config, log *zap.Logger) {
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("api listening", zap.String("port", cfg.Port))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal("server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("shutting down")
			return srv.Shutdown(ctx)
		},
	})
}
