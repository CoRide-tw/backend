package db

import (
	"context"
	. "github.com/CoRide-tw/backend/internal/errors/generated/dberr"
	. "github.com/DenChenn/blunder/pkg/blunder"
	"github.com/jackc/pgx/v5"

	"github.com/CoRide-tw/backend/internal/model"
)

const createTripTableSQL = `
	CREATE TABLE IF NOT EXISTS trips (
		id SERIAL,
		rider_id INT NOT NULL,
		driver_id INT NOT NULL,
		request_id INT NOT NULL,
		route_id INT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
		deleted_at TIMESTAMP
	);
`

func initTripTable() error {
	if _, err := DBClient.pgPool.Exec(context.Background(), createTripTableSQL); err != nil {
		return err
	}
	return nil
}

const listTripByRiderIdSQL = `
	SELECT *
	FROM trips
	WHERE rider_id = $1 AND deleted_at IS NULL;
`

func ListTripByRiderId(riderId int32) ([]*model.Trip, error) {
	var trips []*model.Trip
	rows, err := DBClient.pgPool.Query(context.Background(), listTripByRiderIdSQL, riderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var trip model.Trip
		if err := rows.Scan(
			&trip.Id,
			&trip.RiderId,
			&trip.DriverId,
			&trip.RequestId,
			&trip.RouteId,
			&trip.CreatedAt,
			&trip.DeletedAt,
		); err != nil {
			return nil, ErrUndefined.WithCustomMessage(err.Error())
		}
		trips = append(trips, &trip)
	}
	return trips, nil
}

const listTripByDriverIdSQL = `
	SELECT *
	FROM trips
	WHERE driver_id = $1 AND deleted_at IS NULL;
`

func ListTripByDriverId(driverId int32) ([]*model.Trip, error) {
	var trips []*model.Trip
	rows, err := DBClient.pgPool.Query(context.Background(), listTripByDriverIdSQL, driverId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var trip model.Trip
		if err := rows.Scan(
			&trip.Id,
			&trip.RiderId,
			&trip.DriverId,
			&trip.RequestId,
			&trip.RouteId,
			&trip.CreatedAt,
			&trip.DeletedAt,
		); err != nil {
			return nil, ErrUndefined.WithCustomMessage(err.Error())
		}
		trips = append(trips, &trip)
	}
	return trips, nil
}

const getTripSQL = `
	SELECT *
	FROM trips
	WHERE id = $1 AND deleted_at IS NULL;
`

func GetTrip(id int32) (*model.Trip, error) {
	var trip model.Trip
	if err := DBClient.pgPool.QueryRow(context.Background(), getTripSQL, id).Scan(
		&trip.Id,
		&trip.RiderId,
		&trip.DriverId,
		&trip.RequestId,
		&trip.RouteId,
		&trip.CreatedAt,
		&trip.DeletedAt,
	); err != nil {
		return nil, Match(err, pgx.ErrNoRows, ErrTripNotFound).Return()
	}
	return &trip, nil
}

const createTripSQL = `
	INSERT INTO trips (rider_id, driver_id, request_id, route_id)
	VALUES (
		$1,
		$2, 
		$3, 
		$4
	)
	RETURNING id, created_at;
`

func CreateTrip(trip *model.Trip) (*model.Trip, error) {
	if err := DBClient.pgPool.QueryRow(context.Background(), createTripSQL,
		trip.RiderId,
		trip.DriverId,
		trip.RequestId,
		trip.RouteId,
	).Scan(
		&trip.Id,
		&trip.CreatedAt,
	); err != nil {
		return nil, ErrUndefined.WithCustomMessage(err.Error())
	}
	return trip, nil
}
