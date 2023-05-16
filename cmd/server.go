package main

import (
	"github.com/CoRide-tw/backend/internal/config"
	"github.com/CoRide-tw/backend/internal/router"
	"github.com/CoRide-tw/backend/internal/service"
	"github.com/gin-gonic/gin"
)

func init() {
	config.Env = config.LoadEnv()
}

func main() {
	engine := gin.Default()
	service := service.NewService()

	server := router.NewRouterEngine(engine, service)
	panic(server.Run())
}
