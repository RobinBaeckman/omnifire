// To instrument add this to the beginning of a function
// log, otl, ctx := otelb.Start(ctx)
// defer otl.Span.End()
package otel

import (
	"context"
	"fmt"
	"omnifire/util/logger"

	otelpyroscope "github.com/pyroscope-io/otel-profiling-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Initializes an OTLP exporter, and configures the corresponding trace providers.
// usually used in main function
func NewTracer(ctx context.Context, svcName string, otelHost string, pyroHost string, svcAttr ...attribute.KeyValue) func() {
	log := logger.FromContext(ctx)
	switch {
	case otelHost == "":
		log.Debug("tracing disabled")
		return func() {}
	case pyroHost == "":
		log.Debug("profiling disabled")
	}
	log.Info(fmt.Sprintf("dialing otel-collector, %s..", otelHost))
	exp, err := otelpgrpc.New(ctx,
		otelpgrpc.WithInsecure(),
		otelpgrpc.WithEndpoint(otelHost),
		//otlpgrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		log.Fatalln("dialing collector")
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
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	if pyroHost == "" {
		otel.SetTracerProvider(tp)
		fmt.Println("################6")
	} else {
		fmt.Println("################5")
		otel.SetTracerProvider(otelpyroscope.NewTracerProvider(tp,
			otelpyroscope.WithAppName(svcName),
			otelpyroscope.WithPyroscopeURL(pyroHost),
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
