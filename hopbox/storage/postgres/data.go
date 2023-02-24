package postgres

import (
	"context"
	"errors"
	"omnifire/hopbox/storage"
	hpb "omnifire/proto/hopbox"
	"omnifire/util/logger"
	"omnifire/util/otel"

	"github.com/lib/pq"
	"go.opentelemetry.io/otel/trace"
)

func (db *DB) CreateData(ctx context.Context, sd storage.Hop, req *hpb.HopRequest, srvName string, shouldLog, shouldError bool, ld string) (*storage.Hop, error) {
	log := logger.FromContext(ctx)
	otel.Start(ctx, "",
		otel.WithSpanOpts(trace.WithSpanKind(trace.SpanKindClient)))

	const q = `
INSERT INTO hop (
    body
) VALUES (
	 :body
) RETURNING
    id,created
`
	stmt, err := db.PrepareNamedContext(ctx, q)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&sd, sd); err != nil {
		log.Error(err)
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return nil, storage.Conflict
		}
		return nil, err
	}

	if shouldLog {
		log.Print("log requested: ", ld)
	}
	if shouldError {
		log.Error("error requested")
		return nil, errors.New("request error")
	}
	return &sd, nil
}
