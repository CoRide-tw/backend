package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	cors_config := cors.DefaultConfig()

	cors_config.AllowOrigins = []string{"*"}
	cors_config.AllowCredentials = true
	cors_config.AddAllowHeaders("Authorization")

	return cors.New(cors_config)
}
