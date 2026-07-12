package main

import (
	"github.com/khaivutri/bookmark-service/internal/api"
	"github.com/khaivutri/bookmark-service/pkg/logger"
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
)

//@title Bookmark Service API
//@version 1.0
//@description This is a simple REST API for a bookmark service.
//@BasePath /
func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	//set log level 
	logger.SetLogLevel(cfg.LogLevel)

	//set redis client
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}
	
	engine := api.NewEngine(cfg, redisClient)
	err = engine.Start()

	if err != nil {
		panic(err)
	}
}