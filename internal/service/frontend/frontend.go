package frontend

import (
	"context"
	"embed"

	"github.com/dagu-dev/dagu/internal/config"
	"github.com/dagu-dev/dagu/internal/engine"
	"github.com/dagu-dev/dagu/internal/logger"
	"github.com/dagu-dev/dagu/internal/service/frontend/dag"
	"github.com/dagu-dev/dagu/internal/service/frontend/server"
	"go.uber.org/fx"
)

var (
	//go:embed templates/* assets/*
	assetsFS embed.FS
)

var Module = fx.Options(fx.Provide(NewServer))

type Params struct {
	fx.In

	Config *config.Config
	Logger logger.Logger
	Engine engine.Engine
}

func LifetimeHooks(lc fx.Lifecycle, srv *server.Server) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) (err error) {
				return srv.Serve(ctx)
			},
			OnStop: func(_ context.Context) error {
				srv.Shutdown()
				return nil
			},
		},
	)
}

func NewServer(params Params) *server.Server {

	var hs []server.Handler

	hs = append(hs, dag.NewHandler(
		&dag.NewHandlerArgs{
			Engine:             params.Engine,
			LogEncodingCharset: params.Config.LogEncodingCharset,
		},
	))

	serverParams := server.Params{
		Host:     params.Config.Host,
		Port:     params.Config.Port,
		TLS:      params.Config.TLS,
		Logger:   params.Logger,
		Handlers: hs,
		AssetsFS: assetsFS,
	}

	if params.Config.IsAuthToken {
		serverParams.AuthToken = &server.AuthToken{
			Token: params.Config.AuthToken,
		}
	}

	if params.Config.IsBasicAuth {
		serverParams.BasicAuth = &server.BasicAuth{
			Username: params.Config.BasicAuthUsername,
			Password: params.Config.BasicAuthPassword,
		}
	}

	return server.NewServer(serverParams, params.Config)
}
