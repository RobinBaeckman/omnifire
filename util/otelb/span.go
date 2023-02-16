package otelb

import (
	"context"
	"runtime"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"brank.as/rbac/svcutil/logging"
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
func Start(ctx context.Context, name string, opts ...Option) (*logrus.Entry, trace.Span, context.Context) {
	c := &config{spanName: name, spanOpts: []trace.SpanStartOption{}}
	for _, opt := range opts {
		opt(c)
	}

	// todo: leaving tracename as empty by default for now, it will be set to its default
	// we might want to set this to package name or something, not really sure
	// what the usecase is for this yet
	var t trace.Tracer
	if c.tracerName != "" {
		t = otel.Tracer(c.tracerName)
	} else {
		t = otel.Tracer("")
	}

	var s trace.Span
	if c.spanName != "" {
		ctx, s = t.Start(ctx, c.spanName, c.spanOpts...)
	} else {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		// todo: might be a better way of setting the span name to the function name
		// maybe we just want the function name without the full path
		if ok && details != nil {
			ctx, s = t.Start(ctx, details.Name(), c.spanOpts...)
		} else {
			ctx, s = t.Start(ctx, "error", c.spanOpts...)
		}
	}
	log := logging.FromContext(ctx).WithContext(ctx)
	log = log.WithField("span", c.spanName)
	if s.IsRecording() {
		log = log.WithField("traceid", s.SpanContext().TraceID())
	}
	return log, s, ctx
}
