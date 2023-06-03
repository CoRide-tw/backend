package model

import "github.com/golang-jwt/jwt"

type Claims struct {
	ID int32
	jwt.StandardClaims
}
