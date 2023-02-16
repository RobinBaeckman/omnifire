package otelb

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
	"sync"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"

	"brank.as/rbac/svcutil/logging"
)

func IDGenerator() sdktrace.IDGenerator { return &extractIDer{IDGenerator: defaultIDGenerator()} }

var _ sdktrace.IDGenerator = (*extractIDer)(nil)

type extractIDer struct{ sdktrace.IDGenerator }

func (e *extractIDer) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	if cid := logging.CorrelationIDFromContext(ctx); cid != "" {
		tid := trace.TraceID{}
		copy(tid[:], []byte(cid[:8]+cid[12:]))
		return tid, e.IDGenerator.NewSpanID(ctx, tid)
	}
	return e.IDGenerator.NewIDs(ctx)
}

func (e *extractIDer) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	return e.IDGenerator.NewSpanID(ctx, traceID)
}

type randomIDGenerator struct {
	sync.Mutex
	randSource *rand.Rand
}

// NewSpanID returns a non-zero span ID from a randomly-chosen sequence.
func (gen *randomIDGenerator) NewSpanID(ctx context.Context, traceID trace.TraceID) trace.SpanID {
	gen.Lock()
	defer gen.Unlock()
	sid := trace.SpanID{}
	_, _ = gen.randSource.Read(sid[:])
	return sid
}

// NewIDs returns a non-zero trace ID and a non-zero span ID from a
// randomly-chosen sequence.
func (gen *randomIDGenerator) NewIDs(ctx context.Context) (trace.TraceID, trace.SpanID) {
	gen.Lock()
	defer gen.Unlock()
	tid := trace.TraceID{}
	_, _ = gen.randSource.Read(tid[:])
	sid := trace.SpanID{}
	_, _ = gen.randSource.Read(sid[:])
	return tid, sid
}

func defaultIDGenerator() sdktrace.IDGenerator {
	gen := &randomIDGenerator{}
	var rngSeed int64
	_ = binary.Read(crand.Reader, binary.LittleEndian, &rngSeed)
	gen.randSource = rand.New(rand.NewSource(rngSeed))
	return gen
}
