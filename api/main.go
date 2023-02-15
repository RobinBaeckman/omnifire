package main

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"

	dsvc "omnifire/api/service/data"
	"omnifire/api/storage/postgres"
	"omnifire/util/logger"
	"omnifire/util/srv"

	kc "github.com/ricardo-ch/go-kafka-connect/lib/connectors"
)

func createConnector() {
	cl := kc.NewClient("http://mykafka-cp-kafka-connect.default.svc.cluster.local:8083")
	resp, err := cl.CreateConnector(
		kc.CreateConnectorRequest{
			ConnectorRequest: kc.ConnectorRequest{Name: "inventory-connector"},
			Config: map[string]interface{}{
				"connector.class":                          "io.debezium.connector.postgresql.PostgresConnector",
				"tasks.max":                                "1",
				"database.hostname":                        "postgres",
				"database.port":                            "5432",
				"database.user":                            "postgres",
				"database.password":                        "postgres",
				"database.dbname":                          "inventory",
				"database.server.name":                     "dbserver1",
				"database.whitelist":                       "inventory",
				"database.history.kafka.bootstrap.servers": "kafka:9092",
				"database.history.kafka.topic":             "schema-changes.inventory",
			},
		},
		true,
	)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("################5", resp)
}

func main() {
	createConnector()

	log := logger.New()

	db := postgres.New()
	defer db.Close()
	db.MigrateUp()

	s := grpc.NewServer()
	dsvc.RegisterServer(s, db, log)
	go func() {
		log.Fatalln(s.Serve(srv.Listen("tcp", ":8080")))
	}()
	log.Infoln("Serving gRPC on 0.0.0.0:8080")

	gwServer := dsvc.RegisterGateway(context.Background(), srv.GRPCClientConn(context.Background()))

	log.Infoln("Serving gRPC-Gateway on http://0.0.0.0:8090")
	log.Fatalln(gwServer.ListenAndServe())
}
