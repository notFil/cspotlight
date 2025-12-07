package database

import (
	"github.com/notFil/cspotlight/config"
	"github.com/notFil/cspotlight/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg config.StoreConfig) *gorm.DB {
	dsn := cfg.DatabaseURL
	db, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{},
	)
	if err != nil {
		logger.Logger.Panic("failed to connect database")
	}
	return db
}
