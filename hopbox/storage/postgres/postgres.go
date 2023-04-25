package postgres

import (
	"context"
	"omnifire/util/config"
	"omnifire/util/db"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	// https://www.postgresql.org/docs/9.6/errcodes-appendix.html
	pqUnique   = "23505"
	pqNotFound = "42703"
)

type DB struct {
	*sqlx.DB
}

func New(ctx context.Context, cf *config.Config) *DB {
	db := db.New(ctx, cf)
	return &DB{db}
}

func (st *DB) MigrateUp(ctx context.Context, cf *config.Config) {
	db.MigrateUp(ctx, cf, st.DB)
}
