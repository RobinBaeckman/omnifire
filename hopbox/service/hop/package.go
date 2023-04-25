package hop

import (
	"context"
	"fmt"
	"net/http"
	"omnifire/hopbox/storage/postgres"
	hpb "omnifire/proto/hopbox"
	cf "omnifire/util/config"
	"omnifire/util/logger"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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

func NewConfig(ctx context.Context, cf *cf.Config) config {
	log := logger.FromContext(ctx)
	sn := cf.Server.Name
	hopNo, err := strconv.Atoi(sn[len(sn)-1:])
	if err != nil {
		log.Fatal(err)
	}
	return config{
		hopNo:   hopNo,
		srvName: cf.Server.Name,
		nextHop: cf.NextHop.Host != "",
	}
}

func RegisterServer(s grpc.ServiceRegistrar, db *postgres.DB, cl hpb.HopboxClient, cf config) {
	hpb.RegisterHopboxServer(s, &Server{db: db, cf: cf, cl: cl})
}

func RegisterGateway(ctx context.Context, conn *grpc.ClientConn, addr string) *http.Server {
	log := logger.FromContext(ctx)
	gwmux := runtime.NewServeMux()
	fmt.Println("################19")
	if err := hpb.RegisterHopboxHandler(ctx, gwmux, conn); err != nil {
		log.Fatalln("failed to register gateway:", err)
	}
	fmt.Println("################20")
	return &http.Server{
		Addr:    addr,
		Handler: gwmux,
	}
}
