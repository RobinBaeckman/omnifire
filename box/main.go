package main

import (
	"context"

	"google.golang.org/grpc"

	hsvc "omnifire/box/service/hop"
	"omnifire/box/storage/postgres"
	"omnifire/util/logger"
	"omnifire/util/srv"
)

func main() {
	log := logger.New()

	db := postgres.New()
	defer db.Close()
	db.MigrateUp()

	s := grpc.NewServer()
	hsvc.RegisterServer(s, db, log)
	go func() {
		log.Fatalln(s.Serve(srv.Listen("tcp", ":8080")))
	}()
	log.Infoln("Serving gRPC on 0.0.0.0:8080")

	gwServer := hsvc.RegisterGateway(context.Background(), srv.GRPCClientConn(context.Background()))

	log.Infoln("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
