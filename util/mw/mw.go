package mw

import (
	"context"
	"omnifire/util/logger"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func LoggerInterceptor(e *logrus.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = logger.Inject(ctx, e.WithContext(ctx))
		return handler(ctx, req)
	}
}
