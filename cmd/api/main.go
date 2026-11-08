package main

import (
	_ "bookmark-management/docs"
	"bookmark-management/internal/api"
	"bookmark-management/internal/pkg/redis"
)

// @title Bookmark Management API
// @version 1.0.0
// @description A simple link shortening API
// @BasePath /
func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	redisClient, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}
	defer redisClient.Close()

	app := api.NewEngine(cfg, redisClient)
	err = app.Start()

	if err != nil {
		panic(err)
	}
}
