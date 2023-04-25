package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	hsvc "omnifire/hopbox/service/hop"
	"omnifire/hopbox/storage/postgres"
	hpb "omnifire/proto/hopbox"
	"omnifire/util/config"
	"omnifire/util/db"
	"omnifire/util/logger"
	"omnifire/util/mw"
	"omnifire/util/prof"
	"omnifire/util/srv"

	"omnifire/util/otel"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
)

func main() {
	// setup config
	cf := config.New()

	// setup logger
	log, ctx := logger.New(context.Background(), cf)
	log.Logger.AddHook(&otel.LogHook{})

	// setup profiling
	pf := prof.Start(ctx, cf)
	defer func() {
		if err := pf.Stop(); err != nil {
			log.Error(err)
		}
	}()

	// setup tracing
	shutdown := otel.NewTracer(ctx, cf)
	defer shutdown()

	// setup db
	st := postgres.New(ctx, cf)
	defer st.Close()

	// migrate db
	db.MigrateUp(ctx, cf, st.DB)

	// start grpc server
	errChan := make(chan error)
	go func() {
		errChan <- startGRPCServer(ctx, cf, st)
	}()

	// start http server
	go func() {
		errChan <- startGRPCGatewayServer(ctx, cf)
	}()

	for err := range errChan {
		if err != nil {
			log.Fatalf("Error occurred: %v", err)
		}
	}
}

type conns struct {
	HopConn *grpc.ClientConn
}

func newConns(ctx context.Context, cf *config.Config) *conns {
	log := logger.FromContext(ctx)
	conn := &grpc.ClientConn{}
	nextHop := cf.NextHop.Host

	var opt []grpc.DialOption
	opt = append(opt,
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	opt = append(opt, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if nextHop != "" {
		var err error
		log.Info("dialing hopbox: ", nextHop)
		conn, err = grpc.Dial(nextHop,
			opt...,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
	return &conns{HopConn: conn}
}

func startGRPCServer(ctx context.Context, cf *config.Config, st *postgres.DB) error {
	log := logger.FromContext(ctx)

	var opt []grpc.ServerOption
	opt = append(opt, grpc.ChainUnaryInterceptor(
		otelgrpc.UnaryServerInterceptor(),
		mw.LoggerInterceptor(log),
	))
	s := grpc.NewServer(
		opt...,
	)

	conns := newConns(ctx, cf)
	hcl := hpb.NewHopboxClient(conns.HopConn)

	hsvc.RegisterServer(s, st, hcl, hsvc.NewConfig(ctx, cf))
	log.Infoln("serving gRPC on 0.0.0.0:" + cf.Server.GrpcPort)
	return s.Serve(srv.Listen("tcp", ":"+cf.Server.GrpcPort))
}

func startGRPCGatewayServer(ctx context.Context, cf *config.Config) error {
	log := logger.FromContext(ctx)

	fmt.Println("################17")
	conn, err := grpc.DialContext(
		ctx,
		cf.Server.Name+":"+cf.Server.GrpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("################18")

	gwServer := hsvc.RegisterGateway(
		ctx,
		conn,
		":"+cf.Server.HttpPort,
	)
	log.Infoln("serving gRPC-Gateway on http://0.0.0.0:" + cf.Server.HttpPort)
	return gwServer.ListenAndServe()
}
