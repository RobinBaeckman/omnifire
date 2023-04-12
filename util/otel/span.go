package otel

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Attr = attribute.Key

const (
	MethodKey = attribute.Key("method")
)

type Option func(*config)

type config struct {
	spanName   string
	tracerName string
	spanOpts   []trace.SpanStartOption
}

// WithTracerName option sets the tracer name, if empty a default value is set by
// opentelemetry
func WithTracerName(sn string) Option { return func(c *config) { c.tracerName = sn } }

func WithSpanOpts(opts ...trace.SpanStartOption) Option {
	return func(c *config) { c.spanOpts = append(c.spanOpts, opts...) }
}

// Start a span which is used to trace a function
func Start(ctx context.Context, name string, opts ...Option) (context.Context, trace.Span) {
	c := &config{spanName: name, spanOpts: []trace.SpanStartOption{}}
	for _, opt := range opts {
		opt(c)
	}

	t := otel.Tracer(c.tracerName)

	var s trace.Span
	if c.spanName != "" {
		ctx, s = t.Start(ctx, c.spanName, c.spanOpts...)
	} else {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			ctx, s = t.Start(ctx, details.Name(), c.spanOpts...)
		} else {
			ctx, s = t.Start(ctx, "error", c.spanOpts...)
		}
	}
	return ctx, s
}
