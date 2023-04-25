// To instrument add this to the beginning of a function
// log, otl, ctx := otelb.Start(ctx)
// defer otl.Span.End()
package otel

import (
	"context"
	"fmt"
	cf "omnifire/util/config"
	"omnifire/util/logger"

	otelpyroscope "github.com/pyroscope-io/otel-profiling-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
)

const dockerEnv = "dev"

// Initializes an OTLP exporter, and configures the corresponding trace providers.
// usually used in main function
func NewTracer(ctx context.Context, cf *cf.Config) func() {
	attr := []attribute.KeyValue{
		semconv.ServiceNameKey.String(cf.Server.Name),
		semconv.DeploymentEnvironmentKey.String(cf.Runtime.Env),
		//semconv.ServiceVersionKey.String("todo"),
		attribute.String("app", cf.Server.Name),
	}
	if cf.Runtime.Env == dockerEnv {
		attr = append(attr, attribute.String("container", cf.Server.Name))
	}

	log := logger.FromContext(ctx)
	switch {
	case cf.Trace.CollectorHost == "":
		log.Debug("tracing disabled")
		return func() {}
	case cf.Profile.Host == "":
		log.Debug("profiling disabled")
	}
	log.Info(fmt.Sprintf("dialing otel-collector, %s..", cf.Trace.CollectorHost))
	exp, err := otelpgrpc.New(ctx,
		otelpgrpc.WithInsecure(),
		otelpgrpc.WithEndpoint(cf.Trace.CollectorHost),
		//otlpgrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		log.Fatalln("dialing collector")
	}

	res, err := resource.New(ctx, resource.WithAttributes(attr...))
	if err != nil {
		log.Fatalln("adding resource attributes")
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	if cf.Profile.Host == "" {
		otel.SetTracerProvider(tp)
	} else {
		otel.SetTracerProvider(otelpyroscope.NewTracerProvider(tp,
			otelpyroscope.WithAppName(cf.Server.Name),
			otelpyroscope.WithPyroscopeURL(cf.Profile.Host),
			otelpyroscope.WithRootSpanOnly(true),
			otelpyroscope.WithAddSpanName(true),
			otelpyroscope.WithProfileURL(true),
			otelpyroscope.WithProfileBaselineURL(true),
			//otelpyroscope.WithProfileBaselineLabels(map[string]string{"robin": "rovin"}),
		))
	}

	log.Info("tracing running")

	return func() {
		tp.Shutdown(ctx)
		if err != nil {
			log.Error("stopping provider")
		}
		exp.Shutdown(ctx)
		if err != nil {
			log.Error("stopping exporter")
		}
	}
}
