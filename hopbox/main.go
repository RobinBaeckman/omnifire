package main

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	hsvc "omnifire/hopbox/service/hop"
	"omnifire/hopbox/storage/postgres"
	hpb "omnifire/proto/hopbox"
	"omnifire/util/logger"
	"omnifire/util/mw"
	"omnifire/util/srv"
	"omnifire/util/viper"

	"omnifire/util/otel"

	vpr "github.com/spf13/viper"

	"github.com/pyroscope-io/client/pyroscope"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
)

const (
	dockerEnv = "dev"
)

func main() {
	cf := viper.New()
	log, ctx := logger.New(context.Background(), cf)

	pf, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: cf.GetString("server.name"),
		ServerAddress:   cf.GetString("profile.host"),
		Logger:          pyroscope.StandardLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := pf.Stop(); err != nil {
			log.Error(err)
		}
	}()

	attr := []attribute.KeyValue{
		semconv.ServiceNameKey.String(cf.GetString("server.name")),
		semconv.DeploymentEnvironmentKey.String(cf.GetString("runtime.env")),
		semconv.ServiceVersionKey.String("todo"),
		attribute.String("app", cf.GetString("server.name")),
	}
	if cf.GetString("runtime.env") == dockerEnv {
		attr = append(attr, attribute.String("container", cf.GetString("server.name")))
	}
	shutdown := otel.NewTracerWithProfiler(
		ctx,
		cf.GetString("server.name"),
		cf.GetString("trace.collectorHost"),
		cf.GetString("profile.host"),
		attr...,
	)
	defer shutdown()

	db := postgres.New(ctx, cf)
	defer db.Close()
	db.MigrateUp(ctx, cf)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			otelgrpc.UnaryServerInterceptor(),
			mw.LoggerInterceptor(cf, log),
		),
	)

	conns := newConns(ctx, cf)
	hcl := hpb.NewHopboxClient(conns.HopConn)

	hsvc.RegisterServer(s, db, hcl, hsvc.NewConfig(ctx, cf))
	go func() {
		log.Fatalln(s.Serve(srv.Listen("tcp", ":"+cf.GetString("server.grpcPort"))))
	}()
	log.Infoln("serving gRPC on 0.0.0.0:" + cf.GetString("server.grpcPort"))

	gwServer := hsvc.RegisterGateway(
		ctx,
		srv.GRPCClientConn(
			ctx,
			":"+cf.GetString("server.grpcPort"),
		),
		":"+cf.GetString("server.httpPort"),
	)
	log.Infoln("serving gRPC-Gateway on http://0.0.0.0:" + cf.GetString("server.httpPort"))
	log.Fatalln(gwServer.ListenAndServe())
}

type conns struct {
	HopConn *grpc.ClientConn
}

func newConns(ctx context.Context, cf *vpr.Viper) *conns {
	log := logger.FromContext(ctx)
	conn := &grpc.ClientConn{}
	nextHop := cf.GetString("nexthop.addr")
	if nextHop != "" {
		var err error
		log.Info("dialing hopbox: ", nextHop)
		conn, err = grpc.Dial(nextHop,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		)
		if err != nil {
			log.Fatal(err)
		}
	}
	return &conns{HopConn: conn}
}
