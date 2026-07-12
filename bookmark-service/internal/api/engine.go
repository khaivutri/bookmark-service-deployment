package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/docs"
	"github.com/khaivutri/bookmark-service/internal/handler"
	v1 "github.com/khaivutri/bookmark-service/internal/handler/v1"
	"github.com/khaivutri/bookmark-service/internal/repository"
	"github.com/khaivutri/bookmark-service/internal/service"
	"github.com/khaivutri/bookmark-service/pkg/utils"
	"github.com/redis/go-redis/v9"

	_ "github.com/khaivutri/bookmark-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Engine defines the interface for the API engine.
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app 		*gin.Engine
	cfg 		*Config

	redis 		*redis.Client
}

// NewEngine creates and returns a new Engine instance with initialized routes.
func NewEngine(cfg *Config, client *redis.Client) Engine{
	app := &engine{
		app : gin.Default(),
		cfg : cfg,
		redis : client,
	}
	app.initRoutes()
	return app
}

// Start runs the API server on the configured port.
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP serves the HTTP request using the underlying Gin engine.
func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}
	
func (e *engine) initRoutes(){
	redisPinger := repository.NewRedisPinger(e.redis)
	healthCheckSvc := service.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceId, redisPinger )
	healthCheck := handler.NewHealthCheck(healthCheckSvc)

	urlStorage := repository.NewURLStorage(e.redis)
	shortenURLSvc := service.NewURLStorage(urlStorage, utils.NewGenCode())
	shortenURL := v1.NewShortenURL(shortenURLSvc)
	e.app.HandleMethodNotAllowed = true	
	
	docs.SwaggerInfo.BasePath = e.cfg.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	e.app.GET("/health-check", healthCheck.HealthCheck)

	v1 := e.app.Group("/v1") 
	{
		links := v1.Group("/links")
		{
			links.POST("/shorten", shortenURL.CreateShortenLink)
			links.GET("/redirect/:code", shortenURL.Redirect)
		}
	}
}

