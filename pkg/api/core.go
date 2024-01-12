package api

import (
	"github.com/Frozen-Fantasy/fantasy-backend.git/config"
	"github.com/Frozen-Fantasy/fantasy-backend.git/docs"
	_ "github.com/Frozen-Fantasy/fantasy-backend.git/docs"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Api struct {
	router *gin.Engine
	cfg    config.ServiceConfiguration
}

func NewApi(
	router *gin.Engine,
	cfg config.ServiceConfiguration,
) *Api {
	svc := &Api{
		router: router,
		cfg:    cfg,
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
	//base := api.router.Group(BasePath)

	api.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//baseWithAuth := base.Group("/")
	//baseWithAuth.Use(api.AuthMW())
}
