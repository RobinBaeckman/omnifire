// TODO: probably good to move this into brank.as/serviceutil so it can be shared
//
// To instrument add this to the beginning of a function
// log, otl, ctx := otelb.Start(ctx)
// defer otl.Span.End()
package otelb

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otlpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"brank.as/rbac/svcutil/logging"
)

// Initializes an OTLP exporter, and configures the corresponding trace providers.
// usually used in main function
func InitOTELProvider(ctx context.Context, endpoint string, svcAttr ...attribute.KeyValue) func() {
	if endpoint == "" {
		logging.FromContext(ctx).Debug("tracing disabled")
		return func() {}
	}
	logging.FromContext(ctx).Debug("dialing trace collector")
	exp, err := otlpgrpc.New(ctx,
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(endpoint),
		otlpgrpc.WithDialOption(),
	)
	if err != nil {
		log.Fatalln("dialing OTLP agent")
	}

	res, err := resource.New(ctx, resource.WithAttributes(svcAttr...))
	if err != nil {
		log.Fatalln("adding resource attributes")
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithIDGenerator(IDGenerator()),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)

	return func() {
		tp.Shutdown(ctx)
		if err != nil {
			log.Fatalln("stopping provider")
		}
		exp.Shutdown(ctx)
		if err != nil {
			log.Fatalln("stopping exporter")
		}
	}
}
