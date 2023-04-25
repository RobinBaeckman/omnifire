package logger

import (
	"context"
	"log"
	"omnifire/util/config"

	"github.com/sirupsen/logrus"
)

type (
	loggerKey struct{}
)

func New(ctx context.Context, cf *config.Config) (*logrus.Entry, context.Context) {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetReportCaller(true)

	e := log.
		WithField("app", cf.Server.Name).
		WithField("version", "todo")
	ctx = context.WithValue(ctx, loggerKey{}, e)
	return e, ctx
}

func Inject(ctx context.Context, e *logrus.Entry) context.Context {
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
