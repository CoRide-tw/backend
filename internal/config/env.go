package config

import (
	"log"

	"github.com/joho/godotenv"
)

var Env *env

type env struct {
	CoRideJwtSecret string
}

func LoadEnv() *env {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &env{}
}
