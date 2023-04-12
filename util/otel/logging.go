package otel

import (
	"errors"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"go.opentelemetry.io/otel/trace"
)

type LogHook struct{}

func (h *LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire hook for logging tracing events
func (h *LogHook) Fire(e *logrus.Entry) error {
	s := trace.SpanFromContext(e.Context)
	if s.SpanContext().HasTraceID() {
		e.Data["trace-id"] = s.SpanContext().TraceID()
		e.Data["span-id"] = s.SpanContext().SpanID()
		e.Data["profile_id"] = s.SpanContext().SpanID()
	}
	as := []attribute.KeyValue{}
	if e.HasCaller() {
		as = []attribute.KeyValue{
			semconv.CodeFilepathKey.String(e.Caller.File),
			semconv.CodeLineNumberKey.Int(e.Caller.Line),
		}
	}
	switch e.Level {
	case logrus.ErrorLevel:
		s.RecordError(errors.New(e.Message), trace.WithAttributes(as...))
		s.SetStatus(codes.Error, "internal error")
	default:
		as := append(as, semconv.MessageIDKey.String(e.Message))
		s.AddEvent(e.Level.String(), trace.WithAttributes(as...))
	}
	return nil
}
