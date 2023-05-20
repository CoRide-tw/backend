package db

import (
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

	return nil
}
