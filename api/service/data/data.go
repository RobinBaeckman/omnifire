package data

import (
	"context"
	"log"
	"net/http"
	"omnifire/api/storage/postgres"
	dpb "omnifire/proto/data"
	"omnifire/util/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	dpb.UnimplementedDataServer
	db  *postgres.DB
	log *logger.Log
}

func RegisterServer(s grpc.ServiceRegistrar, db *postgres.DB, log *logger.Log) {
	dpb.RegisterDataServer(s, &Server{db: db, log: log})
}

func RegisterGateway(ctx context.Context, conn *grpc.ClientConn) *http.Server {
	gwmux := runtime.NewServeMux()
	if err := dpb.RegisterDataHandler(ctx, gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	return &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}
}
