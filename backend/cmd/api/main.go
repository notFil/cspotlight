package main

import (
	"fmt"
	"net/http"
	"time"

	"cspotlight/config"
	"cspotlight/internal/logger"
	"cspotlight/internal/store"

	"go.uber.org/zap"

	"cspotlight/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	isProduction := cfg.IsProduction()

	logger.InitializeLogger(isProduction)

	port := cfg.Server.Port

	logger.Logger.Info("Starting server on port %d", zap.Int("port", port))

	db := store.NewDatabase(cfg.Store.DSN)

	router := router.SetUpRouter(db, cfg.Server.BaseURL, cfg.Server.StaticPath, cfg.Session, cfg.CORS, isProduction)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Logger.Fatal("failed to run server: %v", zap.Error(err))
	}
}
