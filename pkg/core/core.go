package core

import (
	"context"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/api"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/user"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func Core() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(storage.NewPostgresStorage, fx.As(new(user.Storage))),
		),
		fx.Provide(
			context.Background,
			storage.NewPostgresStorage,
			config.NewConfig,
			gin.Default,
			api.NewApi,
			user.NewManager,
			user.NewService,
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
