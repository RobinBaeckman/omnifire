package mw

import (
	"context"
	"omnifire/util/logger"
	"omnifire/util/otel"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func LoggerInterceptor(cf *viper.Viper, e *logrus.Entry) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		e = logrus.WithContext(ctx)
		e.Logger.AddHook(&otel.LogHook{})
		ctx = logger.Inject(ctx, e, cf)
		return handler(ctx, req)
	}
}
