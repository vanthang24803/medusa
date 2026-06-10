package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"ecommerce/apps/api/internal/server"
	"ecommerce/packages/config"
	"ecommerce/packages/db"
	"ecommerce/packages/events"
	"ecommerce/packages/logger"
)

func main() {
	cfg := config.Load()
	dev := cfg.Env == "development"
	log := logger.New(dev)
	defer func() { _ = log.Sync() }()

	database := db.MustConnect(db.Config{
		MasterDSN: cfg.DatabaseURL,
		SlaveDSN:  cfg.DatabaseReplicaURL,
	})
	defer database.Close()

	bus, err := events.NewEventBus(events.Config{
		Provider: events.ProviderType(cfg.EventsProvider),
		URL:      cfg.EventsURL,
	}, log)
	if err != nil {
		log.Fatal("failed to initialize event bus", zap.Error(err))
	}
	defer func() { _ = bus.Close() }()

	mods := server.WireModules(database, bus)
	handler := server.New(log, mods)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Info("api listening", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}
