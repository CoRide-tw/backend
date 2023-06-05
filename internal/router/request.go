package router

func (r *router) setRequestRoutes() {
	routeRouter := r.Engine.Group("/request")

	routeRouter.GET("/", r.Service.Request.List)
	routeRouter.GET("/:id", r.Service.Request.Get)
	routeRouter.POST("/", r.Service.Request.Create)
	routeRouter.PATCH("/:id", r.Service.Request.Update)
	routeRouter.PATCH("/:id/status", r.Service.Request.Deny)
	routeRouter.DELETE("/:id", r.Service.Request.Delete)
}
