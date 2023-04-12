package hop

import (
	"context"
	"errors"
	"time"

	"omnifire/hopbox/storage"
	hpb "omnifire/proto/hopbox"

	"omnifire/util/logger"
	"omnifire/util/otel"

	"go.opentelemetry.io/otel/trace"
)

func (s *Server) Hop(ctx context.Context, req *hpb.HopRequest) (*hpb.HopResponse, error) {
	ctx, span := otel.Start(ctx, "",
		otel.WithSpanOpts(trace.WithSpanKind(trace.SpanKindInternal)))
	defer span.End()
	log := logger.FromContext(ctx)

	sp, spd := s.shouldPersist(req)
	if sp {
		sl, sld := s.shouldLog(req, hpb.Section_STORAGE.String())
		_, err := s.db.CreateData(ctx, storage.Hop{
			Body: spd,
		}, req,
			s.cf.srvName,
			sl,
			s.shouldError(req, hpb.Section_STORAGE.String()),
			sld,
		)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	if s.shouldWork(req) {
		time.Sleep(time.Duration(req.Work.DurationSec) * time.Second)
	}
	sl, sld := s.shouldLog(req, hpb.Section_SERVICE.String())
	if sl {
		log.Print("log requested: ", sld)
	}
	if s.shouldError(req, hpb.Section_SERVICE.String()) {
		log.Error("error requested")
		return nil, errors.New("error requested")
	}

	if !s.cf.nextHop {
		return &hpb.HopResponse{Response: s.cf.srvName}, nil
	}

	res, err := s.cl.Hop(ctx, req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &hpb.HopResponse{Response: s.cf.srvName + "->" + res.GetResponse()}, nil
}

func (s *Server) shouldLog(req *hpb.HopRequest, section string) (bool, string) {
	for _, l := range req.GetLog() {
		if l.GetHop() == int32(s.cf.hopNo) && l.GetSection().String() == section {
			return true, l.GetData()
		}
	}
	return false, ""
}

func (s *Server) shouldError(req *hpb.HopRequest, section string) bool {
	if req.GetError().GetHop() == int32(s.cf.hopNo) && req.GetError().GetSection().String() == section {
		return true
	}
	return false
}

func (s *Server) shouldPersist(req *hpb.HopRequest) (bool, string) {
	for _, l := range req.GetPersist() {
		if l.GetHop() == int32(s.cf.hopNo) {
			return true, l.GetData()
		}
	}
	return false, ""
}

func (s *Server) shouldWork(req *hpb.HopRequest) bool {
	if req.GetWork().GetHop() == int32(s.cf.hopNo) {
		return true
	}
	return false
}
