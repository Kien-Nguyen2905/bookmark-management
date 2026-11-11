package api

import (
	"fmt"

	"bookmark-management/internal/handler"
	"bookmark-management/internal/repository"
	"bookmark-management/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Engine is the HTTP application.
type Engine interface {
	Start() error
	GetApp() *gin.Engine
}

type engine struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
}

func NewEngine(c *Config, redisClient *redis.Client) Engine {
	e := &engine{app: gin.Default(), cfg: c, redisClient: redisClient}
	e.initRoutes()
	return e
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%d", e.cfg.AppPort))
}

func (e *engine) GetApp() *gin.Engine {
	return e.app
}

func (e *engine) initRoutes() {
	healthCheckRepo := repository.NewHealthCheck(e.redisClient)
	healthCheckSvc := service.NewHealthCheck(e.cfg.AppName, e.cfg.InstanceID, healthCheckRepo)
	healthCheckHandler := handler.NewHealthCheck(healthCheckSvc)

	genPassSvc := service.NewGenPass()
	urlStorageRepo := repository.NewURLStorage(e.redisClient)
	shortenURLSvc := service.NewShortenUrl(urlStorageRepo, genPassSvc)
	shortenHandler := handler.NewShortenLink(shortenURLSvc)

	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.app.GET("/health-check", healthCheckHandler.Check)
	e.app.POST("/v1/links/shorten", shortenHandler.CreateShortenLink)
}
