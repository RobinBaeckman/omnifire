package db

import (
	"context"
	"fmt"
	"omnifire/util/config"
	"omnifire/util/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

func New(ctx context.Context, cf *config.Config) *sqlx.DB {
	log := logger.FromContext(ctx)

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s",
		cf.Database.User, cf.Database.Password, cf.Database.DBName, cf.Database.Host, cf.Database.SslMode))
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	log.Info("postgreSQL connected, DB stats:  %+v\n", db.Stats())
	return db
}

func MigrateUp(ctx context.Context, cf *config.Config, db *sqlx.DB) {
	log := logger.FromContext(ctx)

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		cf.Database.MigrationPath,
		cf.Database.DBName, driver)
	if err != nil {
		log.Fatalln(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalln(err)
	}
	v, d, err := m.Version()
	if err != nil {
		log.Fatalln(err)
	}
	if d {
		log.Info("PostgreSQL migration is dirty")
	}
	log.Info("PostgreSQL migration version: ", v)
}
