package util

import (
	"time"

	"github.com/CoRide-tw/backend/internal/model"
	"github.com/golang-jwt/jwt"
)

func GenerateJWT(id int32, secret string) (*string, error) {
	now := time.Now()
	claims := model.Claims{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 0,
			IssuedAt:  now.Unix(),
			Issuer:    "CoRide",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenClaims.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &token, nil
}
