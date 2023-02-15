package postgres

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
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

func New() *DB {
	db, err := sqlx.Connect("postgres", "user=postgres password=secret dbname=postgres host=postgresql sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("PostgreSQL connected, DB stats:  %+v\n", db.Stats())
	return &DB{db}
}

func (db *DB) MigrateUp() {
	driver, err := postgres.WithInstance(db.DB.DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://storage/migration",
		"postgres", driver)
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
		log.Println("PostgreSQL migration is dirty")
	}
	log.Println("PostgreSQL migration version: ", v)
}
