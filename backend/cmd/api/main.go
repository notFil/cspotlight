package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/notFil/cspotlight/config"
	"github.com/notFil/cspotlight/pkg/database"
	"github.com/notFil/cspotlight/pkg/logger"
	"go.uber.org/zap"

	"github.com/notFil/cspotlight/internal/router"
)

func main() {
	cfg := config.LoadConfig()
	logger.InitializeLogger(cfg.Server.Environment)

	port := cfg.Server.Port

	logger.Logger.Info("Starting server on port %d", zap.Int("port", port))

	db := database.ConnectDB(cfg.Store)

	router := router.SetUpRouter(db, cfg.JWTAuth)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Logger.Fatal("failed to run server: %v", zap.Error(err))
	}
}
