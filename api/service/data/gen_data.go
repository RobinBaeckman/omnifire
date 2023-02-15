package data

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"

	"omnifire/api/storage"
	dpb "omnifire/proto/data"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GenData(ctx context.Context, req *emptypb.Empty) (*dpb.GenDataResponse, error) {
	s.log.Infoln("method: GenData")
	d, err := s.db.CreateData(ctx, storage.Data{
		Body: strconv.Itoa(rand.Int()),
	})
	if err != nil {
		return nil, err
	}

	res, err := http.Get("http://node1/talk?msg=testing")
	if err != nil {
		return nil, err
	}
	s.log.Infoln(res)
	return &dpb.GenDataResponse{DataGenerated: d.Body}, nil
}
