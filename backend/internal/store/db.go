package store

import (
	"cspotlight/internal/logger"

	gormsessions "github.com/gin-contrib/sessions/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(dsn string) Database {
	return Database{
		DB: openDB(dsn),
	}
}

func openDB(dsn string) *gorm.DB {
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

func (db *Database) CreateSessionStore(secretKey []byte) gormsessions.Store {
	session := gormsessions.NewStore(db.DB, true, secretKey)
	if session == nil {
		logger.Logger.Panic("failed to initialize session store")
	}
	return session
}
