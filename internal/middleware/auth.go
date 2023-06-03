package middleware

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/CoRide-tw/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if len(authHeader) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "No authorization header",
			})
			c.Abort()
			return
		}

		mat, err := regexp.MatchString(`Bearer.*`, authHeader)
		if err != nil || !mat {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Malformed authorization header",
			})
			c.Abort()
			return
		}

		var parsedClaims model.Claims
		token := strings.Split(authHeader, " ")[1]
		tokenClaims, err := jwt.ParseWithClaims(
			token, &parsedClaims,
			func(token *jwt.Token) (i interface{}, err error) {
				return []byte(secret), nil
			},
		)

		// add user id into context
		c.Set("user_id", parsedClaims.ID)

		if err != nil {
			var message string
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					message = "TOKEN_MALFORMED"
				} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
					message = "TOKEN_NOT_VERIFIED"
				} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
					message = "SIGNATURE_VALIDATION_FAILED"
				} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
					message = "TOKEN_EXPIRED"
				} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
					message = "TOKEN_NOT_YET_VALID"
				} else {
					message = "TOKEN_UNEXPECTED_ERROR"
				}
			}
			c.JSON(
				401, gin.H{
					"error": message,
				},
			)
			c.Abort()
			return
		}

		if _, ok := tokenClaims.Claims.(*model.Claims); ok && tokenClaims.Valid {
			c.Next()
		} else {
			c.Abort()
			return
		}
	}
}
