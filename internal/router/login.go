package router

func (r *router) setLoginRoutes() {
	r.Engine.GET("/oauthUrl", r.Service.User.OauthUrl)
	r.Engine.POST("/user/login", r.Service.User.OAuthLogin)
}
