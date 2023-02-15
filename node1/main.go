package main

import (
	"context"

	"google.golang.org/grpc"

	csvc "omnifire/node1/service/chat"
	"omnifire/node1/storage/postgres"
	"omnifire/util/logger"
	"omnifire/util/srv"
)

func main() {
	log := logger.New()

	db := postgres.New()
	defer db.Close()
	db.MigrateUp()

	s := grpc.NewServer()
	csvc.RegisterServer(s, db, log)
	go func() {
		log.Fatalln(s.Serve(srv.Listen("tcp", ":8080")))
	}()
	log.Infoln("Serving gRPC on 0.0.0.0:8080")

	gwServer := csvc.RegisterGateway(context.Background(), srv.GRPCClientConn(context.Background()))

	log.Infoln("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
