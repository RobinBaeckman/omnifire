package mw

import (
	"context"
	"omnifire/util/logger"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func LoggerInterceptor(cf *viper.Viper, e *logrus.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		ctx = logger.Inject(ctx, e.WithContext(ctx), cf)
		return handler(ctx, req)
	}
}
