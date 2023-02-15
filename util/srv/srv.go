// package containing helpers related to network and serving
package srv

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Listen(network, address string) net.Listener {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	return lis
}

func GRPCClientConn(ctx context.Context) *grpc.ClientConn {
	conn, err := grpc.DialContext(
		ctx,
		"0.0.0.0:8080",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}
	return conn
}
