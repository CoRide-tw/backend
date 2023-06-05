package db

import (
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pgPool *pgxpool.Pool
}

var DBClient *DB

func InitDBClient(pgPool *pgxpool.Pool) error {

	// create db client
	DBClient = &DB{
		pgPool: pgPool,
	}

	// init tables
	if err := initUserTable(); err != nil {
		log.Println("Init users table failed")
		return err
	}
	if err := initRouteTable(); err != nil {
		log.Println("Init routes table failed")
		return err
	}
	if err := initRequestTable(); err != nil {
		log.Println("Init requests table failed")
		return err
	}
	if err := initTripTable(); err != nil {
		log.Println("Init trips table failed")
		return err
	}

	return nil
}
