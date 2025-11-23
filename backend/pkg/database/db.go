package database

import (
	"github.com/notFil/cspotlight/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg configs.StoreConfig) *gorm.DB {
	dsn := cfg.DatabaseURL
	db, err := gorm.Open(postgres.New(
		postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{},
	)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
