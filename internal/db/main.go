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
		log.Println("Init user table failed")
		return err
	}
	if err := initRouteTable(); err != nil {
		log.Println("Init route table failed")
		return err
	}
	if err := initRequestTable(); err != nil {
		log.Println("Init request table failed")
		return err
	}

	return nil
}
