package router

func (r *router) setTripRoutes() {
	tripRouter := r.Engine.Group("/trip")

	tripRouter.GET("", r.Service.Trip.List)
	tripRouter.GET("/:id", r.Service.Trip.Get)
	tripRouter.POST("", r.Service.Trip.Create)
}
