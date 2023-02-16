package hop

import (
	"context"
	"math/rand"
	"strconv"

	"omnifire/api/storage"
	apb "omnifire/proto/api"
	bpb "omnifire/proto/box"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Server) Hop(ctx context.Context, req *apb.HopRequest) (*apb.HopResponse, error) {
	if req.Log.Enabled && req.Log.Section == bpb.Section_SERVICE && req.Log.Host == "api" {
		s.log.Println("logging data: ", req.GetData())
	}
	if req.Error.Enabled && req.Error.Section == bpb.Section_SERVICE && req.Error.Host == "api" {
		s.log.Println("logging data: ", req.GetData())
	}
	if req.Persist && req.Host == "api" {
		_, err := s.db.CreateData(ctx, storage.Data{
			Body: strconv.Itoa(rand.Int()),
		})
		if err != nil {
			s.log.Error(err)
			return nil, err
		}
	}

	if req.Host == "api" {
		return &apb.HopResponse{Response: "api"}, nil
	}

	conn, err := grpc.Dial("box:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		s.log.Error(err)
		return nil, err
	}

	cl := bpb.NewBoxClient(conn)
	res, err := cl.Hop(ctx, &bpb.HopRequest{
		Host:    req.GetHost(),
		Persist: req.GetPersist(),
		Data:    req.GetData(),
		Error: &bpb.Error{
			Enabled: req.GetError().GetEnabled(),
			Host:    req.GetError().GetHost(),
			Section: req.GetError().GetSection(),
		},
		Log: &bpb.Log{
			Enabled: req.GetLog().GetEnabled(),
			Host:    req.GetLog().GetHost(),
			Section: req.GetLog().GetSection(),
		},
	})
	if err != nil {
		s.log.Error(err)
		return nil, err
	}
	return &apb.HopResponse{Response: "api" + res.GetResponse()}, nil
}
