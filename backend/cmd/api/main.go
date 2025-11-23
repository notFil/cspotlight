package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/notFil/cspotlight/pkg/database"

	"github.com/notFil/cspotlight/configs"
	"github.com/notFil/cspotlight/internal/router"
)

func main() {
	cfg := configs.LoadConfig()

	port := cfg.Server.Port

	fmt.Printf("Starting server on port %d\n", port)

	db := database.ConnectDB(cfg.Store)

	router := router.SetUpRouter(db, cfg.JWTAuth)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
