package router

import "github.com/gin-gonic/gin"

func (r *router) setLoginRoutes() {
	r.Engine.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to CoRide API service, health check passed.",
		})
	})
	r.Engine.GET("/oauthUrl", r.Service.User.OauthUrl)
	r.Engine.POST("/user/login", r.Service.User.OAuthLogin)
}
