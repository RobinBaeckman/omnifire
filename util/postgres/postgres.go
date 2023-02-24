package postgres

import (
	"context"
	"fmt"
	"omnifire/util/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func New(ctx context.Context, cf *viper.Viper) *sqlx.DB {
	log := logger.FromContext(ctx)

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=%s",
		cf.GetString("database.user"), cf.GetString("database.password"), cf.GetString("database.dbname"), cf.GetString("database.host"), cf.GetString("database.sslmode")))
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	log.Info("postgreSQL connected, DB stats:  %+v\n", db.Stats())
	return db
}

func MigrateUp(ctx context.Context, cf *viper.Viper, db *sqlx.DB) {
	log := logger.FromContext(ctx)

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		cf.GetString("database.migrationPath"),
		cf.GetString("database.dbname"), driver)
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
