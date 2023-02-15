package chat

import (
	"context"
	"log"
	"net/http"
	"omnifire/node1/storage/postgres"
	cpb "omnifire/proto/chat"
	"omnifire/util/logger"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server struct {
	cpb.UnimplementedChatServer
	db  *postgres.DB
	log *logger.Log
}

func RegisterServer(s grpc.ServiceRegistrar, db *postgres.DB, log *logger.Log) {
	cpb.RegisterChatServer(s, &Server{db: db, log: log})
}

func RegisterGateway(ctx context.Context, conn *grpc.ClientConn) *http.Server {
	gwmux := runtime.NewServeMux()
	if err := cpb.RegisterChatHandler(ctx, gwmux, conn); err != nil {
		log.Fatalln("Failed to register gateway:", err)
	}
	return &http.Server{
		Addr:    ":8090",
		Handler: gwmux,
	}
}
