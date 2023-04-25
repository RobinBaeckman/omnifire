package prof

import (
	"context"
	"omnifire/util/config"
	"omnifire/util/logger"

	"github.com/pyroscope-io/client/pyroscope"
)

func Start(ctx context.Context, cf *config.Config) *pyroscope.Profiler {
	log := logger.FromContext(ctx)
	pf, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: cf.Server.Name,
		ServerAddress:   cf.Profile.Host,
		//Logger:          pyroscope.StandardLogger,
	})
	if err != nil {
		log.Fatal(err)
	}
	return pf
}
