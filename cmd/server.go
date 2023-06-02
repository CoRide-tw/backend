package main

import (
	"context"
	"log"

	"github.com/CoRide-tw/backend/internal/config"
	"github.com/CoRide-tw/backend/internal/db"
	"github.com/CoRide-tw/backend/internal/router"
	"github.com/CoRide-tw/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func init() {
	config.Env = config.LoadEnv()
}

func main() {
	// database connection
	pgPool, err := pgxpool.New(context.Background(), config.Env.PostgresDatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer pgPool.Close()
	db.InitDBClient(pgPool)

	engine := gin.Default()
	service := service.NewService()

	server := router.NewRouterEngine(engine, service)
	panic(server.Run())
}
