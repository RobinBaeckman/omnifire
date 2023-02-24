package logger

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type (
	loggerKey struct{}
)

func New(ctx context.Context, cf *viper.Viper) (*logrus.Entry, context.Context) {
	e := logrus.WithContext(ctx)
	e = SetDefaultFields(cf, e)
	ctx = context.WithValue(ctx, loggerKey{}, e)
	return e, ctx
}

func SetDefaultFields(cf *viper.Viper, e *logrus.Entry) *logrus.Entry {
	e.Logger.SetFormatter(&logrus.JSONFormatter{})
	e.Logger.SetReportCaller(true)

	e = e.
		WithField("service", cf.GetString("server.name")).
		WithField("version", "todo")
	return e
}

func Inject(ctx context.Context, e *logrus.Entry, cf *viper.Viper) context.Context {
	e = SetDefaultFields(cf, e)

	ctx = context.WithValue(ctx, loggerKey{}, e)
	return ctx
}

func FromContext(ctx context.Context) *logrus.Entry {
	v := ctx.Value(loggerKey{})
	logr, ok := v.(*logrus.Entry)
	if !ok {
		log.Fatal("missing logger from context")
	}
	return logr
}
