package router

func (r *router) setRequestRoutes() {
	requestRouter := r.Engine.Group("/request")

	requestRouter.GET("", r.Service.Request.List)
	requestRouter.GET("/:id", r.Service.Request.Get)
	requestRouter.POST("", r.Service.Request.Create)
	requestRouter.PATCH("/:id", r.Service.Request.Update)
	requestRouter.PATCH("/:id/status", r.Service.Request.Deny)
	requestRouter.DELETE("/:id", r.Service.Request.Delete)
}
