package db

import (
	"context"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to database")
	}

	if err := ping(db); err != nil {
		return nil, errors.Wrap(err, "failed to ping db")
	}

	if err := bootstrapDatabase(dsn); err != nil {
		return nil, errors.Wrap(err, "failed to bootstrap database")
	}

	return db, nil
}

func ping(db *sqlx.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}

func bootstrapDatabase(dsn string) error {
	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return errors.Wrap(err, "failed to create a migration instance")
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "failed to migrate database")
	}

	return nil
}
