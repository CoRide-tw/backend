package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Env *env

type env struct {
	PostgresDatabaseUrl     string
	GoogleOAuthClientId     string
	GoogleOAuthClientSecret string
	GoogleOAuthRedirectUrl  string
	GoogleOauthScope        string
	CoRideJwtSecret         string
	GoogleMapsApiKey        string
}

func LoadEnv() *env {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &env{
		PostgresDatabaseUrl:     os.Getenv("POSTGRES_DATABASE_URL"),
		GoogleOAuthClientId:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		GoogleOAuthClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		GoogleOAuthRedirectUrl:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
		GoogleOauthScope:        os.Getenv("GOOGLE_OAUTH_SCOPE"),
		CoRideJwtSecret:         os.Getenv("CORIDE_JWT_SECRET"),
		GoogleMapsApiKey:        os.Getenv("GOOGLE_MAPS_API_KEY"),
	}
}
