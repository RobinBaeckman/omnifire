package postgres

import (
	"context"
	"log"
	"omnifire/box/storage"

	"github.com/lib/pq"
)

func (s *DB) CreateData(ctx context.Context, d storage.Data) (*storage.Data, error) {
	// todo fix logging

	const q = `
INSERT INTO box (
    body
) VALUES (
	 :body
) RETURNING
    id,created
`
	stmt, err := s.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&d, d); err != nil {
		log.Println(err)
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, err
	}
	return &d, nil
}
