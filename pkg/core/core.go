package core

import (
	"context"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/api"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/jobs/get_events"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/jobs/update_events"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/events"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func Core() fx.Option {
	return fx.Options(
		fx.Provide(
			fx.Annotate(storage.NewPostgresStorage, fx.As(new(service.UserStorage))),
			fx.Annotate(storage.NewRedisStorage, fx.As(new(service.UserRStorage))),
			fx.Annotate(storage.NewPostgresStorage, fx.As(new(events.EventsStorage))),
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
			events.NewEventsService,
			get_events.NewGetHokeyEvents,
			update_events.NewUpdateHockeyEvents,
			update_events.NewUpdateHockeyEventsKHL,
		),
		fx.Invoke(restAPIHook),
		fx.Invoke(getHokeyEventsHook),
		fx.Invoke(updateHokeyEventsHook),
		fx.Invoke(updateHokeyEventsHookKHL),
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

func getHokeyEventsHook(lifecycle fx.Lifecycle, job *get_events.GetHokeyEvents) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go job.Start(context.Background())
				return nil
			},
		},
	)
}

func updateHokeyEventsHook(lifecycle fx.Lifecycle, job *update_events.UpdateHockeyEvents) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go job.Start(context.Background())
				return nil
			},
		},
	)
}

func updateHokeyEventsHookKHL(lifecycle fx.Lifecycle, job *update_events.UpdateHockeyEventsKHL) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go job.StartKHL(context.Background())
				return nil
			},
		},
	)
}
