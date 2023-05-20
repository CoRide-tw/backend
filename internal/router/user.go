package router

func (r *router) setUserRoutes() {
	userRouter := r.Engine.Group("/user")

	userRouter.GET("/:id", r.Service.User.Get)
	userRouter.PUT("/:id", r.Service.User.Update)
}
