package db

import (
	"context"
	"time"

	. "github.com/CoRide-tw/backend/internal/errors/generated/dberr"
	"github.com/CoRide-tw/backend/internal/model"
	. "github.com/DenChenn/blunder/pkg/blunder"
	"github.com/jackc/pgx/v5"
)

const createRouteTableSQL = `
	CREATE EXTENSION IF NOT EXISTS postgis;

	CREATE TABLE IF NOT EXISTS routes (
		id SERIAL,
		driver_id INT NOT NULL,
		start_location GEOMETRY(Point, 4326) NOT NULL,
		end_location GEOMETRY(Point, 4326) NOT NULL,
		start_time TIMESTAMP WITH TIME ZONE NOT NULL,
		end_time TIMESTAMP WITH TIME ZONE NOT NULL,
		capacity INT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
		deleted_at TIMESTAMP WITH TIME ZONE
	);

	CREATE INDEX IF NOT EXISTS routes_start_location_end_location_idx 
		ON routes USING gist(start_location, end_location);
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
		return nil, Match(err, pgx.ErrNoRows, ErrRouteNotFound).Return()
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
		r.id,
		r.driver_id,
		ST_X(r.start_location),
		ST_Y(r.start_location),
		ST_X(r.end_location),
		ST_Y(r.end_location),
		r.start_time,
		r.end_time,
		r.capacity,
		r.created_at,
		r.updated_at,
		r.deleted_at,
		u.name,
		u.picture_url,
		u.car_type,
		u.car_plate
	FROM rider_requirements, routes r JOIN users u 
	ON r.driver_id = u.id
	WHERE 
		r.deleted_at IS NULL 
		AND start_time <= (SELECT pickup_start_time FROM rider_requirements)
		AND end_time >= (SELECT pickup_end_time FROM rider_requirements)
	ORDER BY (
		ST_Distance(start_location, (SELECT pickup_point FROM rider_requirements)) + 
		ST_Distance(end_location, (SELECT dropoff_point FROM rider_requirements))
	) ASC
	LIMIT 30
`

type ListNearestRoutesQueryResp struct {
	Id               int32      `json:"id"`
	DriverId         int32      `json:"driverId"`
	StartLong        float64    `json:"startLong"`
	StartLat         float64    `json:"startLat"`
	EndLong          float64    `json:"endLong"`
	EndLat           float64    `json:"endLat"`
	StartTime        time.Time  `json:"startTime"`
	EndTime          time.Time  `json:"endTime"`
	Capacity         int32      `json:"capacity"`
	CreatedAt        time.Time  `json:"createdAt"`
	UpdatedAt        time.Time  `json:"updatedAt"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty"`
	DriverName       string     `json:"driverName"`
	DriverPictureUrl string     `json:"driverPictureUrl"`
	DriverCarType    string     `json:"driverCarType"`
	DriverCarPlate   string     `json:"driverCarPlate"`
}

func ListNearestRoutes(
	pickupLong, pickupLat, dropoffLong, dropoffLat float64,
	pickupStartTime, pickupEndTime time.Time,
) ([]*ListNearestRoutesQueryResp, error) {
	rows, err := DBClient.pgPool.Query(context.Background(), listNearestRouteSQL,
		pickupLong, pickupLat, dropoffLong, dropoffLat, pickupStartTime, pickupEndTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*ListNearestRoutesQueryResp
	for rows.Next() {
		var item ListNearestRoutesQueryResp
		if err := rows.Scan(
			&item.Id,
			&item.DriverId,
			&item.StartLong,
			&item.StartLat,
			&item.EndLong,
			&item.EndLat,
			&item.StartTime,
			&item.EndTime,
			&item.Capacity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.DeletedAt,
			&item.DriverName,
			&item.DriverPictureUrl,
			&item.DriverCarType,
			&item.DriverCarPlate,
		); err != nil {
			return nil, ErrUndefined.WithCustomMessage(err.Error())
		}
		items = append(items, &item)
	}
	return items, nil
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
		return nil, ErrUndefined.WithCustomMessage(err.Error())
	}
	return route, nil
}

const deleteRouteSQL = `
	UPDATE routes SET 
		deleted_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL;
`

func DeleteRoute(id int32) error {
	if _, err := DBClient.pgPool.Exec(context.Background(), deleteRouteSQL, id); err != nil {
		return ErrUndefined.WithCustomMessage(err.Error())
	}
	return nil
}
