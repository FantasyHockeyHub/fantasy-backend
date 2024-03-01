package api

import (
	"errors"
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/Frozen-Fantasy/fantasy-backend.git/docs"
	_ "github.com/Frozen-Fantasy/fantasy-backend.git/docs"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/storage"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Api struct {
	router   *gin.Engine
	cfg      config.ServiceConfiguration
	services *service.Services
}

func NewApi(
	router *gin.Engine,
	cfg config.ServiceConfiguration,
) *Api {
	svc := &Api{
		router: router,
		cfg:    cfg,
		services: service.NewServices(service.Deps{
			Cfg:      cfg,
			Storage:  storage.NewPostgresStorage(cfg),
			RStorage: storage.NewRedisStorage(cfg),
			Jwt:      service.NewTokenManager(cfg),
		}),
	}
	svc.router.Use(CORSMiddleware())
	svc.registerRoutes()
	return svc
}

// @BasePath /api/v1/
const BasePath = "/api/v1/"

func (api Api) Run() {
	cfg := config.Load()
	docs.SwaggerInfo.BasePath = BasePath

	api.router.Run(cfg.Api.GetAddr())
}

func (api *Api) registerRoutes() {
	base := api.router.Group(BasePath)

	api.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	auth := base.Group("/auth")
	{
		auth.POST("/sign-up", api.signUp)
		auth.POST("/sign-in", api.signIn)
		auth.POST("/email/send-code", api.sendVerificationCode)
		auth.POST("/refresh-tokens", api.refreshTokens)
		auth.POST("/logout", api.logout)
	}
	user := base.Group("/user")
	{
		user.GET("/exists", api.checkUserDataExists)
		userAuthenticated := user.Group("/", api.userIdentity)
		{
			userAuthenticated.GET("/info", api.userInfo)
			userAuthenticated.PATCH("password/change", api.changePassword)
		}
		password := user.Group("/password")
		{
			password.POST("/forgot", api.forgotPassword)
			password.PATCH("/reset", api.resetPassword)
		}
	}

	team := base.Group("/tournament")
	team.GET("/create_team_nhl", api.CreateTeamsNHL)
	team.GET("/create_team_khl", api.CreateTeamsKHL)
	team.GET("/events_day_khl", api.EventsKHL)
	team.GET("/events_day_nhl", api.EventsNHL)
	team.GET("/get_matches/:league", api.GetMatches)
	team.GET("/create_tournaments", api.CreateTournaments)
}

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

const (
	UnauthorizedErrorTitle     = "Ошибка авторизации"
	InternalServerErrorTitle   = "Ошибка произошла на стороне сервера"
	InternalServerErrorMessage = "Ошибка на сервере. Зайдите позже :("
	NotFoundErrorTitle         = "Ошибка произошла на стороне сервера"
	NotFoundErrorMessage       = "Ошибка на сервере. Зайдите позже :("
	BadRequestErrorTitle       = "Программная ошибка"
)

var (
	InvalidInputBodyError       = errors.New("invalid input body")
	InvalidInputParametersError = errors.New("invalid input parameters")
)

func getUnauthorizedError(err error) Error {
	return Error{
		Error:   UnauthorizedErrorTitle,
		Message: err.Error(),
	}
}

func getInternalServerError() Error {
	return Error{
		Error:   InternalServerErrorTitle,
		Message: InternalServerErrorMessage,
	}
}

func getNotFoundError() Error {
	return Error{
		Error:   NotFoundErrorTitle,
		Message: NotFoundErrorMessage,
	}
}

func getBadRequestError(err error) Error {
	return Error{
		Error:   BadRequestErrorTitle,
		Message: err.Error(),
	}
}

type StatusResponse struct {
	Status string `json:"status"`
}

type CheckEntityExistsResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
