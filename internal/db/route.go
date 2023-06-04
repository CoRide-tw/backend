package db

import (
	"context"
	"time"

	"github.com/CoRide-tw/backend/internal/model"
)

const createRouteTableSQL = `
	CREATE EXTENSION IF NOT EXISTS postgis;

	CREATE TABLE IF NOT EXISTS routes (
		id SERIAL,
		driver_id SERIAL NOT NULL,
		start_location GEOMETRY(Point, 4326) NOT NULL,
		end_location GEOMETRY(Point, 4326) NOT NULL,
		start_time TIMESTAMP WITH TIME ZONE NOT NULL,
		end_time TIMESTAMP WITH TIME ZONE NOT NULL,
		capacity INT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
		deleted_at TIMESTAMP WITH TIME ZONE
	);

	CREATE INDEX ON routes USING gist(start_location, end_location);
`

func initRouteTable() error {
	if _, err := DBClient.pgPool.Exec(context.Background(), createRouteTableSQL); err != nil {
		return err
	}
	return nil
}

const getRouteSQL = `
	SELECT 
		id, 
		driver_id,
		ST_X(start_location), ST_Y(start_location), 
		ST_X(end_location), ST_Y(end_location), 
		start_time, end_time, 
		capacity, 
		created_at, updated_at, deleted_at
	FROM routes
	WHERE id = $1 AND deleted_at IS NULL;
`

func GetRoute(id int32) (*model.Route, error) {
	var route model.Route
	if err := DBClient.pgPool.QueryRow(context.Background(), getRouteSQL, id).Scan(
		&route.Id,
		&route.DriverId,
		&route.StartLong,
		&route.StartLat,
		&route.EndLong,
		&route.EndLat,
		&route.StartTime,
		&route.EndTime,
		&route.Capacity,
		&route.CreatedAt,
		&route.UpdatedAt,
		&route.DeletedAt,
	); err != nil {
		return nil, err
	}
	return &route, nil
}

// Note: ST_MakePoint(longitude, latitude)
const listNearestRouteSQL = `
	WITH rider_requirements AS (
		SELECT 
			ST_SetSRID(ST_MakePoint($1, $2), 4326) AS pickup_point,
			ST_SetSRID(ST_MakePoint($3, $4), 4326) AS dropoff_point,
			$5::timestamp with time zone AS pickup_start_time,
			$6::timestamp with time zone AS pickup_end_time
	)
	SELECT 
		id,
		driver_id,
		ST_X(start_location),
		ST_Y(start_location),
		ST_X(end_location),
		ST_Y(end_location),
		start_time,
		end_time,
		capacity,
		created_at,
		updated_at,
		deleted_at
	FROM 
		routes, rider_requirements
	WHERE 
		deleted_at IS NULL 
		AND start_time <= (SELECT pickup_start_time FROM rider_requirements)
		AND end_time >= (SELECT pickup_end_time FROM rider_requirements)
	ORDER BY (
		ST_Distance(start_location, (SELECT pickup_point FROM rider_requirements)) + 
		ST_Distance(end_location, (SELECT dropoff_point FROM rider_requirements))
	) ASC
	LIMIT 30;
`

func ListNearestRoutes(
	pickupLong, pickupLat, dropoffLong, dropoffLat float64,
	pickupStartTime, pickupEndTime time.Time,
) ([]*model.Route, error) {
	rows, err := DBClient.pgPool.Query(context.Background(), listNearestRouteSQL,
		pickupLong, pickupLat, dropoffLong, dropoffLat, pickupStartTime, pickupEndTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []*model.Route
	for rows.Next() {
		var route model.Route
		if err := rows.Scan(
			&route.Id,
			&route.DriverId,
			&route.StartLong,
			&route.StartLat,
			&route.EndLong,
			&route.EndLat,
			&route.StartTime,
			&route.EndTime,
			&route.Capacity,
			&route.CreatedAt,
			&route.UpdatedAt,
			&route.DeletedAt,
		); err != nil {
			return nil, err
		}
		routes = append(routes, &route)
	}
	return routes, nil
}

const createRouteSQL = `
	INSERT INTO routes (driver_id, start_location, end_location, start_time, end_time, capacity)
	VALUES (
		$1, 
		ST_SetSRID(ST_MakePoint($2, $3), 4326), 
		ST_SetSRID(ST_MakePoint($4, $5), 4326), 
		$6, 
		$7, 
		$8
	)
	RETURNING id, created_at, updated_at;
`

func CreateRoute(route *model.Route) (*model.Route, error) {
	if err := DBClient.pgPool.QueryRow(context.Background(), createRouteSQL,
		route.DriverId,
		route.StartLong,
		route.StartLat,
		route.EndLong,
		route.EndLat,
		route.StartTime,
		route.EndTime,
		route.Capacity,
	).Scan(
		&route.Id,
		&route.CreatedAt,
		&route.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return route, nil
}

const updateRouteSQL = `
	UPDATE routes SET
		start_long = $3,
		start_lat = $4,
		end_long = $5,
		end_lat = $6,
		start_time = $7,
		end_time = $8,
		capacity = $9,
		updated_at = NOW()
	WHERE id = $1, AND driver_id = $2 AND deleted_at IS NULL
	RETURNING
		id, 
		driver_id,
		ST_X(start_location), ST_Y(start_location), 
		ST_X(end_location), ST_Y(end_location), 
		start_time, end_time, 
		capacity, 
		created_at, updated_at;
`

func UpdateRoute(route *model.Route) (*model.Route, error) {
	if err := DBClient.pgPool.QueryRow(context.Background(), updateRouteSQL,
		route.Id,
		route.DriverId,
		route.StartLong,
		route.StartLat,
		route.EndLong,
		route.EndLat,
		route.StartTime,
		route.EndTime,
		route.Capacity,
	).Scan(
		&route.Id,
		&route.DriverId,
		&route.StartLong,
		&route.StartLat,
		&route.EndLong,
		&route.EndLat,
		&route.StartTime,
		&route.EndTime,
		&route.Capacity,
		&route.CreatedAt,
		&route.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return route, nil
}
