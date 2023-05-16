package router

import (
	"github.com/CoRide-tw/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type router struct {
	Engine  *gin.Engine
	Service *service.Service
}

func NewRouterEngine(engine *gin.Engine, service *service.Service) *gin.Engine {
	router := &router{
		Engine:  engine,
		Service: service,
	}

	return router.Engine
}
