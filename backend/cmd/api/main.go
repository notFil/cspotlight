package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/notFil/cspotlight/config"
	"github.com/notFil/cspotlight/internal/logger"
	"github.com/notFil/cspotlight/internal/store"
	"go.uber.org/zap"

	"github.com/notFil/cspotlight/internal/router"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitializeLogger(cfg.Server.Environment)

	port := cfg.Server.Port

	logger.Logger.Info("Starting server on port %d", zap.Int("port", port))

	db := store.ConnectDB(cfg.Store.DSN)

	router := router.SetUpRouter(db, cfg.Server.BaseURL, cfg.Session, cfg.Store.Redis)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Logger.Fatal("failed to run server: %v", zap.Error(err))
	}
}
