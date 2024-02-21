package core

import (
	"context"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/api"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func Core() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(storage.NewPostgresStorage, fx.As(new(service.Storage))),
			fx.Annotate(storage.NewRedisStorage, fx.As(new(service.RStorage))),
		),
		fx.Provide(
			context.Background,
			storage.NewPostgresStorage,
			storage.NewRedisStorage,
			config.NewConfig,
			gin.Default,
			api.NewApi,
			service.NewTokenManager,
			service.NewServices,
		),
		fx.Invoke(restAPIHook),
	)
}

func restAPIHook(lifecycle fx.Lifecycle, api *api.Api) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go api.Run()
				return nil
			},
		},
	)
}
