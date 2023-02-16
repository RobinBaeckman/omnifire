package hop

import (
	"context"
	"log"
	"net/http"
	"omnifire/api/storage/postgres"
	apb "omnifire/proto/api"
	"omnifire/util/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	apb.UnimplementedApiServer
	db  *postgres.DB
	log *logger.Log
}

func RegisterServer(s grpc.ServiceRegistrar, db *postgres.DB, log *logger.Log) {
	apb.RegisterApiServer(s, &Server{db: db, log: log})
}

func RegisterGateway(ctx context.Context, conn *grpc.ClientConn) *http.Server {
	gwmux := runtime.NewServeMux()
	if err := apb.RegisterApiHandler(ctx, gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	return &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}
}
