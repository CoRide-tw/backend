package router

import (
	"github.com/CoRide-tw/backend/internal/middleware"
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

	// use CORS middleware
	router.useCorsMiddleware()

	// set routes
	router.setUserRoutes()

	return router.Engine
}

func (r *router) useCorsMiddleware() {
	r.Engine.Use(middleware.Cors())
}
