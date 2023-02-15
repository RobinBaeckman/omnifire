package storage

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrNotFound is returned when the requested resource does not exist.
	ErrNotFound = status.Error(codes.NotFound, "not found")
	// Conflict is returned when trying to create the same resource twice.
	Conflict = status.Error(codes.AlreadyExists, "conflict")
)

type Data struct {
	ID      string    `db:"id"`
	Body    string    `db:"body"`
	Created time.Time `db:"created_at"`
}
