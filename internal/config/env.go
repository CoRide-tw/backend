package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var Env *env

type env struct {
	PostgresDatabaseUrl string
}

func LoadEnv() *env {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	return &env{
		PostgresDatabaseUrl: os.Getenv("POSTGRES_DATABASE_URL"),
	}
}
