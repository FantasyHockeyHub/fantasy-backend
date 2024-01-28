package api

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/Frozen-Fantasy/fantasy-backend.git/docs"
	_ "github.com/Frozen-Fantasy/fantasy-backend.git/docs"
	"github.com/Frozen-Fantasy/fantasy-backend.git/pkg/service/user"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Api struct {
	router *gin.Engine
	cfg    config.ServiceConfiguration
	user   *user.Service
}

func NewApi(
	router *gin.Engine,
	cfg config.ServiceConfiguration,
	user *user.Service,
) *Api {
	svc := &Api{
		router: router,
		cfg:    cfg,
		user:   user,
	}
	//svc.router.Use(CORSMiddleware())
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

	//baseWithAuth := base.Group("/")
	//baseWithAuth.Use(api.AuthMW())

	auth := base.Group("/auth")
	{
		auth.POST("/sign-up", api.SignUp)
		auth.POST("/sign-in", api.SignIn)
		auth.POST("/email/send-code", api.SendVerificationCode)
		auth.POST("/refresh-tokens", api.RefreshTokens)
	}

	user := base.Group("/user")
	{
		user.POST("/check-email", api.CheckEmailExists)
		user.POST("/check-nickname", api.CheckNicknameExists)
	}

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
