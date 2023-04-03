package hop

import (
	"context"
	"net/http"
	"omnifire/hopbox/storage/postgres"
	hpb "omnifire/proto/hopbox"
	"omnifire/util/logger"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type Server struct {
	hpb.UnimplementedHopboxServer
	db *postgres.DB
	cf config
	cl hpb.HopboxClient
}

type config struct {
	srvName string
	nextHop bool
	hopNo   int
}

func NewConfig(ctx context.Context, cf *viper.Viper) config {
	log := logger.FromContext(ctx)
	sn := cf.GetString("server.name")
	hopNo, err := strconv.Atoi(sn[len(sn)-1:])
	if err != nil {
		log.Fatal(err)
	}
	return config{
		hopNo:   hopNo,
		srvName: cf.GetString("server.name"),
		nextHop: cf.GetString("nexthop.addr") != "",
	}
}

func RegisterServer(s grpc.ServiceRegistrar, db *postgres.DB, cl hpb.HopboxClient, cf config) {
	hpb.RegisterHopboxServer(s, &Server{db: db, cf: cf, cl: cl})
}

func RegisterGateway(ctx context.Context, conn *grpc.ClientConn, addr string) *http.Server {
	log := logger.FromContext(ctx)
	gwmux := runtime.NewServeMux()
	if err := hpb.RegisterHopboxHandler(ctx, gwmux, conn); err != nil {
		log.Fatalln("failed to register gateway:", err)
	}
	return &http.Server{
		Addr:    addr,
		Handler: gwmux,
	}
}
