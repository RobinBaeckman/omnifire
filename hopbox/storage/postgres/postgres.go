package postgres

import (
	"context"
	"omnifire/util/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	// https://www.postgresql.org/docs/9.6/errcodes-appendix.html
	pqUnique   = "23505"
	pqNotFound = "42703"
)

type DB struct {
	*sqlx.DB
}

func New(ctx context.Context, cf *viper.Viper) *DB {
	db := postgres.New(ctx, cf)
	return &DB{db}
}

func (db *DB) MigrateUp(ctx context.Context, cf *viper.Viper) {
	postgres.MigrateUp(ctx, cf, db.DB)
}
