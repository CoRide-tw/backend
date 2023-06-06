package router

func (r *router) setRouteRoutes() {
	routeRouter := r.Engine.Group("/route")

	routeRouter.GET("/ranking", r.Service.Route.ListNearestRoutes)
	routeRouter.GET("/:id", r.Service.Route.Get)
	routeRouter.POST("", r.Service.Route.Create)
	routeRouter.DELETE("/:id", r.Service.Route.Delete)
}
