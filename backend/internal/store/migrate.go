package store

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

func RunMigrations(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
