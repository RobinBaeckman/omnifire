package chat

import (
	"context"

	cpb "omnifire/proto/chat"
)

func (s *Server) Talk(ctx context.Context, req *cpb.TalkRequest) (*cpb.TalkResponse, error) {
	s.log.Infoln("method: Talk")
	s.log.Infoln("talk: ", req.GetMsg())
	return &cpb.TalkResponse{Response: "hello"}, nil
}
