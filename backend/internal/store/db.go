package store

import (
	"github.com/notFil/cspotlight/internal/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(dsn string) *gorm.DB {
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
